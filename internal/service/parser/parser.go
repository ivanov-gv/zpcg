package parser

import (
	"fmt"
	"slices"
	"sort"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"

	"zpcg/internal/model/timetable"
	"zpcg/internal/service/blacklist"
	"zpcg/internal/service/name"
	"zpcg/internal/service/parser/detailed_page"
	"zpcg/internal/service/parser/general_page"
	parser_model "zpcg/internal/service/parser/model"
)

const (
	BaseUrl                 = "https://zpcg.me"
	GeneralTimetablePageUrl = "https://zpcg.me/search"
)

func ParseTimetable() (timetable.ExportFormat, error) {
	generalTimetableResponse, err := retryablehttp.Get(GeneralTimetablePageUrl)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("can not get general timetable page with retryablehttp.Get: %w", err)
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("general_page.ParseGeneralTimetablePage: %w", err)
	}

	detailedTimetableMap := make(map[timetable.TrainId]parser_model.DetailedTimetable, len(generalTimetableMap))
	// do not rewrite this loop with concurrency because zpcg.me do not have enough resources to handle all those requests
	// concurrency version is in the commit f5a2f983ce73fcc74f271d3bc4db51c2c56fe89f
	trainIds := lo.Keys(generalTimetableMap)
	slices.Sort(trainIds) // sort to make output stable
	for _, trainId := range trainIds {
		generalTimetable := generalTimetableMap[trainId]
		detailedTimetableFullLink := BaseUrl + generalTimetable.DetailedTimetableLink
		response, err := retryablehttp.Get(detailedTimetableFullLink)
		if err != nil {
			return timetable.ExportFormat{}, fmt.Errorf("can not get route info with route id = %d, link = %s using retryablehttp.Get: %w",
				trainId, generalTimetable.DetailedTimetableLink, err)
		}
		detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(trainId, detailedTimetableFullLink, response.Body)
		if err != nil {
			return timetable.ExportFormat{}, fmt.Errorf("trainId = %d, link = %s detailed_page.ParseDetailedTimetablePage: %w",
				trainId, generalTimetable.DetailedTimetableLink, err)
		}
		detailedTimetableMap[detailedTimetable.TrainId] = detailedTimetable
	}
	// convert to transfer format
	transferFormat := MapTimetableToTransferFormat(detailedTimetableMap)
	// add blacklisted stations
	transferFormat, err = AddBlacklistedStations(transferFormat)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("AddBlacklistedStations: %w", err)
	}
	// add aliases
	transferFormat, err = AddAliases(transferFormat)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("AddAliases: %w", err)
	}
	return transferFormat, nil
}

func MapTimetableToTransferFormat(routes map[timetable.TrainId]parser_model.DetailedTimetable) timetable.ExportFormat {
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
	var stationIdToStationMap = make(map[timetable.StationId]timetable.Station)
	stationIdToStationNameMap := detailed_page.GetStationIdToNameMap()
	for stationId, stationName := range stationIdToStationNameMap {
		stationIdToStationMap[stationId] = timetable.Station{
			Id:   stationId,
			Name: stationName,
		}
	}
	// fill trainIdToTrainInfoMap
	var trainIdToTrainInfoMap = make(map[timetable.TrainId]timetable.TrainInfo, len(routes))
	for trainId, route := range routes {
		trainIdToTrainInfoMap[trainId] = timetable.TrainInfo{
			TrainId:      trainId,
			TimetableUrl: route.TimetableUrl,
		}
	}
	// fill unifiedStationNameToStationIdMap
	var unifiedStationNameToStationIdMap = make(map[string]timetable.StationId)
	for _, route := range routes {
		for _, station := range route.Stops {
			stationName := detailed_page.GetStationIdToNameMap()[station.Id]
			unifiedStationNameToStationIdMap[name.Unify(stationName)] = station.Id
		}
	}
	// fill unifiedStationNameList
	var unifiedStationNameList []string
	unifiedStationNameList = lo.Keys(unifiedStationNameToStationIdMap)
	// sort for more predictable output
	sort.Strings(unifiedStationNameList)
	// get transfer station id
	var transferStationId = unifiedStationNameToStationIdMap[name.Unify(timetable.TransferStationName)]

	return timetable.ExportFormat{
		StationIdToTrainIdSet:            stationIdToTrainIdSetMap,
		TrainIdToStationMap:              trainIdToStationsMap,
		StationIdToStationMap:            stationIdToStationMap,
		TrainIdToTrainInfoMap:            trainIdToTrainInfoMap,
		UnifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		UnifiedStationNameList:           lo.Map(unifiedStationNameList, func(item string, _ int) []rune { return []rune(item) }),
		TransferStationId:                transferStationId,
	}
}

func AddBlacklistedStations(_timetable timetable.ExportFormat) (timetable.ExportFormat, error) {
	// station name list - UnifiedStationNameList
	_timetable.UnifiedStationNameList = append(_timetable.UnifiedStationNameList,
		lo.Map(blacklist.UnifiedNames, func(item string, _ int) []rune { return []rune(item) })...)

	// map: station name -> station id - UnifiedStationNameToStationIdMap
	newMap := lo.Assign(_timetable.UnifiedStationNameToStationIdMap,
		blacklist.UnifiedStationNameToStationIdMap)
	// ensure the map is not broken
	if oldMapLen := len(_timetable.UnifiedStationNameToStationIdMap) + len(blacklist.UnifiedStationNameToStationIdMap); len(newMap) != oldMapLen {
		return timetable.ExportFormat{},
			fmt.Errorf("seems like some stations got overriden by the black list [diff=%d]", len(newMap)-oldMapLen)
	}
	_timetable.UnifiedStationNameToStationIdMap = newMap
	return _timetable, nil
}

func AddAliases(_timetable timetable.ExportFormat) (timetable.ExportFormat, error) {
	// add aliases to stations list
	_timetable.UnifiedStationNameList = append(_timetable.UnifiedStationNameList,
		lo.Map(AliasesAsUnifiedStationNames, func(item string, _ int) []rune { return []rune(item) })...)
	// add mapping from aliases to station id
	for stationName, aliases := range AliasesOriginalUnifiedStationNameToUnifiedAliasesMap {
		stationId, ok := _timetable.UnifiedStationNameToStationIdMap[stationName]
		if !ok { // station not found
			return timetable.ExportFormat{}, fmt.Errorf("can't find station to add aliases to: stationName = %s", stationName)
		}
		for _, _alias := range aliases {
			_timetable.UnifiedStationNameToStationIdMap[_alias] = stationId
		}
	}
	return _timetable, nil
}
