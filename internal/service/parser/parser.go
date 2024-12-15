package parser

import (
	"fmt"
	"io"
	"sort"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/name"
	"github.com/ivanov-gv/zpcg/internal/service/parser/routes"
	"github.com/ivanov-gv/zpcg/internal/service/parser/stations"
)

const (
	BaseUrl        = "https://api.zpcg.me/api"
	RoutesApiUrl   = BaseUrl + "/routes/cumulative"
	StationsApiUrl = BaseUrl + "/stops"
)

func ParseTimetable(additionalRoutesHttpPaths ...string) (timetable.ExportFormat, error) {
	var (
		zpcgStopIdToStationsMap map[int]timetable.Station
		stationsTypesMap        map[timetable.StationTypeId]timetable.StationType
		detailedTimetableMap    map[timetable.TrainId]timetable.DetailedTimetable
	)
	// parse stations
	stationsResponse, err := retryablehttp.Get(StationsApiUrl)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("retryablehttp.Get[url='%s']: %w", StationsApiUrl, err)
	}
	zpcgStopIdToStationsMap, stationsTypesMap, err = stations.ParseStations(stationsResponse.Body)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("stations.ParseStations: %w", err)
	}

	// parser routes
	routesResponse, err := retryablehttp.Get(RoutesApiUrl)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("retryablehttp.Get [url='%s']: %w", RoutesApiUrl, err)
	}
	var additionalRoutesResponseBodies []io.ReadCloser
	for _, url := range additionalRoutesHttpPaths {
		additionalRoutesResponse, err := retryablehttp.Get(BaseUrl + url)
		if err != nil {
			return timetable.ExportFormat{}, fmt.Errorf("retryablehttp.Get [url='%s']: %w", url, err)
		}
		additionalRoutesResponseBodies = append(additionalRoutesResponseBodies, additionalRoutesResponse.Body)
	}
	detailedTimetableMap, err = routes.ParseRoutes(zpcgStopIdToStationsMap, routesResponse.Body, additionalRoutesResponseBodies...)
	if err != nil {
		return timetable.ExportFormat{}, fmt.Errorf("routes.ParseRoutes: %w", err)
	}

	// convert to transfer format
	transferFormat := MapTimetableToTransferFormat(detailedTimetableMap, lo.Values(zpcgStopIdToStationsMap))
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
	// add station types
	transferFormat.StationTypes = stationsTypesMap
	return transferFormat, nil
}

func MapTimetableToTransferFormat(routes map[timetable.TrainId]timetable.DetailedTimetable, allStations []timetable.Station) timetable.ExportFormat {
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
	// fill unifiedStationNameToStationIdMap
	var unifiedStationNameToStationIdMap = make(map[string]timetable.StationId)
	for _, route := range routes {
		for _, station := range route.Stops {
			stationName := stationIdToStationMap[station.Id].Name
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
