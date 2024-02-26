package alias

import (
	"github.com/samber/lo"

	"zpcg/internal/service/name"
)

type StationAliases struct {
	StationName string
	Aliases     []string
}

var (
	// AliasesAsUnifiedStationNames is a list of unified stations aliases
	AliasesAsUnifiedStationNames = func() []string {
		var result []string
		for _, station := range AliasesStationsList {
			unifiedAliases := lo.Map(station.Aliases, func(alias string, _ int) string {
				return name.Unify(alias)
			})
			result = append(result, unifiedAliases...)
		}
		return result
	}()

	// AliasOriginalUnifiedStationNameToUnifiedAliasesMap is a map of unified station name -> slice of unified aliases
	AliasOriginalUnifiedStationNameToUnifiedAliasesMap = lo.SliceToMap(AliasesStationsList,
		func(stationAliases StationAliases) (string, []string) {
			return name.Unify(stationAliases.StationName), lo.Map(stationAliases.Aliases, name.UnifyIterative)
		})

	AliasesStationsList = []StationAliases{
		{
			StationName: "Beograd Centar",
			Aliases:     []string{"belgrad"},
		},
	}
)
