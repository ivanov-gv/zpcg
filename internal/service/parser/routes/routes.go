package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/utils"
)

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

func ParseRoutes(zpcgStopIdToStationsMap map[int]timetable.Station, rawRoutesBytes io.ReadCloser, additionalRoutesBytes ...io.ReadCloser) (map[timetable.TrainId]timetable.DetailedTimetable, error) {
	defer func() {
		err := rawRoutesBytes.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close raw routes bytes")
		}
		for _, bytes := range additionalRoutesBytes {
			err := bytes.Close()
			if err != nil {
				log.Warn().Err(err).Msg("failed to close additionalRoutes bytes")
			}
		}
	}()

	// read trains info
	bytes, err := io.ReadAll(rawRoutesBytes)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	var rawRoutes map[string]map[string]rawAvailableRoutes
	err = json.Unmarshal(bytes, &rawRoutes)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}
	// get all trains
	var allTrainInfos []rawTrainInfo
	for _, rawRouteMap := range rawRoutes {
		availableRoutes := lo.Values(rawRouteMap)
		trainInfos := lo.Flatten(lo.Map(availableRoutes, func(item rawAvailableRoutes, index int) []rawTrainInfo {
			return item.Direct
		}))
		allTrainInfos = append(allTrainInfos, trainInfos...)
	}
	// add additional trains
	for _, additionalRouteBytes := range additionalRoutesBytes {
		bytes, err := io.ReadAll(additionalRouteBytes)
		if err != nil {
			return nil, fmt.Errorf("io.ReadAll: %w", err)
		}
		var additionalRoute rawAvailableRoutes
		err = json.Unmarshal(bytes, &additionalRoute)
		if err != nil {
			return nil, fmt.Errorf("json.Unmarshal: %w", err)
		}
		allTrainInfos = append(allTrainInfos, additionalRoute.Direct...)
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
