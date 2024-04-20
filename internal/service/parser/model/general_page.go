package model

import (
	"zpcg/internal/model/timetable"
)

type GeneralTimetableRow struct {
	TrainId               timetable.TrainId
	TrainType             timetable.TrainType
	DetailedTimetableLink string
}
