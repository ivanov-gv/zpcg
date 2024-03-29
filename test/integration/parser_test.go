package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/internal/app"
	"zpcg/internal/config"
	"zpcg/internal/service/render"
	"zpcg/resources"
)

func TestApp(t *testing.T) {
	var _config config.Config
	_config.TimetableGobFileName = resources.TimetableGobFileName
	_app, err := app.NewApp(_config)
	assert.NoError(t, err)
	message, _ := _app.GenerateRoute(render.DefaultLanguageTag, "niksic, bar")
	t.Log("\n", message, "\n")
	assert.NotEmpty(t, message)
	numLines := strings.Count(message.Text, "\n")
	assert.Greater(t, numLines, 2) // there is at least header and at least one route
}
