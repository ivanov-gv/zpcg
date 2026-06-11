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
	NiksicStationName              = "Nikšić"
	NiksicWrongStationName         = "niksic"
	NiksicCyrillicStationName      = "никшич"
	DanilovgradStationName         = "Danilovgrad"
	DanilovgradWrongStationName    = "DaNil ovgrad"
	DanilovgradCyrillicStationName = "даниловград"
	BarStationName                 = "Bar"
	BarWrongStationName            = "Barrrrrr"
	BarCyrillicStationName         = "Бар"
)

func TestGenerateRoute(t *testing.T) {
	app, err := NewApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)

	testCases := []struct {
		station1, station2   string
		expected1, expected2 string
	}{
		{
			station1: NiksicWrongStationName, station2: DanilovgradWrongStationName,
			expected1: NiksicStationName, expected2: DanilovgradStationName,
		},
		{
			station1: NiksicWrongStationName, station2: BarWrongStationName,
			expected1: NiksicStationName, expected2: BarStationName,
		},
		{
			station1: DanilovgradCyrillicStationName, station2: BarCyrillicStationName,
			expected1: DanilovgradStationName, expected2: BarStationName,
		},
		{
			station1: NiksicCyrillicStationName, station2: DanilovgradCyrillicStationName,
			expected1: NiksicStationName, expected2: DanilovgradStationName,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.station1+" -> "+testCase.station2, func(t *testing.T) {
			message, _ := app.GenerateRoute(message_render.DefaultLanguageTag, fmt.Sprintf("%s, %s", testCase.station1, testCase.station2))
			t.Log("\n", message)
			assert.NotEmpty(t, message)
			assert.Contains(t, message.Text, testCase.expected1)
			assert.Contains(t, message.Text, testCase.expected2)
			assert.NotEmpty(t, message.InlineKeyboard)
		})
	}
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
	mapExpectedStationToPossibleInput_WinterSeason := map[string][]string{
		"Beograd": {
			"Белград", "Belgrade", "Beograde", "Belgrad", "Београд", "Beograd",
		},
		"Herceg Novi": {
			"Герцег нови", "Херцег Нови", "Херцег новый", "Герцег новый",
			"Herceg novi",
		},
		"Aerodrom": {
			"airport", "аэропорт",
		},
	}

	mapExpectedStationToPossibleInput_SummerSeason := lo.Assign(mapExpectedStationToPossibleInput_WinterSeason,
		map[string][]string{
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
			"Bačka Topola": {
				"бачка топола",
			},
		},
	)

	winterSeason, found := lo.Find(timetable_gen.Timetable.Seasons, func(item timetable.Season) bool {
		return !strings.Contains(item.Name, "summer")
	})

	t.Run(winterSeason.Name, func(t *testing.T) {
		if !found {
			t.Skipf("winter season not found, can't test name clashing properly")
		}

		var option date.Option
		if found {
			option = date.FixedDate(winterSeason.Start)
		}

		app, err := NewApp(option)
		assert.NoError(t, err)
		assert.NotNil(t, app)

		for station, inputs := range mapExpectedStationToPossibleInput_WinterSeason {
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
	})

	summerSeason, found := lo.Find(timetable_gen.Timetable.Seasons, func(item timetable.Season) bool {
		return strings.Contains(item.Name, "summer")
	})

	t.Run(summerSeason.Name, func(t *testing.T) {
		var option date.Option
		if found {
			option = date.FixedDate(summerSeason.Start)
		}

		app, err := NewApp(option)
		assert.NoError(t, err)
		assert.NotNil(t, app)

		for station, inputs := range mapExpectedStationToPossibleInput_SummerSeason {
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
	})
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

func TestNonexistentStationInputWarning(t *testing.T) {
	app, err := NewApp()
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
				Text:   "Podgorica, Istanbul",
				ChatId: 0,
			})
			t.Log("\n", message)
			assert.NotEmpty(t, message)
			assert.Contains(t, lo.FirstOrEmpty(message.Send).Text, "Istanbul: ")
		})
	}
}
