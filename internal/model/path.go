package model

import "time"

type Stop struct {
	StationId   StationId
	StationName string
	Arrival     time.Time
	Departure   time.Time
}

type Path struct {
	Departure, Arrival Station
}
