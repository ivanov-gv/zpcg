package model

type TrainId int

type TrainIdSet map[TrainId]struct{}

type TrainInfo struct {
	TrainId      TrainId
	TimetableUrl string
}
