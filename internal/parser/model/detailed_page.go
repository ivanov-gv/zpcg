package model

import "time"

type DetailedTimetable struct {
	RouteId  int
	Stations []Station
}

type Station struct {
	Name      string
	Arrival   time.Time
	Departure time.Time
}
