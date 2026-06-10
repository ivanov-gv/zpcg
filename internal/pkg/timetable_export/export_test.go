package timetable_export

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	timetable := ImportTimetable()
	// not nil
	assert.NotNil(t, timetable.Seasons)
	assert.NotNil(t, timetable.StationIdToStationMap)
	assert.NotNil(t, timetable.UnifiedStationNameList)
	assert.NotNil(t, timetable.UnifiedStationNameToStationIdMap)
	assert.NotNil(t, timetable.StationTypes)
	for _, season := range timetable.Seasons {
		assert.NotNil(t, season.TrainIdToTrainInfoMap)
		assert.NotNil(t, season.StationIdToTrainIdSet)
		assert.NotNil(t, season.TrainIdToStationMap)
	}
	// not empty
	assert.NotEmpty(t, timetable.Seasons)
	assert.NotEmpty(t, timetable.StationIdToStationMap)
	assert.NotEmpty(t, timetable.UnifiedStationNameList)
	assert.NotEmpty(t, timetable.UnifiedStationNameToStationIdMap)
	assert.NotEmpty(t, timetable.StationTypes)
	for _, season := range timetable.Seasons {
		assert.NotEmpty(t, season.TrainIdToTrainInfoMap)
		assert.NotEmpty(t, season.StationIdToTrainIdSet)
		assert.NotEmpty(t, season.TrainIdToStationMap)
	}
}
