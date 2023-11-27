package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/resources"
)

const TimetableGobFilepath = "timetable.gob"

func TestNewApp(t *testing.T) {
	timetableReader, err := resources.FS.Open(TimetableGobFilepath)
	assert.NoError(t, err)
	app, err := NewApp(timetableReader)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}

const (
	NiksicStationName           = "Nikšić"
	NiksicWrongStationName      = "niksic"
	DanilovgradStationName      = "Danilovgrad"
	DanilovgradWrongStationName = "DaNil ovgrad"
	BarStationName              = "Bar"
	BarWronStationName          = "Barrrrrr"
)

func TestGenerateRoute(t *testing.T) {
	timetableReader, err := resources.FS.Open(TimetableGobFilepath)
	assert.NoError(t, err)
	app, err := NewApp(timetableReader)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	message, _ := app.GenerateRoute(NiksicWrongStationName, DanilovgradWrongStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, DanilovgradStationName)
	message, _ = app.GenerateRoute(NiksicWrongStationName, BarWronStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, BarStationName)
}
