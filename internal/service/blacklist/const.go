package blacklist

import (
	"github.com/samber/lo"
	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/stations"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/name"
)

var (
	UnifiedNames = lo.Map(BlackListedStations, func(s timetable.BlackListedStation, _ int) string {
		return name.Unify(s.Name)
	})

	UnifiedStationNameToStationIdMap = func() map[string]timetable.StationId {
		entries := lo.Map(BlackListedStations, func(item timetable.BlackListedStation, index int) lo.Entry[string, timetable.StationId] {
			stationId := timetable.StationId(-index) // negative, with zero
			unifiedName := name.Unify(item.Name)
			return lo.Entry[string, timetable.StationId]{
				Key:   unifiedName,
				Value: stationId,
			}
		})
		return lo.FromEntries(entries)
	}()

	BlackListedStations = lo.Flatten(
		lo.Map(stations.BlackListedStations, func(item struct {
			Names                              []string
			LanguageTagToCustomErrorMessageMap map[language.Tag]string
		}, _ int) []timetable.BlackListedStation {
			return lo.Map(item.Names, func(name string, _ int) timetable.BlackListedStation {
				return timetable.BlackListedStation{
					Name:                               name,
					LanguageTagToCustomErrorMessageMap: item.LanguageTagToCustomErrorMessageMap,
				}
			})
		}))
)
