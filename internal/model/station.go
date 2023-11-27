package model

import (
	"time"
)

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
