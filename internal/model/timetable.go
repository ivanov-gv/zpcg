package model

const (
	TransferStationName = "Podgorica"
)

type TimetableTransferFormat struct {
	StationIdToTrainIdSet            map[StationId]TrainIdSet
	TrainIdToStationMap              map[TrainId]StationIdToStationMap
	StationIdToStationMap            map[StationId]Station
	TrainIdToTrainInfoMap            map[TrainId]TrainInfo
	UnifiedStationNameToStationIdMap map[string]StationId
	UnifiedStationNameList           []string
	TransferStationId                StationId
	BlacklistedStations              []BlackListedStation
}
