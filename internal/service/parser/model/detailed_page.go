package model

import (
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

type DetailedTimetable struct {
	TrainId      timetable.TrainId
	TimetableUrl string
	Stops        []timetable.Stop
}
