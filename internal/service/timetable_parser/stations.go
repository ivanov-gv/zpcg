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

func (t *TimetableParser) parseStations() error {
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

	var rawStations []rawStation
	err = json.Unmarshal(rawStationsBytes, &rawStations)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	// parse station types
	for _, station := range rawStations {
		savedType, ok := t.stationTypesMap[timetable.StationTypeId(station.StopType.StopTypeID)]
		if ok && !equalStationTypes(station.StopType, savedType) {
			log.Warn().Any("savedType", savedType).Any("parsedType", station.StopType).
				Msg("got two station types with the same id but diff names")
		}
		t.stationTypesMap[timetable.StationTypeId(station.StopType.StopTypeID)] = mapToStationType(station.StopType)
	}

	// parse zpcgStopIdToStationsMap
	lastTakenStationId := lo.Max(lo.Values(gen_timetable.Timetable.UnifiedStationNameToStationIdMap))
	generateStationId := func(stationName string) timetable.StationId {
		if id, found := gen_timetable.Timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(stationName)]; found && id > 0 {
			return id
		}
		lastTakenStationId += 1
		return lastTakenStationId
	}
	t.zpcgStopIdToStationsMap = lo.SliceToMap(rawStations, func(item rawStation) (int, timetable.Station) {
		return item.StopID, timetable.Station{
			Id:         generateStationId(item.NameMe),
			ZpcgStopId: item.StopID,
			Type:       timetable.StationTypeId(item.StopType.StopTypeID),
			Name:       item.NameMe,
			NameEn:     item.NameEn,
			NameCyr:    item.NameMeCyr,
		}
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
