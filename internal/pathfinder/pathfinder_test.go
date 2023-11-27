package pathfinder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/internal/name"
	"zpcg/internal/transfer"
)

const (
	TimetableGobFilepath   = "timetable.gob"
	NiksicStationName      = "Nikšić"
	DanilovgradStationName = "Danilovgrad"
	BarStationName         = "Bar"
)

func TestFindDirectPaths(t *testing.T) {
	timetable, err := transfer.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	pathFinder := NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	paths := pathFinder.findDirectPaths(
		timetable.UnifiedStationNameToStationIdMap[name.Unify(NiksicStationName)],
		timetable.UnifiedStationNameToStationIdMap[name.Unify(DanilovgradStationName)])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
}

func TestFindPaths(t *testing.T) {
	timetable, err := transfer.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	pathFinder := NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	paths := pathFinder.findPathsWithTransfer(
		timetable.UnifiedStationNameToStationIdMap[name.Unify(NiksicStationName)],
		timetable.UnifiedStationNameToStationIdMap[name.Unify(BarStationName)])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
}
