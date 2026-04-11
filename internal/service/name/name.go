package name

import (
	"github.com/samber/lo"

	approxmatch "github.com/ivanov-gv/approximate-match"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

// scoreThreshold is the minimum match score for a station name to be accepted.
// Set slightly below the library default (0.45) to account for floating-point
// rounding in borderline cases (e.g. "Barrrrrr" → "bar" computes to 0.4499...).
// Still well above the noise floor: unrelated European capitals score at most 0.30.
const scoreThreshold = 0.44

func NewStationNameResolver(
	unifiedStationNameToStationIdMap map[string]timetable.StationId,
	unifiedStationNameList [][]rune,
	stationIdToStationMap map[timetable.StationId]timetable.Station,
) *StationNameResolver {
	stationNames := lo.Map(unifiedStationNameList, func(r []rune, _ int) string { return string(r) })
	threshold := scoreThreshold
	matcher := approxmatch.NewMatcher(stationNames, &threshold)
	return &StationNameResolver{
		unifiedStationNameToStationIdMap: unifiedStationNameToStationIdMap,
		matcher:                          matcher,
		stationIdToStationMap:            stationIdToStationMap,
	}
}

type StationNameResolver struct {
	unifiedStationNameToStationIdMap map[string]timetable.StationId
	matcher                          *approxmatch.Matcher
	stationIdToStationMap            map[timetable.StationId]timetable.Station
}

func (s *StationNameResolver) FindStationIdByApproximateName(name string) (timetable.StationId, string, error) {
	matches := s.matcher.Find(name)
	if len(matches) == 0 {
		return 0, "", ErrNoMatchesFound
	}
	bestMatch := matches[0].Word
	stationId := s.unifiedStationNameToStationIdMap[bestMatch]
	stationName := s.stationIdToStationMap[stationId].Name
	return stationId, stationName, nil
}
