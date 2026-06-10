package pathfinder

import (
	"slices"

	timetable_model "github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/pkg/utils"
	"github.com/ivanov-gv/zpcg/internal/service/timetable"
)

func NewPathFinder(timetableService *timetable.TimetableService) *PathFinder {
	return &PathFinder{
		timetableService: timetableService,
	}
}

type PathFinder struct {
	timetableService *timetable.TimetableService
}

func (p *PathFinder) FindRoutes(aStation, bStation timetable_model.StationId) (routes []timetable_model.Path, isDirectRoute bool, err error) {
	season := p.timetableService.Season()
	// try to find direct routes
	if routes = p.findDirectPaths(season, aStation, bStation); len(routes) != 0 {
		return routes, true, nil
	}
	// try to find routes with transfer
	if routes = p.findPathsWithTransfer(season, aStation, bStation); len(routes) != 0 {
		return routes, false, nil
	}
	// no routes found
	return nil, false, ErrNoRoutesFound
}

func (p *PathFinder) findPathsWithTransfer(season timetable_model.Season, aStation, bStation timetable_model.StationId) []timetable_model.Path {
	// there is no direct trains from A to B. lets find one through the transfer station
	pathsAtoTransferStation := p.findDirectPaths(season, aStation, p.timetableService.TransferStationId())
	pathsTransferStationToB := p.findDirectPaths(season, p.timetableService.TransferStationId(), bStation)
	// merge paths
	var (
		paths                                      []timetable_model.Path
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

func (p *PathFinder) findDirectPaths(season timetable_model.Season, aStation, bStation timetable_model.StationId) []timetable_model.Path {
	var trainIdSetA, trainIdSetB timetable_model.TrainIdSet
	trainIdSetA = season.StationIdToTrainIdSet[aStation]
	trainIdSetB = season.StationIdToTrainIdSet[bStation]
	// get intersection of maps of the trains
	possibleRoutes := utils.Intersection(trainIdSetA, trainIdSetB)
	if len(possibleRoutes) == 0 { // I've got no routes~~~
		return nil
	}
	// find suitable routes
	paths := make([]timetable_model.Path, 0, len(possibleRoutes))
	for trainId := range possibleRoutes {
		stations := season.TrainIdToStationMap[trainId]
		origin := stations[aStation]
		destination := stations[bStation]

		if origin.Departure.After(destination.Arrival) { // <- this if is the reason we need to add 24h to the time and then normalize it back
			// Departure is coming after arrival - it is a train in a wrong direction. skip
			continue
		}
		// add found path
		paths = append(paths, timetable_model.Path{
			TrainId:     trainId,
			Origin:      origin,
			Destination: destination,
		})
	}
	// normalize time - make it all in a range 00:00 - 23:59
	utils.NormalizeTimeInPaths(paths)
	slices.SortFunc(paths, func(a, b timetable_model.Path) int {
		return a.Origin.Departure.Compare(b.Origin.Departure)
	})
	return paths
}
