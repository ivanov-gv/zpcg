package blacklist

import (
	"github.com/samber/lo"
	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	"github.com/ivanov-gv/zpcg/internal/service/name"
	"github.com/ivanov-gv/zpcg/internal/service/render"
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
		lo.Map(_blackListedStations, func(item struct {
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

	_blackListedStations = []struct {
		Names                              []string
		LanguageTagToCustomErrorMessageMap map[language.Tag]string
	}{
		{
			Names: []string{"Budva", "Будва"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Будве нет жд станций",
				language.Ukrainian: "В Будві немає залізничних станцій",
				render.Belarusian:  "У Будве няма жд-станцый",
				language.English:   "There is no train station in Budva",
				language.German:    "Es gibt keinen Bahnhof in Budva",
				language.Serbian:   "У Будви нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Budvi",
				language.Slovak:    "V Budve nie je železničná stanica",
				language.Turkish:   "Budva'da tren istasyonu yok",
			},
		},
		{
			Names: []string{"Tivat", "Тиват"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Тивате нет жд станций",
				language.Ukrainian: "В Тіваті немає залізничних станцій",
				render.Belarusian:  "У Тівате няма жд-станцый",
				language.English:   "There is no train station in Tivat",
				language.German:    "Es gibt keinen Bahnhof in Tivat",
				language.Serbian:   "У Тивату нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Tivtu",
				language.Slovak:    "V Tivati nie je železničná stanica",
				language.Turkish:   "Tivat'ta tren istasyonu yok",
			},
		},
		{
			Names: []string{"Kotor", "Котор"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Которе нет жд станций",
				language.Ukrainian: "В Которі немає залізничних станцій",
				render.Belarusian:  "У Котары няма жд-станцый",
				language.English:   "There is no train station in Kotor",
				language.German:    "Es gibt keinen Bahnhof in Kotor",
				language.Serbian:   "У Котору нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Kotoru",
				language.Slovak:    "V Kotore nie je železničná stanica",
				language.Turkish:   "Kotor'da tren istasyonu yok",
			},
		},
		{
			Names: []string{"Cetinje", "Цетине"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Цетине нет жд станций",
				language.Ukrainian: "В Цетинє немає залізничних станцій",
				render.Belarusian:  "У Цетыне няма жд-станцый",
				language.English:   "There is no train station in Cetinje",
				language.German:    "Es gibt keinen Bahnhof in Cetinje",
				language.Serbian:   "У Цетину нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Cetinju",
				language.Slovak:    "V Cetinje nie je železničná stanica",
				language.Turkish:   "Cetinje'de tren istasyonu yok",
			},
		},
		{
			Names: []string{"Perast", "Пераст"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Перасте нет жд станций",
				language.Ukrainian: "В Перасті немає залізничних станцій",
				render.Belarusian:  "У Перасце няма жд-станцый",
				language.English:   "There is no train station in Perast",
				language.German:    "Es gibt keinen Bahnhof in Perast",
				language.Serbian:   "У Перасту нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Perastu",
				language.Slovak:    "V Peraste nie je železničná stanica",
				language.Turkish:   "Perast'ta tren istasyonu yok",
			},
		},
		{
			Names: []string{"Durmitor", "Дурмитор"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Дурмиторе нет жд станций",
				language.Ukrainian: "В Дурміторі немає залізничних станцій",
				render.Belarusian:  "У Дурміторы няма жд-станцый",
				language.English:   "There is no train station in Durmitor",
				language.German:    "Es gibt keinen Bahnhof in Durmitor",
				language.Serbian:   "У Дурмитору нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Durmitoru",
				language.Slovak:    "V Durmitore nie je železničná stanica",
				language.Turkish:   "Durmitor'da tren istasyonu yok",
			},
		},
		{
			Names: []string{"Petrovac", "Петровац"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Петроваце нет жд станций",
				language.Ukrainian: "В Петроваці немає залізничних станцій",
				render.Belarusian:  "У Петроваце няма жд-станцый",
				language.English:   "There is no train station in Petrovac",
				language.German:    "Es gibt keinen Bahnhof in Petrovac",
				language.Serbian:   "У Петровцу нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Petrovcu",
				language.Slovak:    "V Petrovci nie je železničná stanica",
				language.Turkish:   "Petrovac'ta tren istasyonu yok",
			},
		},
		{
			Names: []string{"Ulcinj", "Ульцин"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Ульцине нет жд станций",
				language.Ukrainian: "В Ульцині немає залізничних станцій",
				render.Belarusian:  "У Ульціне няма жд-станцый",
				language.English:   "There is no train station in Ulcinj",
				language.German:    "Es gibt keinen Bahnhof in Ulcinj",
				language.Serbian:   "У Улцињу нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Ulcinju",
				language.Slovak:    "V Ulcinji nie je železničná stanica",
				language.Turkish:   "Ulcinj'de tren istasyonu yok",
			},
		},
		{
			Names: []string{"Sveti Stefan", "Свети Стефан"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Свети Стефан нет жд станций",
				language.Ukrainian: "В Святого Стефана немає залізничних станцій",
				render.Belarusian:  "У Святога Стэфана няма жд-станцый",
				language.English:   "There is no train station in Sveti Stefan",
				language.German:    "Es gibt keinen Bahnhof in Sveti Stefan",
				language.Serbian:   "У Светог Стефана нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Svetom Stefanu",
				language.Slovak:    "V Svetom Štefanovi nie je železničná stanica",
				language.Turkish:   "Sveti Stefan'da tren istasyonu yok",
			},
		},
		{
			Names: []string{"Becici", "Бечичи"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Бечичи нет жд станций",
				language.Ukrainian: "В Бечічі немає залізничних станцій",
				render.Belarusian:  "У Бечыцы няма жд-станцый",
				language.English:   "There is no train station in Becici",
				language.German:    "Es gibt keinen Bahnhof in Becici",
				language.Serbian:   "У Бечићима нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Becićima",
				language.Slovak:    "V Bečičiach nie je železničná stanica",
				language.Turkish:   "Becici'de tren istasyonu yok",
			},
		},
		{
			Names: []string{"Herceg Novi", "Херцег Нови"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Херцег-Нови нет жд станций",
				language.Ukrainian: "У Герцег-Новому немає залізничних станцій",
				render.Belarusian:  "У Герцага-Новым няма жд-станцый",
				language.English:   "There is no train station in Herceg Novi",
				language.German:    "Es gibt keinen Bahnhof in Herceg Novi",
				language.Serbian:   "У Херцег Новом нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Herceg Novom",
				language.Slovak:    "V Herceg Novom nie je železničná stanica",
				language.Turkish:   "Herceg Novi'de tren istasyonu yok",
			},
		},
		{
			Names: []string{"Savnik", "Шавник"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Шавнике нет жд станций",
				language.Ukrainian: "В Шавніку немає залізничних станцій",
				render.Belarusian:  "У Шаўнік няма жд-станцый",
				language.English:   "There is no train station in Savnik",
				language.German:    "Es gibt keinen Bahnhof in Savnik",
				language.Serbian:   "У Шавнику нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Savniku",
				language.Slovak:    "V Šavníku nie je železničná stanica",
				language.Turkish:   "Savnik'te tren istasyonu yok",
			},
		},
		{
			Names: []string{"Zabljak", "Жабляк"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "В Жабляке нет жд станций",
				language.Ukrainian: "У Жабляку немає залізничних станцій",
				render.Belarusian:  "У Жабляку няма жд-станцый",
				language.English:   "There is no train station in Zabljak",
				language.German:    "Es gibt keinen Bahnhof in Zabljak",
				language.Serbian:   "У Жабљаку нема железничке станице",
				language.Croatian:  "Nema željezničke stanice u Žabljaku",
				language.Slovak:    "V Žabľaku nie je železničná stanica",
				language.Turkish:   "Zabljak'ta tren istasyonu yok",
			},
		},
		{
			Names: []string{"Albania", "Албания", "Tirana", "Тирана", "Shkoder", "Шкодер"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "С Албанией нет железнодорожного сообщения из Черногории",
				language.Ukrainian: "З Чорногорії немає залізничного сполучення з Албанією",
				render.Belarusian:  "З Чарнагорыі няма жадзіловага сувязі з Албаніяй",
				language.English:   "There are no routes between Montenegro and Albania",
				language.German:    "Es gibt keine Verbindungen zwischen Montenegro und Albanien",
				language.Serbian:   "Нема маршрута између Црне Горе и Албаније",
				language.Croatian:  "Nema veze između Crne Gore i Albanije",
				language.Slovak:    "Medzi Čiernou Horou a Albánskom nie sú žiadne trasy",
				language.Turkish:   "Karadağ ile Arnavutluk arasında rota yok",
			},
		},
		{
			Names: []string{"Bosnia and Herzegovina", "Босния и Герцеговина", "Sarajevo", "Сараево"},
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian:   "С Боснией нет железнодорожного сообщения из Черногории",
				language.Ukrainian: "З Чорногорії немає залізничного сполучення з Боснією і Герцеговиною",
				render.Belarusian:  "З Чарнагорыі няма жадзіловага сувязі з Боснія і Герцагавіна",
				language.English:   "There are no routes between Montenegro and Bosnia and Herzegovina",
				language.German:    "Es gibt keine Verbindungen zwischen Montenegro und Bosnien und Herzegowina",
				language.Serbian:   "Нема маршрута између Црне Горе и Босне и Херцеговине",
				language.Croatian:  "Nema veze između Crne Gore i Bosne i Hercegovine",
				language.Slovak:    "Medzi Čiernou Horou a Bosnou a Hercegovinou nie sú žiadne trasy",
				language.Turkish:   "Karadağ ile Bosna-Hersek arasında rota yok",
			},
		},
		{
			Names: []string{"Novi sad", "Нови сад", "Indjija", "Индия", "Stara Pazova", "Стара пазова", "Nova pazova", "Нова пазова"}, // Novi Beograd - was not added to not interfere with Beograd centar
			LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
				language.Russian: `
Поезда 1130 и 1131 с 15.06 по 16.09 будут ходить по продленному маршруту Бар - Белград - Нови Сад. 
В остальные дни поезда из Черногории до Нови Сада не ходят.

Ссылки на эти маршруты в расписании:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.Ukrainian: `
Поїзди 1130 і 1131 з 15.06 по 16.09 будуть курсувати за продовженим маршрутом Бар - Белград - Нови Сад. 
У інші дні поїзди з Чорногорії до Нови Сада не курсують.

Посилання на ці маршрути в розкладі:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				render.Belarusian: `
Паезды 1130 і 1131 з 15.06 па 16.09 будуць курсаваць па прадоўжаным маршруце Бар - Белград - Новы Сад. 
У іншыя дні паезды з Чарнагорыі да Новага Сада не курсуюць.

Спасылкі на гэтыя маршруты ў раскладзе:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.English: `
Trains 1130 and 1131 from 15.06 to 16.09 will operate on an extended route Bar - Belgrade - Novi Sad. 
On other days, trains from Montenegro to Novi Sad do not operate.

Links to these routes in the schedule:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.German: `
Die Züge 1130 und 1131 verkehren vom 15.06. bis 16.09. auf einer verlängerten Route Bar - Belgrad - Novi Sad. 
An anderen Tagen verkehren die Züge von Montenegro nach Novi Sad nicht.

Links zu diesen Routen im Fahrplan:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.Serbian: `
Возови 1130 и 1131 од 15.06. до 16.09. саобраћаће на продуженом правцу Бар – Београд – Нови Сад.
Осталим данима возови из Црне Горе за Нови Сад не саобраћају.

Линкови до ових рута у распореду:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.Croatian: `
Vlakovi 1130 i 1131 od 15.06. do 16.09. će voziti na produženoj ruti Bar - Beograd - Novi Sad. 
U ostalim danima, vlakovi iz Crne Gore za Novi Sad ne voze.

Poveznice na ove rute u rasporedu:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.Slovak: `
Vlaky 1130 a 1131 od 15.06. do 16.09. budú jazdiť na predĺženej trase Bar - Beograd - Novi Sad. 
V ostatné dni vlaky z Čiernogorska do Nového Sadu nejazdia.

Odkazy na tieto trasy v rozvrhu:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
				language.Turkish: `
15.06 - 16.09 tarihleri arasında trenler 1130 ve 1131 Bar - Belgrad - Novi Sad üzerinde uzatılmış bir rota üzerinde işleyecek. 
Diğer günlerde Karadağ'dan Novi Sad'a tren seferleri yapılmamaktadır.

Bu rotalara ilişkin bağlantılar:
https://zpcg.me/red-voznje?start=Bar&finish=Novi+Sad&date=2024-06-16`,
			},
		},
	}
)
