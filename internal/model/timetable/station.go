package timetable

import (
	"time"

	"golang.org/x/text/language"
)

// StationId is an id for stations. (-inf, 0] - for black listed stations, (0, +inf) - for regular stations
type StationId int

type StationIdToStationMap map[StationId]Stop

type Stop struct {
	Id        StationId
	Arrival   time.Time
	Departure time.Time
}

type StationTypeId int

type StationType struct {
	Id     StationTypeId
	Name   string
	NameEn string
}

type Station struct {
	Id         StationId
	ZpcgStopId int
	Type       StationTypeId
	Name       string
	NameEn     string
	NameCyr    string
}

type BlackListedStation struct {
	Name                               string
	LanguageTagToCustomErrorMessageMap map[language.Tag]string
}
