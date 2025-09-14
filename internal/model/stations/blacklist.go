package stations

import (
	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/render"
)

// don't forget to rebuild the timetable after changing these lines

var BlackListedStations = []struct {
	Names                              []string
	LanguageTagToCustomErrorMessageMap map[language.Tag]string
}{

	// summer season stations

	{
		Names: []string{"Novi sad", "Нови сад", "Indjija", "Индия", "Stara Pazova", "Стара пазова", "Nova pazova", "Нова пазова", "Subotica", "Суботица"}, // Novi Beograd - was not added to not interfere with Beograd centar
		LanguageTagToCustomErrorMessageMap: map[language.Tag]string{
			language.Russian: `
	2026 добавятся поезда по маршруту Бар - Белград - Нови Сад - Суботица.
	В остальные дни поезда из Черногории ходят только до Beograd Centar.
	`,
			language.Ukrainian: `
	2026 додадуться поїзди за маршрутом Бар - Белград - Нові Сад - Суботиця.
	В інші дні поїзди з Чорногорії ходять тільки до Beograd Centar.
	`,
			render.Belarusian: `
	2026 дададуцца цягнікі па маршруце Бар - Белград - Новы Сад - Суботыца.
	У астатнія дні цягнікі з Чарнагорыі ходзяць толькі да Beograd Centar.
	`,
			language.English: `
	In 2026 trains will be added on the route Bar - Belgrade - Novi Sad - Subotica.
	On other days, trains from Montenegro only run to Beograd Centar.
	`,
			language.German: `
	2026 werden Züge auf der Strecke Bar - Belgrad - Novi Sad - Subotica hinzugefügt.
	An anderen Tagen verkehren Züge aus Montenegro nur bis Beograd Centar.
	`,
			language.Serbian: `
	2026 биће додати возови на релацији Бар – Београд – Нови Сад – Суботица.
	Осталим данима возови из Црне Горе иду само до Beograd Centar.
	`,
			language.Croatian: `
	2026 bit će dodani vlakovi na ruti Bar - Beograd - Novi Sad - Subotica.
	Ostalim danima vlakovi iz Crne Gore voze samo do Beograd Centar.
	`,
			language.Slovak: `
	2026 budú pridané vlaky na trase Bar - Beograd - Novi Sad - Subotica.
	V ostatné dni vlaky z Čiernej Hory jazdia len do Beograd Centar.
	`,
			language.Turkish: `
	2026 tarihleri arasında Bar - Belgrad - Novi Sad - Subotica rotasında ek trenler olacak.
	Diğer günlerde Karadağ'dan trenler sadece Beograd Centar'a kadar çalışır.
	`,
		},
	},

	// summer season stations

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
}
