package model

import "zpcg/internal/model"

type GeneralTimetableRow struct {
	TrainId               model.TrainId
	TrainType             model.TrainType
	DetailedTimetableLink string
}
