package render

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func TestDirectRoutes(t *testing.T) {
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
	message := NewRender(stationsMap, trainMap).DirectRoutes(DefaultLanguageTag, paths, time.Time{},
		"updateCallback", "reverseCallback")
	t.Logf("\n%v\n", message)
	assert.Contains(t, message.Text, "1111](https:/somesite.com/timetable/1111")
	assert.Contains(t, message.Text, "222](https:/somesite.com/timetable/222")
	assert.Contains(t, message.Text, "12:10") // origin departure
	assert.Contains(t, message.Text, "12:30") // destination arrival
	assert.Contains(t, message.Text, "08:10") // origin departure
	assert.Contains(t, message.Text, "08:30") // destination arrival
}

func TestTransferRoutes(t *testing.T) {
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
	message := NewRender(stationsMap, trainMap).TransferRoutes(DefaultLanguageTag, paths, time.Time{}, 1, 2, 3,
		"updateCallback", "reverseCallback")
	t.Logf("\n%v\n", message)
	assert.Contains(t, message.Text, "1111](https:/somesite.com/timetable/1111")
	assert.Contains(t, message.Text, "222](https:/somesite.com/timetable/222")
	assert.Contains(t, message.Text, "12:10")
	assert.Contains(t, message.Text, "12:30")
	assert.Contains(t, message.Text, "08:10")
	assert.Contains(t, message.Text, "08:30")
}

func TestConstants(t *testing.T) {
	var constantsToTest = map[string]map[language.Tag]string{
		"ErrorMessageMap":                     ErrorMessageMap,
		"StartMessageMap":                     StartMessageMap,
		"StationDoesNotExistMessageMap":       StationDoesNotExistMessageMap,
		"StationDoesNotExistMessageSuffixMap": StationDoesNotExistMessageSuffixMap,
	}

	for name, _map := range constantsToTest {
		t.Run(name, func(t *testing.T) {
			// all the supported languages are present
			languagesSortFunction := func(a, b language.Tag) int { return strings.Compare(a.String(), b.String()) }
			actualLanguages := lo.Keys(_map)
			slices.SortFunc(actualLanguages, languagesSortFunction)
			expectedLanguages := SupportedLanguages
			slices.SortFunc(SupportedLanguages, languagesSortFunction)
			assert.EqualValuesf(t, expectedLanguages, actualLanguages, "all the supported languages are present")
			// there is no repeating values (i.e. set of keys ~ set of values)
			valuesSet := lo.SliceToMap(lo.Values(_map), func(item string) (string, struct{}) { return item, struct{}{} })
			assert.Equal(t, len(lo.Keys(_map)), len(valuesSet), "there are no equal messages accidentally mapped for different languages")
		})
	}
}

func TestAlertMessage(t *testing.T) {
	const MaxTextLen = 200
	for _, text := range lo.Values(AlertUpdateNotificationTextMap) {
		assert.Less(t, len(text), MaxTextLen)
	}
	for _, text := range lo.Values(SimpleUpdateNotificationTextMap) {
		assert.Less(t, len(text), MaxTextLen)
	}
}
