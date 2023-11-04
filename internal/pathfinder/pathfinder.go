package pathfinder

import (
	"slices"
	"zpcg/internal/model"
	"zpcg/internal/utils"
)

func NewPathFinder(
	stationIdToTrainIdSetMap map[model.StationId]model.TrainIdSet,
	trainIdToStationsMap map[model.TrainId]model.StationIdToStationMap,
	stationIdToStationName map[model.StationId]string) *PathFinder {
	return &PathFinder{
		stationIdToTrainIdSetMap: stationIdToTrainIdSetMap,
		trainIdToStationsMap:     trainIdToStationsMap,
		stationIdToStationName:   stationIdToStationName,
	}
}

type PathFinder struct {
	stationIdToTrainIdSetMap map[model.StationId]model.TrainIdSet
	trainIdToStationsMap     map[model.TrainId]model.StationIdToStationMap
	stationIdToStationName   map[model.StationId]string
}

func (p *PathFinder) FindPaths(aStation, bStation model.StationId) []model.Path {
	panic("not implemented")
}

func (p *PathFinder) FindStraightPaths(aStation, bStation model.StationId) []model.Path {
	var trainIdSetA, trainIdSetB model.TrainIdSet
	trainIdSetA = p.stationIdToTrainIdSetMap[aStation]
	trainIdSetB = p.stationIdToTrainIdSetMap[bStation]
	// get intersection of maps of the trains
	possibleRoutes := utils.Intersection(trainIdSetA, trainIdSetB)
	// find suitable routes
	paths := make([]model.Path, 0, len(possibleRoutes))
	for trainId := range possibleRoutes {
		stations := p.trainIdToStationsMap[trainId]
		departure := stations[aStation]
		arrival := stations[bStation]
		if departure.Departure.Before(arrival.Arrival) { // if it is a route in a right direction - add
			paths = append(paths, model.Path{
				Departure: departure,
				Arrival:   arrival,
			})
		}
	}
	slices.SortFunc(paths, func(a, b model.Path) int {
		return a.Departure.Departure.Compare(b.Departure.Departure)
	})
	return paths
}
