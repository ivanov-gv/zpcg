package name

import (
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func NewStationNameResolver(
	unifiedStationNameToStationIdMap map[string]timetable.StationId,
	unifiedStationNameList [][]rune,
) *StationNameResolver {
	return &StationNameResolver{
		unifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		unifiedStationNameList:           unifiedStationNameList,
	}
}

type StationNameResolver struct {
	unifiedStationNameToStationIdMap map[string]timetable.StationId
	unifiedStationNameList           [][]rune
}

func (s *StationNameResolver) FindStationIdByApproximateName(name string) (timetable.StationId, string, error) {
	unifiedName := Unify(name)
	match, err := findBestMatch([]rune(unifiedName), s.unifiedStationNameList)
	if err != nil {
		return 0, "", err
	}
	return s.unifiedStationNameToStationIdMap[match], match, nil
}
