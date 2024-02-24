package name

import "zpcg/internal/model"

func NewStationNameResolver(
	unifiedStationNameToStationIdMap map[string]model.StationId,
	unifiedStationNameList []string,
) *StationNameResolver {
	return &StationNameResolver{
		unifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		unifiedStationNameList:           unifiedStationNameList,
	}
}

type StationNameResolver struct {
	unifiedStationNameToStationIdMap map[string]model.StationId
	unifiedStationNameList           []string
}

func (s *StationNameResolver) FindStationIdByApproximateName(name string) (model.StationId, error) {
	unifiedName := Unify(name)
	match, err := findBestMatch(unifiedName, s.unifiedStationNameList)
	if err != nil {
		return 0, err
	}
	return s.unifiedStationNameToStationIdMap[match], nil
}
