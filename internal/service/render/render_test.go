package render

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"zpcg/internal/model"
)

func TestDirectRoutes(t *testing.T) {
	paths := []model.Path{
		{
			TrainId: 1111,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.UTC),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.UTC),
			},
		},
		{
			TrainId: 222,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.UTC),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.UTC),
			},
		},
	}
	stationsMap := map[model.StationId]model.Station{
		1: {
			Id:   1,
			Name: "Station1",
		},
		2: {
			Id:   2,
			Name: "Station2",
		},
	}
	trainMap := map[model.TrainId]model.TrainInfo{
		1111: {
			TrainId:      1111,
			TimetableUrl: "https:/somesite.com/timetable/1111",
		},
		222: {
			TrainId:      222,
			TimetableUrl: "https:/somesite.com/timetable/222",
		},
	}
	message, _ := NewRender(stationsMap, trainMap).DirectRoutes(paths)
	t.Logf("\n%s\n", message)
	assert.Contains(t, message, "1111](https:/somesite.com/timetable/1111")
	assert.Contains(t, message, "222](https:/somesite.com/timetable/222")
	assert.Contains(t, message, "12:00")
	assert.Contains(t, message, "12:40")
	assert.Contains(t, message, "08:00")
	assert.Contains(t, message, "08:40")
}

func TestTransferRoutes(t *testing.T) {
	paths := []model.Path{
		{
			TrainId: 1111,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.UTC),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.UTC),
			},
		},
		{
			TrainId: 222,
			Origin: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.UTC),
			},
			Destination: model.Stop{
				Id:        3,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.UTC),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.UTC),
			},
		},
	}
	stationsMap := map[model.StationId]model.Station{
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
	trainMap := map[model.TrainId]model.TrainInfo{
		1111: {
			TrainId:      1111,
			TimetableUrl: "https:/somesite.com/timetable/1111",
		},
		222: {
			TrainId:      222,
			TimetableUrl: "https:/somesite.com/timetable/222",
		},
	}
	message, _ := NewRender(stationsMap, trainMap).TransferRoutes(paths, 1, 2, 3)
	t.Logf("\n%s\n", message)
	assert.Contains(t, message, "1111](https:/somesite.com/timetable/1111")
	assert.Contains(t, message, "222](https:/somesite.com/timetable/222")
	assert.Contains(t, message, "12:00")
	assert.Contains(t, message, "12:40")
	assert.Contains(t, message, "08:00")
	assert.Contains(t, message, "08:40")
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
