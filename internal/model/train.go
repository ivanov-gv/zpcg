package model

type TrainId int

type TrainIdSet map[TrainId]struct{}

type TrainType int

const (
	LocalTrain TrainType = iota
	FastTrain
)

type TrainInfo struct {
	TrainId      TrainId
	TimetableUrl string
}
