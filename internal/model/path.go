package model

type Stop struct {
	Station
	TrainId TrainId
}

type Path struct {
	Origin, Destination Stop
}
