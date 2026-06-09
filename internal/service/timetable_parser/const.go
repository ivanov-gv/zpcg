package timetable_parser

import (
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/stations"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
)

var (
	// AliasesAsUnifiedStationNames is a list of unified stations aliases
	AliasesAsUnifiedStationNames = func() []string {
		var result []string
		for _, station := range stations.AliasesStationsList {
			unifiedAliases := lo.Map(station.Aliases, func(alias string, _ int) string {
				return station_name_resolver.Unify(alias)
			})
			result = append(result, unifiedAliases...)
		}
		return result
	}()

	// AliasesOriginalUnifiedStationNameToUnifiedAliasesMap is a map of unified station name -> slice of unified aliases
	AliasesOriginalUnifiedStationNameToUnifiedAliasesMap = lo.SliceToMap(stations.AliasesStationsList,
		func(stationAliases timetable.StationAliases) (string, []string) {
			return station_name_resolver.Unify(stationAliases.StationName), lo.Map(stationAliases.Aliases, station_name_resolver.UnifyIterative)
		})
)
