package app

import (
	"fmt"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	timetable_gen "github.com/ivanov-gv/zpcg/gen/timetable"
	message_model "github.com/ivanov-gv/zpcg/internal/model/message"
	"github.com/ivanov-gv/zpcg/internal/model/message_render"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/date"
	"github.com/ivanov-gv/zpcg/internal/service/station_blacklist"
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
	message, _ := app.GenerateRoute(message_render.DefaultLanguageTag, fmt.Sprintf("%s, %s", NiksicWrongStationName, DanilovgradWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, DanilovgradStationName)
	assert.NotEmpty(t, message.InlineKeyboard)
	message, _ = app.GenerateRoute(message_render.DefaultLanguageTag, fmt.Sprintf("%s, %s", NiksicWrongStationName, BarWrongStationName))
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
	message, _ := app.GenerateRoute(message_render.DefaultLanguageTag, fmt.Sprintf("%s %s %s", NiksicWrongStationName, string(lo.SpecialCharset), DanilovgradWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, DanilovgradStationName)
	assert.NotEmpty(t, message.InlineKeyboard)
	message, _ = app.GenerateRoute(message_render.DefaultLanguageTag, fmt.Sprintf("%s %s %s", NiksicWrongStationName, string(lo.SpecialCharset), BarWrongStationName))
	t.Log("\n", message)
	assert.NotEmpty(t, message)
	assert.Contains(t, message.Text, NiksicStationName)
	assert.Contains(t, message.Text, BarStationName)
}

func TestBlackList(t *testing.T) {
	app, err := NewApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)

	for _, station := range station_blacklist.BlackListedStations {
		t.Run(station.Name, func(t *testing.T) {
			for _, language := range message_render.SupportedLanguages {
				t.Run(language.String(), func(t *testing.T) {
					message, err := app.GenerateRoute(language, fmt.Sprintf("%s, %s", BarStationName, station.Name))
					assert.NoError(t, err)
					assert.Contains(t, message.Text, station.LanguageTagToCustomErrorMessageMap[language])
				})
			}
		})
	}

}

func TestNameClashing(t *testing.T) {
	summerSeason, found := lo.Find(timetable_gen.Timetable.Seasons, func(item timetable.Season) bool {
		return strings.Contains(item.Name, "summer")
	})

	var option date.Option
	if found {
		option = date.FixedDate(summerSeason.Start)
	}

	app, err := NewApp(option)
	assert.NoError(t, err)
	assert.NotNil(t, app)

	mapExpectedStationToPossibleInput := map[string][]string{
		"Beograd": {
			"Белград", "Belgrade", "Beograde", "Belgrad", "Београд", "Beograd",
		},
		"Herceg Novi": {
			"Герцег нови", "Херцег Нови", "Херцег новый", "Герцег новый",
			"Herceg novi",
		},
		"Novi Beograd": {
			"Нови Белград", "Novi Belgrade", "Novi Beograde", "Novi Belgrad", "Нови Београд",
			"Новый Белград", "New Belgrade", "Novij Beograde", " Novii Belgrad", "Новый Београд",
		},
		"Stara Pazova": {
			"Стара пазова", "Старая пазова",
			"Stara Pazova",
		},
		"Novi Sad": {
			"Novi Sad", "New Sad", "Novij Sad",
			"Новый Сад", "Нови сад",
		},
	}

	for station, inputs := range mapExpectedStationToPossibleInput {
		t.Run(station, func(t *testing.T) {
			for _, originInput := range inputs {
				input := originInput + ", Podgorica"
				message, err := app.GenerateRoute(message_render.DefaultLanguageTag, input)
				assert.NoError(t, err)
				assert.Contains(t, message.Text, station,
					"'%s': '%s'", input, message.Text)
			}
		})
	}
}

func TestNoTrainsWarning(t *testing.T) {
	summerSeason, found := lo.Find(timetable_gen.Timetable.Seasons, func(item timetable.Season) bool {
		return strings.Contains(item.Name, "summer")
	})
	if !found {
		t.Skipf("summer season not found, can't test no trains warning properly")
	}

	app, err := NewApp(date.FixedDate(summerSeason.Start))
	assert.NoError(t, err)
	assert.NotNil(t, app)

	for _, languageTag := range message_render.SupportedLanguages {
		t.Run(languageTag.String(), func(t *testing.T) {
			message, _ := app.HandleMessage(message_model.Message{
				IsFilled: true,
				From: message_model.From{
					IsFilled:     true,
					LanguageCode: languageTag.String(),
				},
				Text:   "Podgorica, Novi Sad", // Novi sad is a summer train station
				ChatId: 0,
			})
			t.Log("\n", message)
			assert.NotEmpty(t, message)
		})
	}
}
