package timetable

const (
	TransferStationName = "Podgorica"
)

type ExportFormat struct {
	Seasons                          []Season
	UnifiedStationNameList           [][]rune
	UnifiedStationNameToStationIdMap map[string]StationId
	StationIdToStationMap            map[StationId]Station
	StationTypes                     map[StationTypeId]StationType
	TransferStationId                StationId
}
