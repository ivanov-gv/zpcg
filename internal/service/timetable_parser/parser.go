package timetable_parser

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/config/timetable_parser_config"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/station_blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
	"github.com/ivanov-gv/zpcg/internal/service/timetable_parser/routes"
)

const (
	BaseUrl        = "https://api.zpcg.me/api"
	StationsApiUrl = BaseUrl + "/stops"
)

func New(config timetable_parser_config.Config) TimetableParser {

	return TimetableParser{
		timetableConfig: config.Timetable,

		parsedTimetable: parsedTimetable{
			zpcgStopIdToStationsMap:          make(map[int]timetable.Station),
			stationTypesMap:                  make(map[timetable.StationTypeId]timetable.StationType),
			unifiedStationNameToStationIdMap: make(map[string]timetable.StationId),
			unifiedStationNameList:           make([][]rune, 0),
			transferStationId:                0,
			seasonTimetables:                 make([]timetable.Season, 0, len(config.Timetable.Seasons)),
		},
	}
}

type parsedTimetable struct {
	zpcgStopIdToStationsMap          map[int]timetable.Station
	stationTypesMap                  map[timetable.StationTypeId]timetable.StationType
	unifiedStationNameToStationIdMap map[string]timetable.StationId
	unifiedStationNameList           [][]rune
	transferStationId                timetable.StationId
	seasonTimetables                 []timetable.Season
}

type TimetableParser struct {
	timetableConfig timetable_parser_config.Timetable

	parsedTimetable
}

func (t *TimetableParser) ParseTimetable() (timetable.ExportFormat, error) {
	// parse stations
	err := t.parseStations()
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("stations.ParseStations: %w", err)
	}

	// parse seasons
	for _, season := range t.timetableConfig.Seasons {
		seasonTimetable, err := t.parseSeason(season)
		if err != nil {
			return timetable.ExportFormat{}, fmt.Errorf("t.parseSeason [name='%s']: %w", season.Name, err)
		}
		t.seasonTimetables = append(t.seasonTimetables, seasonTimetable)
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
		StationTypes:                     t.stationTypesMap,
		TransferStationId:                t.transferStationId,
	}, nil
}

func url(start, end string, date time.Time) string {
	start = strings.ReplaceAll(start, " ", "+")
	end = strings.ReplaceAll(end, " ", "+")
	dateString := date.Format("2006-02-01")
	return fmt.Sprintf(BaseUrl+"/routes?start=%s&finish=%s&date=%s", start, end, dateString)
}

func (t *TimetableParser) parseSeason(season timetable_parser_config.Season) (timetable.Season, error) {
	var routesResponseBodies []io.ReadCloser
	for _, route := range t.timetableConfig.Routes {
		routesResponse, err := retryablehttp.Get(url(route.Start, route.Finish, season.FetchDate))
		if err != nil {
			return timetable.Season{}, fmt.Errorf("retryablehttp.Get [url='%s']: %w", url, err)
		}
		routesResponseBodies = append(routesResponseBodies, routesResponse.Body)
	}

	detailedTimetableMap, err := routes.ParseRoutes(t.zpcgStopIdToStationsMap, routesResponseBodies...)
	if err != nil {
		return timetable.Season{}, fmt.Errorf("routes.ParseRoutes: %w", err)
	}
	parsedSeason := MapTimetableToSeason(season, detailedTimetableMap, lo.Values(t.zpcgStopIdToStationsMap))

	// fill unifiedStationNameToStationIdMap
	for _, route := range detailedTimetableMap {
		for _, station := range route.Stops {
			stationName := parsedSeason.StationIdToStationMap[station.Id].Name
			if stationId, ok := t.unifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)]; ok && stationId != station.Id {
				return timetable.Season{}, fmt.Errorf("seems like station.id is not unique across seasons: "+
					"season='%s', stationName='%s', stationId='%d', savedStationId='%d'",
					season.Name, stationName, station.Id, stationId)
			}
			t.unifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)] = station.Id
		}
	}
	return parsedSeason, nil
}

func MapTimetableToSeason(season timetable_parser_config.Season, routes map[timetable.TrainId]timetable.DetailedTimetable, allStations []timetable.Station) timetable.Season {
	// fill stationIdToTrainIdSetMap
	var stationIdToTrainIdSetMap = make(map[timetable.StationId]timetable.TrainIdSet)
	for trainId, route := range routes {
		for _, station := range route.Stops {
			// add route to the station
			if trainIdSet, ok := stationIdToTrainIdSetMap[station.Id]; ok { // found
				trainIdSet[trainId] = struct{}{}
			} else { // not created yet
				trainIdSet = make(timetable.TrainIdSet)
				trainIdSet[trainId] = struct{}{}
				stationIdToTrainIdSetMap[station.Id] = trainIdSet
			}
		}
	}
	// fill trainIdToStationsMap
	var trainIdToStationsMap = make(map[timetable.TrainId]timetable.StationIdToStationMap, len(routes))
	for trainId, route := range routes {
		var routeStationsMap = make(timetable.StationIdToStationMap, len(route.Stops))
		for _, station := range route.Stops {
			// add station to the route
			routeStationsMap[station.Id] = station
		}
		trainIdToStationsMap[trainId] = routeStationsMap
	}
	// fill stationIdToStationMap
	var stationIdToStationMap = lo.SliceToMap(allStations, func(item timetable.Station) (timetable.StationId, timetable.Station) {
		return item.Id, item
	})
	// fill trainIdToTrainInfoMap
	var trainIdToTrainInfoMap = make(map[timetable.TrainId]timetable.TrainInfo, len(routes))
	for trainId, route := range routes {
		trainIdToTrainInfoMap[trainId] = timetable.TrainInfo{
			TrainId:      trainId,
			TimetableUrl: route.TimetableUrl,
		}
	}

	return timetable.Season{
		Name:  season.Name,
		Start: season.Start,
		End:   season.End,
		Warning: timetable.Warning{
			Be: season.Warning.Be,
			De: season.Warning.De,
			En: season.Warning.En,
			Hr: season.Warning.Hr,
			Ru: season.Warning.Ru,
			Sk: season.Warning.Sk,
			Sr: season.Warning.Sr,
			Tr: season.Warning.Tr,
			Uk: season.Warning.Uk,
		},
		StationIdToTrainIdSet: stationIdToTrainIdSetMap,
		TrainIdToStationMap:   trainIdToStationsMap,
		StationIdToStationMap: stationIdToStationMap,
		TrainIdToTrainInfoMap: trainIdToTrainInfoMap,
	}
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

func AddAliases(unifiedStationNameList [][]rune, unifiedStationNameToStationIdMap map[string]timetable.StationId) ([][]rune, map[string]timetable.StationId, error) {
	// add aliases to stations list
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
