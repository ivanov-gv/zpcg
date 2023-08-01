package graph

import (
	"github.com/yourbasic/graph"
	"zpcg/internal/parser/model"
)

func BuildRouteGraph(routes map[int]model.DetailedTimetable) (*graph.Mutable, map[string]int) {
	stationsNameToIdMap, numOfStations := GetStationsFromRoutes(routes)
	fullGraph := graph.New(numOfStations)
	for _, route := range routes {
		for i, currentStation := range route.Stations {
			if i == 0 {
				continue
			}
			// 0 < i < len(route.Stations)
			previousStation := route.Stations[i-1]
			previousStationId := stationsNameToIdMap[previousStation.Name]
			currentStationId := stationsNameToIdMap[currentStation.Name]
			fullGraph.Add(previousStationId, currentStationId)
		}
	}
	return fullGraph, stationsNameToIdMap
}
