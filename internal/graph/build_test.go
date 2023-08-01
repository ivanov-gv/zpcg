package graph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yourbasic/graph"
	"testing"
	"zpcg/internal/parser"
	"zpcg/internal/utils"
)

const (
	TimetableGobFilepath         = "../../resources/timetable.gob"
	NumOfStationsFromNiksicToBar = 14
)

func TestBuildRouteGraph(t *testing.T) {
	routes, err := parser.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	fullGraph, nameToIdMap := BuildRouteGraph(routes)
	assert.NotNil(t, fullGraph)
	assert.NotEmpty(t, nameToIdMap)
	for _, route := range routes {
		for _, station := range route.Stations {
			assert.Contains(t, nameToIdMap, station.Name)
		}
	}
	assert.LessOrEqual(t, len(nameToIdMap), MaxNumberOfStations)
	assert.Truef(t, graph.Connected(fullGraph), "graph.Connected")
	barStationId := nameToIdMap["Bar"]
	niksicStationId := nameToIdMap["Nikšić"]
	path, _ := graph.ShortestPath(fullGraph, niksicStationId, barStationId)
	assert.NotEmptyf(t, path, "graph.ShortestPath(fullGraph, niksicStationId, barStationId)")
	assert.Equalf(t, NumOfStationsFromNiksicToBar, len(path), "check num of stations in path")

	stationIdToName := utils.RevertMap(nameToIdMap)
	var pathString string
	for _, stationId := range path {
		pathString = fmt.Sprintf("%s -> %s", pathString, stationIdToName[stationId])
	}
	t.Logf("the path from Niksic to Bar: %s", pathString)
}
