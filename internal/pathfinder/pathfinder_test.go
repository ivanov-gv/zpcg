package pathfinder

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
	"zpcg/internal/transfer"
)

const (
	TimetableGobFilepath = "../../resources/timetable.gob"
)

func TestFindDirectPaths(t *testing.T) {
	param1, param2, param3, err := transfer.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	stationNameToIdMap := lo.Invert(param3)
	transferStationId := stationNameToIdMap["Podgorica"]
	pathFinder := NewPathFinder(param1, param2, param3, transferStationId)
	paths := pathFinder.FindDirectPaths(stationNameToIdMap["Nikšić"], stationNameToIdMap["Danilovgrad"])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
}

func TestFindPaths(t *testing.T) {
	param1, param2, param3, err := transfer.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	stationNameToIdMap := lo.Invert(param3)
	transferStationId := stationNameToIdMap["Podgorica"]
	pathFinder := NewPathFinder(param1, param2, param3, transferStationId)
	paths, withTransfer := pathFinder.FindPaths(stationNameToIdMap["Nikšić"], stationNameToIdMap["Bar"])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
	assert.True(t, withTransfer)
}
