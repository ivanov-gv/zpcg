package detailed_page

import (
	"github.com/samber/lo"

	"zpcg/internal/model"
)

var (
	lastId      model.StationId = 0
	nameToIdMap                 = map[string]model.StationId{}
)

func generateStationId(stationName string) model.StationId {
	if id, found := nameToIdMap[stationName]; found {
		return id
	}
	lastId++
	nameToIdMap[stationName] = lastId
	return lastId
}

func GetStationIdToNameMap() map[model.StationId]string {
	return lo.Invert(nameToIdMap)
}
