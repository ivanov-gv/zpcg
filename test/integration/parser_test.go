package integration

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"zpcg/internal/parser"
)

func TestParser(t *testing.T) {
	timetable, err := parser.ParseTimetable()
	assert.NoError(t, err)
	assert.NotNil(t, timetable)

	for routeId, timetable := range timetable {
		assert.NotEmptyf(t, timetable.Stations, "there is no parsed stations for route %d", routeId)
		t.Logf("route: %d , stations: %v", routeId, timetable.Stations)
	}
}
