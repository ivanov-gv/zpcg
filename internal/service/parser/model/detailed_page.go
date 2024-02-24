package model

import "zpcg/internal/model"

type DetailedTimetable struct {
	TrainId      model.TrainId
	TimetableUrl string
	Stops        []model.Stop
}
