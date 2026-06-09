package pathfinder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
	"github.com/ivanov-gv/zpcg/internal/service/timetable_export"
)

const (
	NiksicStationName      = "Nikšić"
	DanilovgradStationName = "Danilovgrad"
	BarStationName         = "Bar"
)

func TestFindDirectPaths(t *testing.T) {
	timetable := timetable_export.ImportTimetable()
	pathFinder := NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	paths := pathFinder.findDirectPaths(
		timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(NiksicStationName)],
		timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(DanilovgradStationName)])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
}

func TestFindPaths(t *testing.T) {
	timetable := timetable_export.ImportTimetable()
	pathFinder := NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	paths := pathFinder.findPathsWithTransfer(
		timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(NiksicStationName)],
		timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(BarStationName)])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
}
