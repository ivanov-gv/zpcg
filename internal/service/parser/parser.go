package parser

import (
	"fmt"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"

	"zpcg/internal/model"
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

func ParseTimetable() (model.TimetableTransferFormat, error) {
	generalTimetableResponse, err := retryablehttp.Get(GeneralTimetablePageUrl)
	if err != nil {
		return model.TimetableTransferFormat{}, fmt.Errorf("can not get general timetable page with retryablehttp.Get: %w", err)
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return model.TimetableTransferFormat{}, fmt.Errorf("general_page.ParseGeneralTimetablePage: %w", err)
	}

	detailedTimetableMap := make(map[model.TrainId]parser_model.DetailedTimetable, len(generalTimetableMap))
	// do not rewrite this loop with concurrency because zpcg.me do not have enough resources to handle all those requests
	// concurrency version is in the commit f5a2f983ce73fcc74f271d3bc4db51c2c56fe89f
	for trainId, generalTimetable := range generalTimetableMap {
		detailedTimetableFullLink := BaseUrl + generalTimetable.DetailedTimetableLink
		response, err := retryablehttp.Get(detailedTimetableFullLink)
		if err != nil {
			return model.TimetableTransferFormat{}, fmt.Errorf("can not get route info with route id = %d, link = %s using retryablehttp.Get: %w",
				trainId, generalTimetable.DetailedTimetableLink, err)
		}
		detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(trainId, detailedTimetableFullLink, response.Body)
		if err != nil {
			return model.TimetableTransferFormat{}, fmt.Errorf("trainId = %d, link = %s detailed_page.ParseDetailedTimetablePage: %w",
				trainId, generalTimetable.DetailedTimetableLink, err)
		}
		detailedTimetableMap[detailedTimetable.TrainId] = detailedTimetable
	}
	// convert to transfer format
	transferFormat := MapTimetableToTransferFormat(detailedTimetableMap)
	// add blacklisted stations
	transferFormat, err = AddBlacklistedStations(transferFormat)
	if err != nil {
		return model.TimetableTransferFormat{}, fmt.Errorf("AddBlacklistedStations: %w", err)
	}
	// add aliases
	transferFormat, err = AddAliases(transferFormat)
	if err != nil {
		return model.TimetableTransferFormat{}, fmt.Errorf("AddAliases: %w", err)
	}
	return transferFormat, nil
}

func MapTimetableToTransferFormat(routes map[model.TrainId]parser_model.DetailedTimetable) model.TimetableTransferFormat {
	// fill stationIdToTrainIdSetMap
	var stationIdToTrainIdSetMap = make(map[model.StationId]model.TrainIdSet)
	for trainId, route := range routes {
		for _, station := range route.Stops {
			// add route to the station
			if trainIdSet, ok := stationIdToTrainIdSetMap[station.Id]; ok { // found
				trainIdSet[trainId] = struct{}{}
			} else { // not created yet
				trainIdSet = make(model.TrainIdSet)
				trainIdSet[trainId] = struct{}{}
				stationIdToTrainIdSetMap[station.Id] = trainIdSet
			}
		}
	}
	// fill trainIdToStationsMap
	var trainIdToStationsMap = make(map[model.TrainId]model.StationIdToStationMap, len(routes))
	for trainId, route := range routes {
		var routeStationsMap = make(model.StationIdToStationMap, len(route.Stops))
		for _, station := range route.Stops {
			// add station to the route
			routeStationsMap[station.Id] = station
		}
		trainIdToStationsMap[trainId] = routeStationsMap
	}
	// fill stationIdToStationMap
	var stationIdToStationMap = make(map[model.StationId]model.Station)
	stationIdToStationNameMap := detailed_page.GetStationIdToNameMap()
	for stationId, stationName := range stationIdToStationNameMap {
		stationIdToStationMap[stationId] = model.Station{
			Id:   stationId,
			Name: stationName,
		}
	}
	// fill trainIdToTrainInfoMap
	var trainIdToTrainInfoMap = make(map[model.TrainId]model.TrainInfo, len(routes))
	for trainId, route := range routes {
		trainIdToTrainInfoMap[trainId] = model.TrainInfo{
			TrainId:      trainId,
			TimetableUrl: route.TimetableUrl,
		}
	}
	// fill unifiedStationNameToStationIdMap
	var unifiedStationNameToStationIdMap = make(map[string]model.StationId)
	for _, route := range routes {
		for _, station := range route.Stops {
			stationName := detailed_page.GetStationIdToNameMap()[station.Id]
			unifiedStationNameToStationIdMap[name.Unify(stationName)] = station.Id
		}
	}
	// fill unifiedStationNameList
	var unifiedStationNameList []string
	unifiedStationNameList = lo.Keys(unifiedStationNameToStationIdMap)
	// get transfer station id
	var transferStationId = unifiedStationNameToStationIdMap[name.Unify(model.TransferStationName)]

	return model.TimetableTransferFormat{
		StationIdToTrainIdSet:            stationIdToTrainIdSetMap,
		TrainIdToStationMap:              trainIdToStationsMap,
		StationIdToStationMap:            stationIdToStationMap,
		TrainIdToTrainInfoMap:            trainIdToTrainInfoMap,
		UnifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		UnifiedStationNameList:           unifiedStationNameList,
		TransferStationId:                transferStationId,
	}
}

func AddBlacklistedStations(timetable model.TimetableTransferFormat) (model.TimetableTransferFormat, error) {
	// station name list - UnifiedStationNameList
	timetable.UnifiedStationNameList = append(timetable.UnifiedStationNameList, blacklist.UnifiedNames...)

	// map: station name -> station id - UnifiedStationNameToStationIdMap
	newMap := lo.Assign(timetable.UnifiedStationNameToStationIdMap,
		blacklist.UnifiedStationNameToStationIdMap)
	// ensure the map is not broken
	if oldMapLen := len(timetable.UnifiedStationNameToStationIdMap) + len(blacklist.UnifiedStationNameToStationIdMap); len(newMap) != oldMapLen {
		return model.TimetableTransferFormat{},
			fmt.Errorf("seems like some stations got overriden by the black list [diff=%d]", len(newMap)-oldMapLen)
	}
	timetable.UnifiedStationNameToStationIdMap = newMap
	return timetable, nil
}

func AddAliases(timetable model.TimetableTransferFormat) (model.TimetableTransferFormat, error) {
	// add aliases to stations list
	timetable.UnifiedStationNameList = append(timetable.UnifiedStationNameList, AliasesAsUnifiedStationNames...)
	// add mapping from aliases to station id
	for stationName, aliases := range AliasesOriginalUnifiedStationNameToUnifiedAliasesMap {
		stationId, ok := timetable.UnifiedStationNameToStationIdMap[stationName]
		if !ok { // station not found
			return model.TimetableTransferFormat{}, fmt.Errorf("can't find station to add aliases to: stationName = %s", stationName)
		}
		for _, _alias := range aliases {
			timetable.UnifiedStationNameToStationIdMap[_alias] = stationId
		}
	}
	return timetable, nil
}
