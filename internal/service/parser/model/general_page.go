package model

import (
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

type GeneralTimetableRow struct {
	TrainId               timetable.TrainId
	TrainType             timetable.TrainType
	DetailedTimetableLink string
}
