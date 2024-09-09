package timetable

type StationAliases struct {
	StationName string
	Aliases     []string
}

var AliasesStationsList = []StationAliases{
	{
		StationName: "Beograd Centar",
		Aliases:     []string{"Beograd", "belgrad", "Београд Центар", "Белград Центр"},
	},
	{
		StationName: "Novi Sad",
		Aliases:     []string{"Нови сад", "Новый Сад"},
	},
	{
		StationName: "Novi Beograd",
		Aliases:     []string{"Novi Beograd", "Novi belgrad", "New Belgrade", "Нови Београд", "Новый Белград"},
	},
	{
		StationName: "Stara Pazova",
		Aliases:     []string{"Стара Пазова"},
	},
	{
		StationName: "Nova Pazova",
		Aliases:     []string{"Нова Пазова"},
	},
	{
		StationName: "Bar",
		Aliases:     []string{"Бар"},
	},
	{
		StationName: "Sutomore",
		Aliases:     []string{"Сутоморе"},
	},
	{
		StationName: "Golubovci",
		Aliases:     []string{"Голубовци", "Голубовцы"},
	},
	{
		StationName: "Podgorica",
		Aliases:     []string{"Подгорица"},
	},
	{
		StationName: "Kolašin",
		Aliases:     []string{"Колашин"},
	},
	{
		StationName: "Mojkovac",
		Aliases:     []string{"Мојковац", "Мойковац"},
	},
	{
		StationName: "Bijelo Polje",
		Aliases:     []string{"Бијело Поље", "Бело Поле"},
	},
	{
		StationName: "Prijepolje teretna",
		Aliases:     []string{"Пријепоље теретна", "Приеполье теретна"},
	},
	{
		StationName: "Prijepolje",
		Aliases:     []string{"Пријепоље", "Приеполье"},
	},
	{
		StationName: "Priboj",
		Aliases:     []string{"Прибој", "Прибой"},
	},
	{
		StationName: "Užice",
		Aliases:     []string{"Ужице"},
	},
	{
		StationName: "Požega",
		Aliases:     []string{"Пожега"},
	},
	{
		StationName: "Kosjerić",
		Aliases:     []string{"Косјерић", "Косерич"},
	},
	{
		StationName: "Valjevo",
		Aliases:     []string{"Ваљево", "Вальево"},
	},
	{
		StationName: "Lajkovac",
		Aliases:     []string{"Лајковац", "Лайковац"},
	},
	{
		StationName: "Lazarevac",
		Aliases:     []string{"Лазаревац", "Лазаревац"},
	},
	{
		StationName: "Rakovica",
		Aliases:     []string{"Раковица", "Раковица"},
	},
	{
		StationName: "Zemun",
		Aliases:     []string{"Земун"},
	},
	{
		StationName: "Šušanj",
		Aliases:     []string{"Шушањ", "Шушань"},
	},
	{
		StationName: "Crmnica",
		Aliases:     []string{"Црмница"},
	},
	{
		StationName: "Virpazar",
		Aliases:     []string{"Вирпазар"},
	},
	{
		StationName: "Vranjina",
		Aliases:     []string{"Врањина", "Враньина"},
	},
	{
		StationName: "Zeta",
		Aliases:     []string{"Зета"},
	},
	{
		StationName: "Morača",
		Aliases:     []string{"Морача"},
	},
	{
		StationName: "Aerodrom",
		Aliases:     []string{"аэродром"},
	},
	{
		StationName: "Zlatica",
		Aliases:     []string{"Златица"},
	},
	{
		StationName: "Bioče",
		Aliases:     []string{"Биоче"},
	},
	{
		StationName: "Bratonožići",
		Aliases:     []string{"Братоношићи", "Братоношичи"},
	},
	{
		StationName: "Lutovo",
		Aliases:     []string{"Лутово"},
	},
	{
		StationName: "Kruševački Potok",
		Aliases:     []string{"Крушевачки поток", "Крушевацки поток"},
	},
	{
		StationName: "Trebešica",
		Aliases:     []string{"Требешица"},
	},
	{
		StationName: "Selište",
		Aliases:     []string{"Селиште"},
	},
	{
		StationName: "Kos",
		Aliases:     []string{"Кос"},
	},
	{
		StationName: "Mateševo",
		Aliases:     []string{"Матешево", "Матесево"},
	},
	{
		StationName: "Padež",
		Aliases:     []string{"Падеж"},
	},
	{
		StationName: "Oblutak",
		Aliases:     []string{"Облутак"},
	},
	{
		StationName: "Trebaljevo",
		Aliases:     []string{"Требаљево", "Требальево"},
	},
	{
		StationName: "Štitarička Rijeka",
		Aliases:     []string{"Штитарица река", "Штитаричка река"}, // TODO: recheck
	},
	{
		StationName: "Žari",
		Aliases:     []string{"Зари", "Зари"}, // TODO: recheck
	},
	{
		StationName: "Mijatovo Kolo",
		Aliases:     []string{"Мијатово коло", "Миятово коло"},
	},
	{
		StationName: "Slijepač Most",
		Aliases:     []string{"Слијепац мост", "Слепец мост"}, // TODO: recheck
	},
	{
		StationName: "Ravna Rijeka",
		Aliases:     []string{"Равна ријека", "Равна река"},
	},
	{
		StationName: "Kruševo",
		Aliases:     []string{"Крушево"},
	},
	{
		StationName: "Lješnica",
		Aliases:     []string{"Љешница", "Льешница"},
	},
	{
		StationName: "Pričelje",
		Aliases:     []string{"Прицеље", "Причелье"},
	},
	{
		StationName: "Spuž",
		Aliases:     []string{"Спуж"},
	},
	{
		StationName: "Ljutotuk",
		Aliases:     []string{"Љутотук", "Лютотук"},
	},
	{
		StationName: "Danilovgrad",
		Aliases:     []string{"Даниловград"},
	},
	{
		StationName: "Slap",
		Aliases:     []string{"Слап"},
	},
	{
		StationName: "Bare Šumanovića",
		Aliases:     []string{"Баре Шумановица"},
	},
	{
		StationName: "Šobajići",
		Aliases:     []string{"Собајићи", "Шобайичи"}, // TODO: recheck
	},
	{
		StationName: "Ostrog",
		Aliases:     []string{"Острог"},
	},
	{
		StationName: "Dabovići",
		Aliases:     []string{"Дабовићи", "Дабовичи"},
	},
	{
		StationName: "Stubica",
		Aliases:     []string{"Стубица"},
	},
	{
		StationName: "Nikšić",
		Aliases:     []string{"Никшић", "Никшич"},
	},
}
