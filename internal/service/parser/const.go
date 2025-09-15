package parser

import (
	"github.com/samber/lo"

	"github.com/ivanov-gv/zpcg/internal/model/stations"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/name"
)

var (
	// AliasesAsUnifiedStationNames is a list of unified stations aliases
	AliasesAsUnifiedStationNames = func() []string {
		var result []string
		for _, station := range stations.AliasesStationsList {
			unifiedAliases := lo.Map(station.Aliases, func(alias string, _ int) string {
				return name.Unify(alias)
			})
			result = append(result, unifiedAliases...)
		}
		return result
	}()

	// AliasesWithUnifiedNames is a list of unified stations aliases
	AliasesWithUnifiedNames = lo.Map(stations.AliasesStationsList,
		func(item timetable.StationAliases, _ int) timetable.StationAliases {
			return timetable.StationAliases{
				StationName: name.Unify(item.StationName),
				Aliases:     lo.Map(item.Aliases, name.UnifyIterative),
			}
		})
)
