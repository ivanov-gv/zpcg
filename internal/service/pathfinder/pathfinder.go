package pathfinder

import (
	"slices"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/utils"
)

func NewPathFinder(
	stationIdToTrainIdSetMap map[timetable.StationId]timetable.TrainIdSet,
	trainIdToStationsMap map[timetable.TrainId]timetable.StationIdToStationMap,
	transferStation timetable.StationId) *PathFinder {
	return &PathFinder{
		stationIdToTrainIdSetMap: stationIdToTrainIdSetMap,
		trainIdToStationsMap:     trainIdToStationsMap,
		transferStation:          transferStation,
	}
}

type PathFinder struct {
	stationIdToTrainIdSetMap map[timetable.StationId]timetable.TrainIdSet
	trainIdToStationsMap     map[timetable.TrainId]timetable.StationIdToStationMap
	transferStation          timetable.StationId
}

func (p *PathFinder) FindRoutes(aStation, bStation timetable.StationId) (routes []timetable.Path, isDirectRoute bool) {
	var (
		paths    []timetable.Path
		isDirect bool
	)
	// try to find direct paths
	if paths = p.findDirectPaths(aStation, bStation); len(paths) != 0 {
		isDirect = true
	} else {
		// find paths with transfer
		paths = p.findPathsWithTransfer(aStation, bStation)
		isDirect = false
	}
	return paths, isDirect
}

func (p *PathFinder) findPathsWithTransfer(aStation, bStation timetable.StationId) []timetable.Path {
	// there is no direct trains from A to B. lets find one through the transfer station
	pathsAtoTransferStation := p.findDirectPaths(aStation, p.transferStation)
	pathsTransferStationToB := p.findDirectPaths(p.transferStation, bStation)
	// merge paths
	var (
		paths                                      []timetable.Path
		indexPathAtoTransfer, indexPathTransferToB int
	)
	// run through all the paths and add them to the result
	for indexPathAtoTransfer < len(pathsAtoTransferStation) && indexPathTransferToB < len(pathsTransferStationToB) {
		currentPathAtoTransfer := pathsAtoTransferStation[indexPathAtoTransfer]
		currentPathTransferToB := pathsTransferStationToB[indexPathTransferToB]
		if currentPathAtoTransfer.Destination.Arrival.Before(currentPathTransferToB.Origin.Departure) {
			// arrival in the transfer station is before departure of the next train - from the transfer to destination
			// which means this could be a complete route from A to B
			// append currentPathAtoTransfer then
			paths = append(paths, currentPathAtoTransfer)
			indexPathAtoTransfer++
		} else {
			paths = append(paths, currentPathTransferToB)
			indexPathTransferToB++
		}
	}
	// add remaining paths
	if indexPathAtoTransfer < len(pathsAtoTransferStation) {
		paths = append(paths, pathsAtoTransferStation[indexPathAtoTransfer:]...)
	}
	if indexPathTransferToB < len(pathsTransferStationToB) {
		paths = append(paths, pathsTransferStationToB[indexPathTransferToB:]...)
	}
	return paths
}

func (p *PathFinder) findDirectPaths(aStation, bStation timetable.StationId) []timetable.Path {
	var trainIdSetA, trainIdSetB timetable.TrainIdSet
	trainIdSetA = p.stationIdToTrainIdSetMap[aStation]
	trainIdSetB = p.stationIdToTrainIdSetMap[bStation]
	// get intersection of maps of the trains
	possibleRoutes := utils.Intersection(trainIdSetA, trainIdSetB)
	if len(possibleRoutes) == 0 { // I've got no routes~~~
		return nil
	}
	// find suitable routes
	paths := make([]timetable.Path, 0, len(possibleRoutes))
	for trainId := range possibleRoutes {
		stations := p.trainIdToStationsMap[trainId]
		origin := stations[aStation]
		destination := stations[bStation]

		if origin.Departure.After(destination.Arrival) { // <- this if is the reason we need to add 24h to the time and then normalize it back
			// Departure is coming after arrival - it is a train in a wrong direction. skip
			continue
		}
		// add found path
		paths = append(paths, timetable.Path{
			TrainId:     trainId,
			Origin:      origin,
			Destination: destination,
		})
	}
	// normalize time - make it all in a range 00:00 - 23:59
	utils.NormalizeTimeInPaths(paths)
	slices.SortFunc(paths, func(a, b timetable.Path) int {
		return a.Origin.Departure.Compare(b.Origin.Departure)
	})
	return paths
}
