package timetable

const (
	TransferStationName = "Podgorica"
)

type ExportFormat struct {
	Seasons                          []Season
	UnifiedStationNameToStationIdMap map[string]StationId
	UnifiedStationNameList           [][]rune
	StationTypes                     map[StationTypeId]StationType
	TransferStationId                StationId
}
