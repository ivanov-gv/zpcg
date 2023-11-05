package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const TimetableGobFilepath = "../../resources/timetable.gob"

func TestNewApp(t *testing.T) {
	app, err := NewApp(TimetableGobFilepath)
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
	app, err := NewApp(TimetableGobFilepath)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	message := app.GenerateRoute(NiksicWrongStationName, DanilovgradWrongStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, DanilovgradStationName)
	message = app.GenerateRoute(NiksicWrongStationName, BarWronStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, BarStationName)
}
