package parser

import (
	"github.com/samber/lo"
	"zpcg/internal/model"
	"zpcg/internal/name"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"

	"zpcg/internal/parser/detailed_page"
	"zpcg/internal/parser/general_page"
	parser_model "zpcg/internal/parser/model"
)

const (
	BaseUrl                 = "https://zpcg.me"
	GeneralTimetablePageUrl = "https://zpcg.me/search"
)

func ParseTimetable() (model.TimetableTransferFormat, error) {
	generalTimetableResponse, err := retryablehttp.Get(GeneralTimetablePageUrl)
	if err != nil {
		return model.TimetableTransferFormat{}, errors.Wrap(err, "can not get general timetable page with retryablehttp.Get")
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return model.TimetableTransferFormat{}, errors.Wrap(err, "general_page.ParseGeneralTimetablePage")
	}

	detailedTimetableMap := make(map[model.TrainId]parser_model.DetailedTimetable, len(generalTimetableMap))
	// do not rewrite this loop with concurrency because zpcg.me do not have enough resources to handle all those requests
	// concurrency version is in the commit f5a2f983ce73fcc74f271d3bc4db51c2c56fe89f
	for trainId, generalTimetable := range generalTimetableMap {
		detailedTimetableFullLink := BaseUrl + generalTimetable.DetailedTimetableLink
		response, err := retryablehttp.Get(detailedTimetableFullLink)
		if err != nil {
			return model.TimetableTransferFormat{}, errors.Wrapf(err, "can not get route info with route id = %d, link = %s using retryablehttp.Get",
				trainId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(model.TrainId(trainId), detailedTimetableFullLink, response.Body)
		if err != nil {
			return model.TimetableTransferFormat{}, errors.Wrapf(err, "trainId = %d, link = %s detailed_page.ParseDetailedTimetablePage",
				trainId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetableMap[detailedTimetable.TrainId] = detailedTimetable
	}
	// convert to transfer format
	transferFormat := MapTimetableToTransferFormat(detailedTimetableMap)
	return transferFormat, nil
}

func MapTimetableToTransferFormat(routes map[model.TrainId]parser_model.DetailedTimetable) model.TimetableTransferFormat {
	// fill stationIdToTrainIdSetMap
	var stationIdToTrainIdSetMap = make(map[model.StationId]model.TrainIdSet)
	for trainId, route := range routes {
		for _, station := range route.Stations {
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
		var routeStationsMap = make(model.StationIdToStationMap, len(route.Stations))
		for _, station := range route.Stations {
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
		for _, station := range route.Stations {
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
		StationIdToStaionMap:             stationIdToStationMap,
		TrainIdToTrainInfoMap:            trainIdToTrainInfoMap,
		UnifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		UnifiedStationNameList:           unifiedStationNameList,
		TransferStationId:                transferStationId,
	}
}
