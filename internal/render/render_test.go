package render

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"zpcg/internal/model"
)

func TestDirectRoutes(t *testing.T) {
	paths := []model.Path{
		{
			TrainId: 1111,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.Local),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.Local),
			},
		},
		{
			TrainId: 222,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.Local),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.Local),
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
	assert.Contains(t, message, "1111](https:/somesite.com/timetable/1111)")
	assert.Contains(t, message, "222](https:/somesite.com/timetable/222)")
	assert.Contains(t, message, "12:10")
	assert.Contains(t, message, "12:30")
	assert.Contains(t, message, "08:10")
	assert.Contains(t, message, "08:30")
}

func TestTransferRoutes(t *testing.T) {
	paths := []model.Path{
		{
			TrainId: 1111,
			Origin: model.Stop{
				Id:        1,
				Arrival:   time.Date(0, 0, 0, 12, 00, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 12, 10, 0, 0, time.Local),
			},
			Destination: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 12, 30, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 12, 40, 0, 0, time.Local),
			},
		},
		{
			TrainId: 222,
			Origin: model.Stop{
				Id:        2,
				Arrival:   time.Date(0, 0, 0, 8, 00, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 8, 10, 0, 0, time.Local),
			},
			Destination: model.Stop{
				Id:        3,
				Arrival:   time.Date(0, 0, 0, 8, 30, 0, 0, time.Local),
				Departure: time.Date(0, 0, 0, 8, 40, 0, 0, time.Local),
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
	assert.Contains(t, message, "1111](https:/somesite.com/timetable/1111)")
	assert.Contains(t, message, "222](https:/somesite.com/timetable/222)")
	assert.Contains(t, message, "12:10")
	assert.Contains(t, message, "12:30")
	assert.Contains(t, message, "08:10")
	assert.Contains(t, message, "08:30")
}
