package blacklist

import (
	"github.com/samber/lo"
	"golang.org/x/text/language"

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

	BlackListedStations = []timetable.BlackListedStation{
		{
			Name: "Budva",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Будве нет жд станций",
				language.English.String(): "There is no train station in Budva",
			},
		},
		{
			Name: "Будва",
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
			Name: "Тиват",
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
			Name: "Котор",
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
			Name: "Цетине",
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
			Name: "Пераст",
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
			Name: "Дурмитор",
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
			Name: "Петровац",
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
			Name: "Ульцин",
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
			Name: "Свети Стефан",
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
			Name: "Бечичи",
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
			Name: "Херцег Нови",
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
			Name: "Шавник",
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
			Name: "Жабляк",
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
			Name: "Дурмитор",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "В Дурмиторе нет жд станций",
				language.English.String(): "There is no train station in Durmitor",
			},
		},
		{
			Name: "Albania",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Албания",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Tirana",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Тирана",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Shkoder",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Шкодер",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Албанией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Albania",
			},
		},
		{
			Name: "Bosnia and Herzegovina",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina",
			},
		},
		{
			Name: "Босния и Герцеговина",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina",
			},
		},
		{
			Name: "Sarajevo",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina",
			},
		},
		{
			Name: "Сараево",
			LanguageTagToCustomErrorMessageMap: map[string]string{
				language.Russian.String(): "С Боснией нет железнодорожного сообщения из Черногории",
				language.English.String(): "There are no routes between Montenegro and Bosnia and Herzegovina",
			},
		},
	}
)
