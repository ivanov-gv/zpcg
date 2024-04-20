package utils

import (
	"time"

	"github.com/samber/lo"

	"zpcg/internal/model/timetable"
)

var (
	Midnight = time.Date(0, 1, 1, 23, 59, 59, 0, time.UTC)
	Day      = time.Hour * 24
)

func NormalizeTime(t time.Time) time.Time {
	// some arrivals/destinations are after the midnight for route path algo simplicity reasons
	// but in order to get consistent timetable with right sort
	// we need to extract 24h from the time in this case
	if t.After(Midnight) {
		return t.Add(-Day)
	}
	return t
}

// NormalizeTimeInPaths makes all the times in a range 00:00 - 23:59
func NormalizeTimeInPaths(paths []timetable.Path) {
	lo.ForEach(paths, func(item timetable.Path, index int) {
		item.Origin.Arrival = NormalizeTime(item.Origin.Arrival)
		item.Origin.Departure = NormalizeTime(item.Origin.Departure)
		item.Destination.Arrival = NormalizeTime(item.Destination.Arrival)
		item.Destination.Departure = NormalizeTime(item.Destination.Departure)
		paths[index] = item
	})
}
