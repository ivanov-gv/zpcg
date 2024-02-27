package blacklist

import (
	"github.com/samber/lo"
	"golang.org/x/text/language"

	"zpcg/internal/model"
	"zpcg/internal/service/name"
)

var (
	UnifiedNames = lo.Map(BlackListedStations, func(s model.BlackListedStation, _ int) string {
		return name.Unify(s.Name)
	})

	UnifiedStationNameToStationIdMap = func() map[string]model.StationId {
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
		{
			Name: "Budva",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Будве нет жд станций",
				language.English.String(): "There is no train station in Budva",
			},
		},
		{
			Name: "Tivat",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Тивате нет жд станций",
				language.English.String(): "There is no train station in Tivat",
			},
		},
		{
			Name: "Kotor",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Которе нет жд станций",
				language.English.String(): "There is no train station in Kotor",
			},
		},
		{
			Name: "Cetinje",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Цетине нет жд станций",
				language.English.String(): "There is no train station in Cetinje",
			},
		},
		{
			Name: "Perast",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Перасте нет жд станций",
				language.English.String(): "There is no train station in Perast",
			},
		},
		{
			Name: "Durmitor",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Дурмиторе нет жд станций",
				language.English.String(): "There is no train station in Durmitor",
			},
		},
		{
			Name: "Petrovac",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Петроваце нет жд станций",
				language.English.String(): "There is no train station in Petrovac",
			},
		},
		{
			Name: "Ulcinj",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Ульцине нет жд станций",
				language.English.String(): "There is no train station in Ulcinj",
			},
		},
		{
			Name: "Sveti Stefan",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Свети Стефан нет жд станций",
				language.English.String(): "There is no train station in Sveti Stefan",
			},
		},
		{
			Name: "Becici",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Бечичи нет жд станций",
				language.English.String(): "There is no train station in Becici",
			},
		},
		{
			Name: "Herceg Novi",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Херцег-Нови нет жд станций",
				language.English.String(): "There is no train station in Herceg Novi",
			},
		},
		{
			Name: "Savnik",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Шавнике нет жд станций",
				language.English.String(): "There is no train station in Savnik",
			},
		},
		{
			Name: "Zabljak",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Жабляке нет жд станций",
				language.English.String(): "There is no train station in Zabljak",
			},
		},
		{
			Name: "Durmitor",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Дурмиторе нет жд станций",
				language.English.String(): "There is no train station in Durmitor",
			},
		},
		{
			Name: "Tirana",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории.",
				language.English.String(): "There are no routes between Montenegro and Albania.",
			},
		},
		{
			Name: "Sarajevo",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории.",
				language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina.",
			},
		},
	}
)
