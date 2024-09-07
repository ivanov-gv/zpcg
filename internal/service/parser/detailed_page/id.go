package detailed_page

import (
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

var (
	lastId      timetable.StationId = 0
	nameToIdMap                     = map[string]timetable.StationId{}
)

// generateStationId generates a unique StationId for the given stationName.
// If the station is already present in the nameToIdMap,
// it returns the existing id. Otherwise, it generates a new id, which is more or equal to 1.
func generateStationId(stationName string) timetable.StationId {
	if id, found := nameToIdMap[stationName]; found {
		return id
	}
	lastId++
	nameToIdMap[stationName] = lastId
	return lastId
}

func getStationIdToNameMap() map[timetable.StationId]string {
	return lo.Invert(nameToIdMap)
}
