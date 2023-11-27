package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/internal/app"
	"zpcg/resources"
)

const (
	timetableGobFilePath = "timetable.gob"
)

func TestApp(t *testing.T) {
	timetableReader, err := resources.FS.Open(timetableGobFilePath)
	assert.NoError(t, err)
	_app, err := app.NewApp(timetableReader)
	assert.NoError(t, err)
	message, _ := _app.GenerateRoute("niksic", "bar")
	t.Log("\n", message, "\n")
	assert.NotEmpty(t, message)
	numLines := strings.Count(message, "\n")
	assert.Greater(t, numLines, 2) // there is at least header and at least one route
}
