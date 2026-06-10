package timetable_parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	gen_timetable "github.com/ivanov-gv/zpcg/gen/timetable"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
)

func newStationParser() *stationParser {
	return &stationParser{
		zpcgStopIdToStationsMap: nil,
		stationIdToStationMap:   nil,
		stationTypesMap:         nil,
	}
}

type stationParser struct {
	zpcgStopIdToStationsMap map[int]timetable.Station
	stationIdToStationMap   map[timetable.StationId]timetable.Station
	stationTypesMap         map[timetable.StationTypeId]timetable.StationType
}

type rawStationType struct {
	StopTypeID int    `json:"StopTypeID"`
	NameMe     string `json:"Name_me"`
	NameEn     string `json:"Name_en"`
}

type rawStation struct {
	StopID     int            `json:"StopID"`
	NameMe     string         `json:"Name_me"`
	NameEn     string         `json:"Name_en"`
	NameMeCyr  string         `json:"Name_me_cyr"`
	StopTypeID int            `json:"StopTypeID"`
	Latitude   float64        `json:"Latitude"`
	Longitude  float64        `json:"Longitude"`
	Local      int            `json:"local"`
	StopType   rawStationType `json:"stop_type"`
}

func (s *stationParser) parseStations() error {
	// get stations from api
	stationsResponse, err := retryablehttp.Get(StationsApiUrl)
	if err != nil {
		return fmt.Errorf("retryablehttp.Get[url='%s']: %w", StationsApiUrl, err)
	}
	rawStationsReader := stationsResponse.Body
	defer func() {
		err := rawStationsReader.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close raw zpcgStopIdToStationsMap reader")
		}
	}()
	rawStationsBytes, err := io.ReadAll(rawStationsReader)
	if err != nil {
		return fmt.Errorf("io.ReadAll: %w", err)
	}

	// unmarshal
	var rawStations []rawStation
	err = json.Unmarshal(rawStationsBytes, &rawStations)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	// fill stationTypesMap
	s.stationTypesMap = make(map[timetable.StationTypeId]timetable.StationType)
	for _, station := range rawStations {
		savedType, ok := s.stationTypesMap[timetable.StationTypeId(station.StopType.StopTypeID)]
		if ok && !equalStationTypes(station.StopType, savedType) {
			log.Warn().Any("savedType", savedType).Any("parsedType", station.StopType).
				Msg("got two station types with the same id but diff names")
		}
		s.stationTypesMap[timetable.StationTypeId(station.StopType.StopTypeID)] = mapToStationType(station.StopType)
	}

	// fill zpcgStopIdToStationsMap
	lastTakenStationId := lo.Max(lo.Values(gen_timetable.Timetable.UnifiedStationNameToStationIdMap))
	generateStationId := func(stationName string) timetable.StationId {
		if id, found := gen_timetable.Timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)]; found && id > 0 {
			return id
		}
		lastTakenStationId += 1
		return lastTakenStationId
	}
	s.zpcgStopIdToStationsMap = lo.SliceToMap(rawStations, func(item rawStation) (int, timetable.Station) {
		return item.StopID, timetable.Station{
			Id:         generateStationId(item.NameMe),
			ZpcgStopId: item.StopID,
			Type:       timetable.StationTypeId(item.StopType.StopTypeID),
			Name:       item.NameMe,
			NameEn:     item.NameEn,
			NameCyr:    item.NameMeCyr,
		}
	})

	// fill stationIdToStationMap
	allStations := lo.Values(s.zpcgStopIdToStationsMap)
	s.stationIdToStationMap = lo.SliceToMap(allStations, func(item timetable.Station) (timetable.StationId, timetable.Station) {
		return item.Id, item
	})

	return nil
}

func equalStationTypes(raw rawStationType, parsed timetable.StationType) bool {
	return mapToStationType(raw) == parsed
}

func mapToStationType(raw rawStationType) timetable.StationType {
	return timetable.StationType{
		Id:     timetable.StationTypeId(raw.StopTypeID),
		Name:   raw.NameMe,
		NameEn: raw.NameEn,
	}
}
