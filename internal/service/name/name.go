package name

import (
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func NewStationNameResolver(
	unifiedStationNameToStationIdMap map[string]timetable.StationId,
	unifiedStationNameList [][]rune,
	stationIdToStationMap map[timetable.StationId]timetable.Station,
) *StationNameResolver {
	return &StationNameResolver{
		unifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		unifiedStationNameList:           unifiedStationNameList,
		stationIdToStationMap:            stationIdToStationMap,
	}
}

type StationNameResolver struct {
	unifiedStationNameToStationIdMap map[string]timetable.StationId
	unifiedStationNameList           [][]rune
	stationIdToStationMap            map[timetable.StationId]timetable.Station
}

func (s *StationNameResolver) FindStationIdByApproximateName(name string) (timetable.StationId, string, error) {
	unifiedName := Unify(name)
	match, err := findBestMatch([]rune(unifiedName), s.unifiedStationNameList)
	if err != nil {
		return 0, "", err
	}
	stationId := s.unifiedStationNameToStationIdMap[match]
	stationName := s.stationIdToStationMap[stationId].Name
	return stationId, stationName, nil
}
