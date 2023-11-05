package parser

import (
	"zpcg/internal/model"

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

func ParseTimetable() (
	map[model.StationId]model.TrainIdSet,
	map[model.TrainId]model.StationIdToStationMap,
	error,
) {
	generalTimetableResponse, err := retryablehttp.Get(GeneralTimetablePageUrl)
	if err != nil {
		return nil, nil, errors.Wrap(err, "can not get general timetable page with retryablehttp.Get")
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "general_page.ParseGeneralTimetablePage")
	}

	detailedTimetableMap := make(map[model.TrainId]parser_model.DetailedTimetable, len(generalTimetableMap))
	// do not rewrite this loop with concurrency because zpcg.me do not have enough resources to handle all those requests
	// concurrency version is in the commit f5a2f983ce73fcc74f271d3bc4db51c2c56fe89f
	for trainId, generalTimetable := range generalTimetableMap {
		detailedTimetableFullLink := BaseUrl + generalTimetable.DetailedTimetableLink
		response, err := retryablehttp.Get(detailedTimetableFullLink)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "can not get route info with route id = %d, link = %s using retryablehttp.Get",
				trainId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(model.TrainId(trainId), detailedTimetableFullLink, response.Body)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "trainId = %d, link = %s detailed_page.ParseDetailedTimetablePage",
				trainId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetableMap[detailedTimetable.TrainId] = detailedTimetable
	}
	// map timetable to the maps needed for work
	stationIdToTrainIdSet, trainIdToStationMap := GetStationsAndTrainsMaps(detailedTimetableMap)
	return stationIdToTrainIdSet, trainIdToStationMap, nil
}

func GetStationsAndTrainsMaps(routes map[model.TrainId]parser_model.DetailedTimetable) (
	map[model.StationId]model.TrainIdSet,
	map[model.TrainId]model.StationIdToStationMap,
// map[model.TrainId]string, TODO
// map[model.StationId]string,
) {
	var (
		stationIdToTrainIdSetMap = make(map[model.StationId]model.TrainIdSet)
		trainIdToStationsMap     = make(map[model.TrainId]model.StationIdToStationMap, len(routes))
	)
	// fill maps
	for trainId, route := range routes {
		var routeStationsMap = make(model.StationIdToStationMap, len(route.Stations))
		for _, station := range route.Stations {
			// add station to the route
			routeStationsMap[station.Id] = station
			// add route to the station
			if trainIdSet, ok := stationIdToTrainIdSetMap[station.Id]; ok { // found
				trainIdSet[trainId] = struct{}{}
			} else { // not created yet
				trainIdSet = make(model.TrainIdSet)
				trainIdSet[trainId] = struct{}{}
				stationIdToTrainIdSetMap[station.Id] = trainIdSet
			}
		}
		trainIdToStationsMap[trainId] = routeStationsMap
	}
	return stationIdToTrainIdSetMap, trainIdToStationsMap
}
