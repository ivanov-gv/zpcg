package model

import "zpcg/internal/model"

type DetailedTimetable struct {
	TrainId  model.TrainId
	Stations []model.Station
}
