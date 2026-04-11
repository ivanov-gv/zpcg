package render

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/zelenin/go-tdlib/client"
	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/message"
	model_render "github.com/ivanov-gv/zpcg/internal/model/render"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func renderTestDirectRoutes(tag language.Tag, _time time.Time, updateCallback, reverseCallback string) message.Response {
	paths := []timetable.Path{
		{
			TrainId: 1111,
			Origin: timetable.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.UTC),
			},
			Destination: timetable.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.UTC),
			},
		},
		{
			TrainId: 222,
			Origin: timetable.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.UTC),
			},
			Destination: timetable.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.UTC),
			},
		},
	}
	stationsMap := map[timetable.StationId]timetable.Station{
		1: {
			Id:   1,
			Name: "Station1",
		},
		2: {
			Id:   2,
			Name: "Station2",
		},
	}
	trainMap := map[timetable.TrainId]timetable.TrainInfo{
		1111: {
			TrainId:      1111,
			TimetableUrl: "https:/somesite.com/timetable/1111",
		},
		222: {
			TrainId:      222,
			TimetableUrl: "https:/somesite.com/timetable/222",
		},
	}
	return NewRender(stationsMap, trainMap).DirectRoutes(tag, paths, _time, updateCallback, reverseCallback)
}

func TestDirectRoutes(t *testing.T) {
	_message := renderTestDirectRoutes(model_render.DefaultLanguageTag, time.Time{}, "updateCallback", "reverseCallback")
	t.Logf("\n%v\n", _message)
	assert.Contains(t, _message.Text, "1111")
	assert.Contains(t, _message.Text, "222")
	assert.Contains(t, _message.Text, "12:10") // origin departure
	assert.Contains(t, _message.Text, "12:30") // destination arrival
	assert.Contains(t, _message.Text, "08:10") // origin departure
	assert.Contains(t, _message.Text, "08:30") // destination arrival
}

func renderTestTransferRoutes(tag language.Tag, _time time.Time, updateCallback, reverseCallback string) message.Response {
	paths := []timetable.Path{
		{
			TrainId: 1111,
			Origin: timetable.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.UTC),
			},
			Destination: timetable.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.UTC),
			},
		},
		{
			TrainId: 222,
			Origin: timetable.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.UTC),
			},
			Destination: timetable.Stop{
				Id:        3,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.UTC),
			},
		},
	}
	stationsMap := map[timetable.StationId]timetable.Station{
		1: {
			Id:   1,
			Name: "Station1",
		},
		2: {
			Id:   2,
			Name: "Station2",
		},
		3: {
			Id:   3,
			Name: "Station3",
		},
	}
	trainMap := map[timetable.TrainId]timetable.TrainInfo{
		1111: {
			TrainId:      1111,
			TimetableUrl: "https:/somesite.com/timetable/1111",
		},
		222: {
			TrainId:      222,
			TimetableUrl: "https:/somesite.com/timetable/222",
		},
	}
	return NewRender(stationsMap, trainMap).TransferRoutes(tag, paths, _time, 1, 2, 3,
		updateCallback, reverseCallback)
}

func TestTransferRoutes(t *testing.T) {
	_message := renderTestTransferRoutes(model_render.DefaultLanguageTag, time.Time{}, "updateCallback", "reverseCallback")
	t.Logf("\n%v\n", _message)
	assert.Contains(t, _message.Text, "1111")
	assert.Contains(t, _message.Text, "222")
	assert.Contains(t, _message.Text, "12:10")
	assert.Contains(t, _message.Text, "12:30")
	assert.Contains(t, _message.Text, "08:10")
	assert.Contains(t, _message.Text, "08:30")
}

func TestConstants(t *testing.T) {
	var constantsToTest = map[string]map[language.Tag]string{
		"ErrorMessageMap":                 model_render.ErrorMessageMap,
		"StartMessageMap":                 model_render.StartMessageMap,
		"StationDoesNotExistMessageMap":   model_render.StationDoesNotExistMessageMap,
		"RailwayMapButtonTextMap":         model_render.RailwayMapButtonTextMap,
		"ReverseRouteInlineButtonTextMap": model_render.ReverseRouteInlineButtonTextMap,
		"AlertUpdateNotificationTextMap":  model_render.AlertUpdateNotificationTextMap,
		"SimpleUpdateNotificationTextMap": model_render.SimpleUpdateNotificationTextMap,
		"OfficialTimetableUrlTextMap":     model_render.OfficialTimetableUrlTextMap,
		"MapMessageMap":                   model_render.MapMessageMap,
		"HelpMessageMap":                  model_render.HelpMessageMap,
		"AboutMessageMap":                 model_render.AboutMessageMap,
	}

	for name, _map := range constantsToTest {
		t.Run(name, func(t *testing.T) {
			// all the supported languages are present
			languagesSortFunction := func(a, b language.Tag) int { return strings.Compare(a.String(), b.String()) }
			actualLanguages := lo.Keys(_map)
			slices.SortFunc(actualLanguages, languagesSortFunction)
			expectedLanguages := model_render.SupportedLanguages
			slices.SortFunc(model_render.SupportedLanguages, languagesSortFunction)
			assert.EqualValuesf(t, expectedLanguages, actualLanguages, "all the supported languages are present")
			// there is no repeating values (i.e. set of keys ~ set of values)
			valuesSet := lo.SliceToMap(lo.Values(_map), func(item string) (string, struct{}) { return item, struct{}{} })
			assert.NotContainsf(t, lo.Keys(valuesSet), "", "no empty values")
			assert.Equal(t, len(lo.Keys(_map)), len(valuesSet), "there are no equal messages accidentally mapped for different languages")
		})
	}
}

func TestAlertMessage(t *testing.T) {
	const MaxTextLen = 200
	for lang, text := range model_render.AlertUpdateNotificationTextMap {
		assert.Less(t, len(text), MaxTextLen, "lang tag: %s", lang.String())
	}
	for lang, text := range model_render.SimpleUpdateNotificationTextMap {
		assert.Less(t, len(text), MaxTextLen, "lang tag: %s", lang.String())
	}
}

func TestBelarusianLanguage(t *testing.T) {
	parsed := ParseLanguageTag("be")
	assert.Equal(t, parsed, model_render.Belarusian)
	assert.Equal(t, parsed.String(), "be")
}

func TestMarkdownMessages(t *testing.T) {
	t.Run("Markdown constants", func(t *testing.T) {
		var constantsToTest = map[string]map[language.Tag]string{
			"StartMessageMap": model_render.StartMessageMap,
		}
		for name, _map := range constantsToTest {
			for languageTag, _message := range _map {
				t.Run(name+" "+languageTag.String(), func(t *testing.T) {
					_, err := client.ParseTextEntities(&client.ParseTextEntitiesRequest{
						Text: _message,
						ParseMode: &client.TextParseModeMarkdown{
							Version: 2,
						},
					})
					assert.NoError(t, err, "%+v", err)
				})
			}
		}
	})

	t.Run("Markdown direct route", func(t *testing.T) {
		for _, languageTag := range model_render.SupportedLanguages {
			t.Run(languageTag.String(), func(t *testing.T) {
				_message := renderTestDirectRoutes(languageTag, time.Now(), "", "")
				_, err := client.ParseTextEntities(&client.ParseTextEntitiesRequest{
					Text: _message.Text,
					ParseMode: &client.TextParseModeMarkdown{
						Version: 2,
					},
				})
				assert.NoError(t, err, "%+v", err)
			})
		}
	})
	t.Run("Markdown transfer route", func(t *testing.T) {
		for _, languageTag := range model_render.SupportedLanguages {
			t.Run(languageTag.String(), func(t *testing.T) {
				_message := renderTestTransferRoutes(languageTag, time.Now(), "", "")
				_, err := client.ParseTextEntities(&client.ParseTextEntitiesRequest{
					Text: _message.Text,
					ParseMode: &client.TextParseModeMarkdown{
						Version: 2,
					},
				})
				assert.NoError(t, err, "%+v", err)
			})
		}
	})
}

func TestBotInfo(t *testing.T) {
	for languageTag, bot := range model_render.BotInfoMap {
		t.Run(languageTag.String(), func(t *testing.T) {
			assert.Less(t, len([]rune(bot.Name)), 64)
			assert.Less(t, len([]rune(bot.Description)), 512)
			assert.Less(t, len([]rune(bot.ShortDescription)), 120)
			for _, command := range model_render.AllCommands {
				assert.Contains(t, bot.CommandNames, command)
			}
		})
	}
}
