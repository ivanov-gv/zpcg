package integration

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"zpcg/internal/app"
)

const (
	timetableGobFilePath = "../../resources/timetable.gob"
)

func TestApp(t *testing.T) {
	_app, err := app.NewApp(timetableGobFilePath)
	assert.NoError(t, err)
	message, _ := _app.GenerateRoute("niksic", "bar")
	t.Log("\n", message, "\n")
	assert.NotEmpty(t, message)
	numLines := strings.Count(message, "\n")
	assert.Greater(t, numLines, 2) // there is at least header and at least one route
}
