package timetable

import (
	"time"
)

type Season struct {
	Name    string
	Warning Warning

	Start time.Time // zero value is -inf
	End   time.Time // zero value is +inf

	StationIdToTrainIdSet map[StationId]TrainIdSet
	TrainIdToStationMap   map[TrainId]StationIdToStationMap
	StationIdToStationMap map[StationId]Station
	TrainIdToTrainInfoMap map[TrainId]TrainInfo
}

type Warning struct {
	Be string
	De string
	En string
	Hr string
	Ru string
	Sk string
	Sr string
	Tr string
	Uk string
}

func (s Season) IsInSeason(t time.Time) bool {
	afterStart := s.Start.IsZero() || !t.Before(s.Start)
	beforeEnd := s.End.IsZero() || !t.After(s.End)
	return afterStart && beforeEnd
}

func (s Season) Before(season Season) bool {
	return !s.End.IsZero() && s.End.Before(season.Start)
}

func (s Season) Intersects(season Season) bool {
	// s ends strictly before season starts
	if !s.End.IsZero() && !season.Start.IsZero() && s.End.Before(season.Start) {
		return false
	}
	// season ends strictly before s starts
	if !season.End.IsZero() && !s.Start.IsZero() && season.End.Before(s.Start) {
		return false
	}
	return true
}
