package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/internal/app"
	"zpcg/internal/service/render"
)

func TestApp(t *testing.T) {
	_app, err := app.NewApp()
	assert.NoError(t, err)
	message, _ := _app.GenerateRoute(render.DefaultLanguageTag, "niksic, bar")
	t.Log("\n", message, "\n")
	assert.NotEmpty(t, message)
	numLines := strings.Count(message.Text, "\n")
	assert.Greater(t, numLines, 2) // there is at least header and at least one route
}
