package timetable_export

import (
	"fmt"
	"time"

	"github.com/samber/lo"

	model_timetable "github.com/ivanov-gv/zpcg/internal/model/timetable"
)

type TimeRange struct {
	Name       string
	Start, End time.Time
}

func VerifySeasons(seasons []model_timetable.Season) error {
	timeRanges := lo.Map(seasons, func(item model_timetable.Season, _ int) TimeRange {
		return TimeRange{
			Name:  item.Name,
			Start: item.Start,
			End:   item.End,
		}
	})
	return VerifyTimeranges(timeRanges)
}

func VerifyTimeranges(timeRanges []TimeRange) error {
	if len(timeRanges) < 1 {
		return fmt.Errorf("at least one season must be added")
	}
	if !lo.FirstOrEmpty(timeRanges).Start.IsZero() {
		return fmt.Errorf("the first season's start date must be zero (i.e. -inf)")
	}
	if !lo.LastOrEmpty(timeRanges).End.IsZero() {
		return fmt.Errorf("the last season's end date must be zero (i.e. +inf)")
	}

	previousEnd := time.Time{}.Add(-24 * time.Hour)
	for i, season := range timeRanges {
		if !previousEnd.Add(24 * time.Hour).Equal(season.Start) {
			return fmt.Errorf("seasons must be consecutive. season '%s' must start '%s'",
				season.Name, previousEnd.Add(24*time.Hour).Format(time.DateOnly))
		}
		if i < len(timeRanges)-1 && season.Start.After(season.End) {
			return fmt.Errorf("'%s' season's end date must be after start date. currently start='%s', end='%s'",
				season.Name, season.Start.Format(time.DateOnly), season.End.Format(time.DateOnly))
		}
		previousEnd = season.End
	}
	return nil
}
