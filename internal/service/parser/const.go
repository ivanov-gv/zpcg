package parser

import (
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/render"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/name"
)

var (
	// AliasesAsUnifiedStationNames is a list of unified stations aliases
	AliasesAsUnifiedStationNames = func() []string {
		var result []string
		for _, station := range render.AliasesStationsList {
			unifiedAliases := lo.Map(station.Aliases, func(alias string, _ int) string {
				return name.Unify(alias)
			})
			result = append(result, unifiedAliases...)
		}
		return result
	}()

	// AliasesOriginalUnifiedStationNameToUnifiedAliasesMap is a map of unified station name -> slice of unified aliases
	AliasesOriginalUnifiedStationNameToUnifiedAliasesMap = lo.SliceToMap(render.AliasesStationsList,
		func(stationAliases timetable.StationAliases) (string, []string) {
			return name.Unify(stationAliases.StationName), lo.Map(stationAliases.Aliases, name.UnifyIterative)
		})
)
