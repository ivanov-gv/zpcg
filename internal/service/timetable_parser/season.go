package timetable_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/config/timetable_parser_config"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/pkg/utils"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
)

func newSeasonParser(config timetable_parser_config.Timetable, parser *stationParser) *seasonParser {
	return &seasonParser{
		timetableConfig:                  config,
		seasonTimetables:                 make([]timetable.Season, 0, len(config.Seasons)),
		unifiedStationNameToStationIdMap: make(map[string]timetable.StationId),
		stationParser:                    parser,
	}
}

type seasonParser struct {
	timetableConfig timetable_parser_config.Timetable

	seasonTimetables                 []timetable.Season
	unifiedStationNameToStationIdMap map[string]timetable.StationId

	*stationParser
}

func (p *seasonParser) parseSeasons() error {
	// parse stations
	err := p.stationParser.parseStations()
	if err != nil {
		return fmt.Errorf("stations.ParseStations: %w", err)
	}

	// then seasons
	for _, season := range p.timetableConfig.Seasons {
		seasonTimetable, err := p.parseSeason(season)
		if err != nil {
			return fmt.Errorf("t.parseSeason [name='%s']: %w", season.Name, err)
		}
		p.seasonTimetables = append(p.seasonTimetables, seasonTimetable)
	}
	return nil
}

const (
	BaseUrl        = "https://api.zpcg.me/api"
	StationsApiUrl = BaseUrl + "/stops"
)

func buildUrl(start, end string, date time.Time) string {
	start = strings.ReplaceAll(start, " ", "+")
	end = strings.ReplaceAll(end, " ", "+")
	dateString := date.Format("2006-01-02")
	return fmt.Sprintf(BaseUrl+"/routes?start=%s&finish=%s&date=%s", start, end, dateString)
}

func (p *seasonParser) parseSeason(season timetable_parser_config.Season) (timetable.Season, error) {
	var routesResponseBodies []io.ReadCloser
	for _, route := range p.timetableConfig.Routes {
		url := buildUrl(route.Start, route.Finish, season.FetchDate)
		routesResponse, err := retryablehttp.Get(url)
		if err != nil {
			return timetable.Season{}, fmt.Errorf("retryablehttp.Get [url='%s']: %w", url, err)
		}
		routesResponseBodies = append(routesResponseBodies, routesResponse.Body)
	}

	detailedTimetableMap, err := parseRoutes(p.zpcgStopIdToStationsMap, routesResponseBodies...)
	if err != nil {
		return timetable.Season{}, fmt.Errorf("routes.ParseRoutes: %w", err)
	}
	parsedSeason := mapTimetableToSeason(season, detailedTimetableMap)

	// fill unifiedStationNameToStationIdMap
	for _, route := range detailedTimetableMap {
		for _, station := range route.Stops {
			stationName := p.stationIdToStationMap[station.Id].Name
			if stationId, ok := p.unifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)]; ok && stationId != station.Id {
				return timetable.Season{}, fmt.Errorf("seems like station.id is not unique across seasons: "+
					"season='%s', stationName='%s', stationId='%d', savedStationId='%d'",
					season.Name, stationName, station.Id, stationId)
			}
			p.unifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)] = station.Id
		}
	}
	return parsedSeason, nil
}

type rawTimetableItem struct {
	ArrivalTime   string `json:"ArrivalTime"`
	DepartureTime string `json:"DepartureTime"`
	Routestop     struct {
		StopID int `json:"StopID"`
	} `json:"routestop"`
}

type rawTrainInfo struct {
	TrainNumber    string             `json:"TrainNumber"`
	International  int                `json:"International"`
	TimetableItems []rawTimetableItem `json:"timetable_items"`
	Route          struct {
		ValidFrom string `json:"ValidFrom"`
		ValidTo   string `json:"ValidTo"`
	} `json:"route"`
}

type rawAvailableRoutes struct {
	Direct []rawTrainInfo `json:"direct"`
}

func parseRoutes(zpcgStopIdToStationsMap map[int]timetable.Station, routesBytes ...io.ReadCloser) (map[timetable.TrainId]timetable.DetailedTimetable, error) {
	defer func() {
		for _, bytes := range routesBytes {
			err := bytes.Close()
			if err != nil {
				log.Warn().Err(err).Msg("failed to close additionalRoutes bytes")
			}
		}
	}()

	// read trains info
	var allTrainInfos []rawTrainInfo
	for _, routeBytes := range routesBytes {
		bytes, err := io.ReadAll(routeBytes)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll: %w", err)
		}
		var route rawAvailableRoutes
		err = json.Unmarshal(bytes, &route)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		allTrainInfos = append(allTrainInfos, route.Direct...)
	}
	// build a set of trains
	var trainInfoMap = map[string]rawTrainInfo{}
	for _, trainInfo := range allTrainInfos {
		savedTrainInfo, found := trainInfoMap[trainInfo.TrainNumber]
		if found && len(savedTrainInfo.TimetableItems) > len(trainInfo.TimetableItems) {
			continue // the already saved one has more info about stations than the current one - skip
		}
		trainInfoMap[trainInfo.TrainNumber] = trainInfo
	}
	// build result
	var resultErr error
	var result = lo.MapEntries(trainInfoMap,
		func(key string, value rawTrainInfo) (timetable.TrainId, timetable.DetailedTimetable) {
			trainId, detailedTimetable, err := rawTrainInfoToDetailedTimetable(zpcgStopIdToStationsMap, value)
			if err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("strconv.Atoi [key='%s']: %w", key, err))
			}
			return trainId, detailedTimetable
		},
	)

	if resultErr != nil {
		return nil, resultErr
	}
	return result, nil
}

func rawTrainInfoToDetailedTimetable(zpcgStopIdToStationsMap map[int]timetable.Station, trainInfo rawTrainInfo) (timetable.TrainId, timetable.DetailedTimetable, error) {
	// valid from, valid to
	validFrom, err := parseValidity(trainInfo.Route.ValidFrom)
	if err != nil {
		return 0, timetable.DetailedTimetable{}, fmt.Errorf("parseValidity [validFrom]: %w", err)
	}
	validTo, err := parseValidity(trainInfo.Route.ValidTo)
	if err != nil {
		return 0, timetable.DetailedTimetable{}, fmt.Errorf("parseValidity [validTo]: %w", err)
	}
	// trainId
	trainId, err := strconv.Atoi(trainInfo.TrainNumber)
	if err != nil {
		return 0, timetable.DetailedTimetable{}, fmt.Errorf("strconv.Atoi: %w", err)
	}
	// stops
	stops, err := parseStops(zpcgStopIdToStationsMap, trainInfo.TimetableItems)
	if err != nil {
		return 0, timetable.DetailedTimetable{}, fmt.Errorf("parseStops: %w", err)
	}

	return timetable.TrainId(trainId), timetable.DetailedTimetable{
		TrainId:       timetable.TrainId(trainId),
		TimetableUrl:  "",
		International: trainInfo.International > 0,
		ValidFrom:     validFrom,
		ValidTo:       validTo,
		Stops:         stops,
	}, nil
}

const routeValidityDateLayout = "2006-01-02"

func parseValidity(validity string) (time.Time, error) {
	if len(validity) == 0 {
		return time.Time{}, fmt.Errorf("validity is empty")
	}
	parsed, err := time.Parse(routeValidityDateLayout, validity)
	if err != nil {
		return time.Time{}, fmt.Errorf("time.Parse [validity='%s']: %w", validity, err)
	}
	return parsed, nil
}

func parseStops(zpcgStopIdToStationsMap map[int]timetable.Station, timetableItems []rawTimetableItem) ([]timetable.Stop, error) {
	var result = make([]timetable.Stop, 0, len(timetableItems))
	for _, item := range timetableItems {
		// if the time is not present in the timetable - just get the last station departure time or empty one
		var fallbackTime time.Time
		if prevStop, found := lo.Last(result); found {
			fallbackTime = prevStop.Departure
		}
		// parse
		station, err := parseStop(zpcgStopIdToStationsMap, item, fallbackTime)
		if err != nil {
			return nil, fmt.Errorf("parseStop: %w", err)
		}

		// there might be a train with stops after midnight
		// in this case we need to add 24h to the arrival/departure time
		if prevStop, found := lo.Last(result); found &&
			(prevStop.Departure.After(utils.Midnight) || // previous stop departure is after midnight
				prevStop.Departure.After(station.Arrival) || // previous stop departure is before midnight, but current stop arrival is after midnight. like 23:58 -> 00:02
				station.Arrival.After(station.Departure)) { // arrival at a current stop is after departure: 23:58 -> 00:02
			if !station.Arrival.After(station.Departure) { // if both arrival and departure are after midnight - add 24h to both. example: 00:02 -> 00:05
				station.Arrival = station.Arrival.Add(utils.Day)
			}
			station.Departure = station.Departure.Add(utils.Day)
		}
		result = append(result, station)
	}
	return result, nil
}

const stopTimeLayout = "15:04:05"

func parseStop(zpcgStopIdToStationsMap map[int]timetable.Station, timetableItem rawTimetableItem, fallbackTime time.Time) (timetable.Stop, error) {
	var (
		arrivalTime, departureTime     time.Time
		arrivalFilled, departureFilled bool
	)
	// arrival
	if len(timetableItem.ArrivalTime) != 0 {
		_arrivalTime, err := time.Parse(stopTimeLayout, timetableItem.ArrivalTime)
		if err != nil {
			return timetable.Stop{}, fmt.Errorf("time.Parse: %w", err)
		}
		arrivalTime = _arrivalTime
		arrivalFilled = true
	}
	// departure
	if len(timetableItem.DepartureTime) != 0 {
		_departureTime, err := time.Parse(stopTimeLayout, timetableItem.DepartureTime)
		if err != nil {
			return timetable.Stop{}, fmt.Errorf("time.Parse: %w", err)
		}
		departureTime = _departureTime
		departureFilled = true
	}
	// fix missing times
	if !arrivalFilled && !departureFilled {
		arrivalTime = fallbackTime
		departureTime = fallbackTime
	}
	if !arrivalFilled {
		arrivalTime = departureTime
	}
	if !departureFilled {
		departureTime = arrivalTime
	}
	// for some reason all the departure and arrival times have seconds. it is not present in real timetable, so we should get rid of it.
	// it will help keep the parsing result consistent
	departureTime = departureTime.Truncate(time.Minute)
	arrivalTime = arrivalTime.Truncate(time.Minute)

	// station
	var stationId timetable.StationId
	if station, found := zpcgStopIdToStationsMap[timetableItem.Routestop.StopID]; found {
		stationId = station.Id
	} else {
		return timetable.Stop{}, fmt.Errorf("station not found: stopID='%d'", timetableItem.Routestop.StopID)
	}
	return timetable.Stop{
		Id:        stationId,
		Arrival:   arrivalTime,
		Departure: departureTime,
	}, nil
}

func mapTimetableToSeason(season timetable_parser_config.Season, routes map[timetable.TrainId]timetable.DetailedTimetable) timetable.Season {
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
		TrainIdToTrainInfoMap: trainIdToTrainInfoMap,
	}
}
