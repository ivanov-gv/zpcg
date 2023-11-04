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

func TestPathfinder(t *testing.T) {
	param1, param2, param3, err := transfer.ImportTimetable(TimetableGobFilepath)
	assert.NoError(t, err)
	pathFinder := NewPathFinder(param1, param2, param3)
	stationNameToIdMap := lo.Invert(param3)
	paths := pathFinder.FindStraightPaths(stationNameToIdMap["Nikšić"], stationNameToIdMap["Danilovgrad"])
	assert.NotNil(t, paths)
	assert.NotEmpty(t, paths)
	assert.Len(t, paths, 5)
}
