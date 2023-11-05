package transfer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const filepath = "./../../resources/timetable.gob"

func TestImport(t *testing.T) {
	timetable, err := ImportTimetable(filepath)
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
