package app

import (
	"fmt"
	"log"
	"zpcg/internal/model"
	"zpcg/internal/name"
	"zpcg/internal/render"

	"zpcg/internal/pathfinder"
	"zpcg/internal/transfer"
)

func NewApp(timetableFilepath string) (*App, error) {
	timetable, err := transfer.ImportTimetable(timetableFilepath)
	if err != nil {
		return nil, fmt.Errorf("transfer.ImportTimetable: %w", err)
	}
	finder := pathfinder.NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	stationNameResolver := name.NewStationNameResolver(timetable.UnifiedStationNameToStationIdMap, timetable.UnifiedStationNameList)
	_render := render.NewRender(timetable.StationIdToStaionMap, timetable.TrainIdToTrainInfoMap)
	return &App{
		finder:              finder,
		stationNameResolver: stationNameResolver,
		render:              _render,
		transferStationId:   timetable.TransferStationId,
	}, nil
}

type App struct {
	finder              *pathfinder.PathFinder
	stationNameResolver *name.StationNameResolver
	render              *render.Render
	transferStationId   model.StationId
}

func (a *App) GenerateRoute(origin, destination string) string {
	// find station ids
	originStationId, err := a.stationNameResolver.FindStationIdByApproximateName(origin)
	if err != nil {
		log.Println("err", err, "origin", origin)
		return a.render.StationResolveError(origin)
	}
	destinationStationId, err := a.stationNameResolver.FindStationIdByApproximateName(destination)
	if err != nil {
		log.Println("err", err, "destination", destination)
		return a.render.StationResolveError(destination)
	}
	// find route
	routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId)
	if isDirect {
		return a.render.DirectRoutes(routes)
	}
	return a.render.TransferRoutes(routes, originStationId, a.transferStationId, destinationStationId)
}
