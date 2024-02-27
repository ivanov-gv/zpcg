package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/internal/config"
	"zpcg/internal/service/blacklist"
	"zpcg/internal/service/render"
	"zpcg/resources"
)

func TestNewApp(t *testing.T) {
	var _config config.Config
	_config.TimetableGobFileName = resources.TimetableGobFileName
	app, err := NewApp(_config)
	assert.NoError(t, err)
	assert.NotNil(t, app)
}

const (
	NiksicStationName           = "Nikšić"
	NiksicWrongStationName      = "niksic"
	DanilovgradStationName      = "Danilovgrad"
	DanilovgradWrongStationName = "DaNil ovgrad"
	BarStationName              = "Bar"
	BarWrongStationName         = "Barrrrrr"
)

func TestGenerateRoute(t *testing.T) {
	var _config config.Config
	_config.TimetableGobFileName = resources.TimetableGobFileName
	app, err := NewApp(_config)
	assert.NoError(t, err)
	assert.NotNil(t, app)
	message, _, _ := app.GenerateRoute(render.DefaultLanguageTag, NiksicWrongStationName, DanilovgradWrongStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, DanilovgradStationName)
	message, _, _ = app.GenerateRoute(render.DefaultLanguageTag, NiksicWrongStationName, BarWrongStationName)
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message, NiksicStationName)
	assert.Contains(t, message, BarStationName)
}

func TestBlackList(t *testing.T) {
	var _config config.Config
	_config.TimetableGobFileName = resources.TimetableGobFileName
	app, err := NewApp(_config)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	for _, station := range blacklist.BlackListedStations {
		for _, language := range render.SupportedLanguages {
			t.Run(station.Name+"/"+language.String(), func(t *testing.T) {
				message, _, err := app.GenerateRoute(language, BarStationName, station.Name)
				assert.NoError(t, err)
				assert.Contains(t, message, station.LanguageTagToCustomErrorMessageMap[language.String()])
			})
		}
	}

}
