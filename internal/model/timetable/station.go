package timetable

import (
	"time"
)

// StationId is an id for stations. (-inf, 0] - for black listed stations, (0, +inf) - for regular stations
type StationId int

type StationIdToStationMap map[StationId]Stop

type Stop struct {
	Id        StationId
	Arrival   time.Time
	Departure time.Time
}

type Station struct {
	Id   StationId
	Name string
}

type BlackListedStation struct {
	Name                               string
	LanguageTagToCustomErrorMessageMap map[string]string
}
