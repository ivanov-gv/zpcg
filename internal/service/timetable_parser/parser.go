package timetable_parser

import (
	"fmt"
	"sort"

	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/config/timetable_parser_config"
	"github.com/ivanov-gv/zpcg/internal/model/stations"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/station_blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
)

func New(config timetable_parser_config.Config) TimetableParser {
	return TimetableParser{
		seasonParser:           newSeasonParser(config.Timetable, newStationParser()),
		unifiedStationNameList: nil,
		transferStationId:      0,
	}
}

type TimetableParser struct {
	*seasonParser
	unifiedStationNameList [][]rune
	transferStationId      timetable.StationId
}

func (t *TimetableParser) ParseTimetable() (timetable.ExportFormat, error) {
	// parse seasons' timetables
	err := t.parseSeasons() //
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("seasons.ParseSeasons: %w", err)
	}

	// fill unifiedStationNameList
	unifiedStationNameStrings := lo.Keys(t.unifiedStationNameToStationIdMap)
	// sort for more predictable output
	sort.Strings(unifiedStationNameStrings)
	t.unifiedStationNameList = lo.Map(unifiedStationNameStrings, func(item string, _ int) []rune { return []rune(item) })

	// get transfer station id
	t.transferStationId = t.unifiedStationNameToStationIdMap[station_name_resolver.Unify(timetable.TransferStationName)]

	// add blacklisted stations
	t.unifiedStationNameList, t.unifiedStationNameToStationIdMap, err =
		AddBlacklistedStations(t.unifiedStationNameList, t.unifiedStationNameToStationIdMap)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("AddBlacklistedStations: %w", err)
	}
	// add aliases
	t.unifiedStationNameList, t.unifiedStationNameToStationIdMap, err = AddAliases(t.unifiedStationNameList, t.unifiedStationNameToStationIdMap)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("AddAliases: %w", err)
	}

	return timetable.ExportFormat{
		Seasons:                          t.seasonTimetables,
		UnifiedStationNameToStationIdMap: t.unifiedStationNameToStationIdMap,
		UnifiedStationNameList:           t.unifiedStationNameList,
		StationIdToStationMap:            t.stationIdToStationMap,
		StationTypes:                     t.stationTypesMap,
		TransferStationId:                t.transferStationId,
	}, nil
}

var (
	// AliasesAsUnifiedStationNames is a list of unified stations aliases
	AliasesAsUnifiedStationNames = func() []string {
		var result []string
		for _, station := range stations.AliasesStationsList {
			unifiedAliases := lo.Map(station.Aliases, func(alias string, _ int) string {
				return station_name_resolver.Unify(alias)
			})
			result = append(result, unifiedAliases...)
		}
		return result
	}()

	// AliasesOriginalUnifiedStationNameToUnifiedAliasesMap is a map of unified station name -> slice of unified aliases
	AliasesOriginalUnifiedStationNameToUnifiedAliasesMap = lo.SliceToMap(stations.AliasesStationsList,
		func(stationAliases timetable.StationAliases) (string, []string) {
			return station_name_resolver.Unify(stationAliases.StationName), lo.Map(stationAliases.Aliases, station_name_resolver.UnifyIterative)
		})
)

func AddAliases(unifiedStationNameList [][]rune, unifiedStationNameToStationIdMap map[string]timetable.StationId) ([][]rune, map[string]timetable.StationId, error) {
	// add aliases to station list
	unifiedStationNameList = append(unifiedStationNameList,
		lo.Map(AliasesAsUnifiedStationNames, func(item string, _ int) []rune { return []rune(item) })...)
	// add mapping from aliases to station id
	for stationName, aliases := range AliasesOriginalUnifiedStationNameToUnifiedAliasesMap {
		stationId, ok := unifiedStationNameToStationIdMap[stationName]
		if !ok { // station not found
			return nil, nil, fmt.Errorf("can't find station to add aliases to: stationName = %s", stationName)
		}
		for _, _alias := range aliases {
			unifiedStationNameToStationIdMap[_alias] = stationId
		}
	}
	return unifiedStationNameList, unifiedStationNameToStationIdMap, nil
}

func AddBlacklistedStations(_unifiedStationNameList [][]rune, _unifiedStationNameToStationIdMap map[string]timetable.StationId) (unifiedStationNameList [][]rune, unifiedStationNameToStationIdMap map[string]timetable.StationId, err error) {
	// station name list - UnifiedStationNameList
	_unifiedStationNameList = append(_unifiedStationNameList, lo.Map(station_blacklist.UnifiedNames,
		func(item string, _ int) []rune { return []rune(item) })...)

	// map: station name -> station id - UnifiedStationNameToStationIdMap
	newMap := lo.Assign(_unifiedStationNameToStationIdMap, station_blacklist.UnifiedStationNameToStationIdMap)
	// ensure the map is not broken
	if oldMapLen := len(_unifiedStationNameToStationIdMap) + len(station_blacklist.UnifiedStationNameToStationIdMap); len(newMap) != oldMapLen {
		return nil, nil,
			fmt.Errorf("seems like some stations got overridden by the black list [diff=%d]", len(newMap)-oldMapLen)
	}
	return _unifiedStationNameList, newMap, nil
}
