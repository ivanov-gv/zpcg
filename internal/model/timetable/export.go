package timetable

const (
	TransferStationName = "Podgorica"
)

type ExportFormat struct {
	StationIdToTrainIdSet            map[StationId]TrainIdSet
	TrainIdToStationMap              map[TrainId]StationIdToStationMap
	StationIdToStationMap            map[StationId]Station
	TrainIdToTrainInfoMap            map[TrainId]TrainInfo
	UnifiedStationNameToStationIdMap map[string]StationId
	UnifiedStationNameList           [][]rune
	StationTypes                     map[StationTypeId]StationType
	TransferStationId                StationId
}
