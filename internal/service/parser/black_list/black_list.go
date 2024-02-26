package black_list

import (
	"github.com/samber/lo"

	"zpcg/internal/model"
	"zpcg/internal/service/name"
)

var (
	BlackListUnifiedNames = lo.Map(BlackListedStations, func(s model.BlackListedStation, _ int) string {
		return name.Unify(s.Name)
	})

	BlackListUnifiedStationNameToStationIdMap = func() map[string]model.StationId {
		entries := lo.Map(BlackListedStations, func(item model.BlackListedStation, index int) lo.Entry[string, model.StationId] {
			stationId := model.StationId(-index) // negative, with zero
			unifiedName := name.Unify(item.Name)
			return lo.Entry[string, model.StationId]{
				Key:   unifiedName,
				Value: stationId,
			}
		})
		return lo.FromEntries(entries)
	}()

	BlackListedStations = []model.BlackListedStation{
		//{ TODO: uncomment when ready
		//	Name:                               "Budva",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Tivat",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Kotor",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Cetinje",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Perast",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Durmitor",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Petrovac",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Ulcinj",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Sveti Stefan",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Becici",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Herceg Novi",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Savnik",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Zabljak",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name:                               "Durmitor",
		//	LanguageTagToCustomErrorMessageMap: nil,
		//},
		//{
		//	Name: "Tirana",
		//	LanguageTagToCustomErrorMessageMap: map[string]string{
		//		language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории.",
		//		language.English.String(): "There are no routes between Montenegro and Albania.",
		//	},
		//},
		//{
		//	Name: "Sarajevo",
		//	LanguageTagToCustomErrorMessageMap: map[string]string{
		//		language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории.",
		//		language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina.",
		//	},
		//},
	}
)
