package app

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/ivanov-gv/zpcg/internal/service/blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/render"
)

func TestNewApp(t *testing.T) {
	app, err := NewApp()
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
	app, err := NewApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)
	message, _ := app.GenerateRoute(render.DefaultLanguageTag, fmt.Sprintf("%s, %s", NiksicWrongStationName, DanilovgradWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, DanilovgradStationName)
	assert.NotEmpty(t, message.InlineKeyboard)
	message, _ = app.GenerateRoute(render.DefaultLanguageTag, fmt.Sprintf("%s, %s", NiksicWrongStationName, BarWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, BarStationName)
	assert.NotEmpty(t, message.InlineKeyboard)
}

func TestGenerateRouteWithCustomDelimiter(t *testing.T) {
	app, err := NewApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)
	message, _ := app.GenerateRoute(render.DefaultLanguageTag, fmt.Sprintf("%s %s %s", NiksicWrongStationName, string(lo.SpecialCharset), DanilovgradWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, DanilovgradStationName)
	assert.NotEmpty(t, message.InlineKeyboard)
	message, _ = app.GenerateRoute(render.DefaultLanguageTag, fmt.Sprintf("%s %s %s", NiksicWrongStationName, string(lo.SpecialCharset), BarWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, BarStationName)
}

func TestBlackList(t *testing.T) {
	app, err := NewApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)

	for _, station := range blacklist.BlackListedStations {
		t.Run(station.Name, func(t *testing.T) {
			for _, language := range render.SupportedLanguages {
				t.Run(language.String(), func(t *testing.T) {
					message, err := app.GenerateRoute(language, fmt.Sprintf("%s, %s", BarStationName, station.Name))
					assert.NoError(t, err)
					assert.Contains(t, message.Text, station.LanguageTagToCustomErrorMessageMap[language.String()])
				})
			}
		})
	}

}
