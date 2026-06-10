package timetable

import (
	"fmt"

	timetable_model "github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/pkg/timetable_export"
	"github.com/ivanov-gv/zpcg/internal/service/date"
)

func New(timetable timetable_model.ExportFormat, dateService *date.DateService) (*TimetableService, error) {
	err := timetable_export.VerifySeasons(timetable.Seasons)
	if err != nil {
		return nil, fmt.Errorf("timetable.VerifySeasons: %w", err)
	}
	return &TimetableService{
		timetable:   timetable,
		dateService: dateService,
	}, nil
}

type TimetableService struct {
	timetable   timetable_model.ExportFormat
	dateService *date.DateService
}

func (s *TimetableService) TransferStationId() timetable_model.StationId {
	return s.timetable.TransferStationId
}

func (s *TimetableService) Season() timetable_model.Season {
	for _, season := range s.timetable.Seasons {
		if season.IsInSeason(s.dateService.CurrentDateAsTime()) {
			return season
		}
	}
	return timetable_model.Season{}
}
