package model

import "time"

type StationId int

type StationIdToStationMap map[StationId]Station

type Station struct {
	Id        StationId
	Arrival   time.Time
	Departure time.Time
}
