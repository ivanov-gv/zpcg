package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	timetable := ImportTimetable()
	// not nil
	assert.NotNil(t, timetable.StationIdToStationMap)
	assert.NotNil(t, timetable.TrainIdToTrainInfoMap)
	assert.NotNil(t, timetable.StationIdToTrainIdSet)
	assert.NotNil(t, timetable.TrainIdToStationMap)
	assert.NotNil(t, timetable.UnifiedStationNameList)
	assert.NotNil(t, timetable.UnifiedStationNameToStationIdMap)
	// not empty
	assert.NotEmpty(t, timetable.StationIdToStationMap)
	assert.NotEmpty(t, timetable.TrainIdToTrainInfoMap)
	assert.NotEmpty(t, timetable.StationIdToTrainIdSet)
	assert.NotEmpty(t, timetable.TrainIdToStationMap)
	assert.NotEmpty(t, timetable.UnifiedStationNameList)
	assert.NotEmpty(t, timetable.UnifiedStationNameToStationIdMap)
}
