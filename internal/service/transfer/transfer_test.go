package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/test/resources"
)

func TestImport(t *testing.T) {
	f, err := resources.TestFS.Open(resources.TestTimetableGobFilepath)
	assert.NoError(t, err)
	timetable, err := ImportTimetableFromReader(f)
	assert.NoError(t, err)
	// not nil
	assert.NotNil(t, timetable.StationIdToStaionMap)
	assert.NotNil(t, timetable.TrainIdToTrainInfoMap)
	assert.NotNil(t, timetable.StationIdToTrainIdSet)
	assert.NotNil(t, timetable.TrainIdToStationMap)
	assert.NotNil(t, timetable.UnifiedStationNameList)
	assert.NotNil(t, timetable.UnifiedStationNameToStationIdMap)
	// not empty
	assert.NotEmpty(t, timetable.StationIdToStaionMap)
	assert.NotEmpty(t, timetable.TrainIdToTrainInfoMap)
	assert.NotEmpty(t, timetable.StationIdToTrainIdSet)
	assert.NotEmpty(t, timetable.TrainIdToStationMap)
	assert.NotEmpty(t, timetable.UnifiedStationNameList)
	assert.NotEmpty(t, timetable.UnifiedStationNameToStationIdMap)
}
