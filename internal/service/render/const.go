package render

import (
	"time"

	"github.com/samber/lo"
	"golang.org/x/text/language"
)

var DefaultLanguageTag = language.English

const BelarusianLanguageCode = "be"

var Belarusian = lo.Must(language.Parse(BelarusianLanguageCode)) // there is no var language.Belarusian, so we have to improvise

var SupportedLanguages = []language.Tag{
	language.Russian,
	language.Ukrainian,
	Belarusian,
	language.English,
	language.German,
	language.Serbian,
	language.Croatian,
	language.Slovak,
	language.Turkish,
}

var (
	BotInfoMap = map[language.Tag]Bot{
		language.English: {
			Name:             BotNameEn,
			Description:      BotDescriptionEn,
			ShortDescription: BotShortDescriptionEn,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartEn,
			},
		},
		language.Russian: {
			Name:             BotNameRu,
			Description:      BotDescriptionRu,
			ShortDescription: BotShortDescriptionRu,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartRu,
			},
		},
		language.Serbian: {
			Name:             BotNameSr,
			Description:      BotDescriptionSr,
			ShortDescription: BotShortDescriptionSr,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartSr,
			},
		},
		language.Turkish: {
			Name:             BotNameTr,
			Description:      BotDescriptionTr,
			ShortDescription: BotShortDescriptionTr,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartTr,
			},
		},
		language.Slovak: {
			Name:             BotNameSk,
			Description:      BotDescriptionSk,
			ShortDescription: BotShortDescriptionSk,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartSk,
			},
		},
		language.Croatian: {
			Name:             BotNameHr,
			Description:      BotDescriptionHr,
			ShortDescription: BotShortDescriptionHr,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartHr,
			},
		},
		language.German: {
			Name:             BotNameDe,
			Description:      BotDescriptionDe,
			ShortDescription: BotShortDescriptionDe,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartDe,
			},
		},
		language.Ukrainian: {
			Name:             BotNameUa,
			Description:      BotDescriptionUa,
			ShortDescription: BotShortDescriptionUa,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartUa,
			},
		},
		Belarusian: {
			Name:             BotNameBe,
			Description:      BotDescriptionBe,
			ShortDescription: BotShortDescriptionBe,
			CommandNames: map[BotCommand]string{
				BotCommandStart: BotCommandNameStartBe,
			},
		},
	}
	ErrorMessageMap = map[language.Tag]string{
		language.Russian:   ErrorMessageRu,
		language.Ukrainian: ErrorMessageUa,
		Belarusian:         ErrorMessageBe,
		language.English:   ErrorMessageEn,
		language.German:    ErrorMessageDe,
		language.Serbian:   ErrorMessageSr,
		language.Croatian:  ErrorMessageHr,
		language.Slovak:    ErrorMessageSk,
		language.Turkish:   ErrorMessageTr,
	}
	StartMessageMap = map[language.Tag]string{
		language.Russian:   StartMessageRu,
		language.Ukrainian: StartMessageUa,
		Belarusian:         StartMessageBe,
		language.English:   StartMessageEn,
		language.German:    StartMessageDe,
		language.Serbian:   StartMessageSr,
		language.Croatian:  StartMessageHr,
		language.Slovak:    StartMessageSk,
		language.Turkish:   StartMessageTr,
	}
	StationDoesNotExistMessageMap = map[language.Tag]string{
		language.Russian:   StationDoesNotExistMessageRu,
		language.Ukrainian: StationDoesNotExistMessageUa,
		Belarusian:         StationDoesNotExistMessageBe,
		language.English:   StationDoesNotExistMessageEn,
		language.German:    StationDoesNotExistMessageDe,
		language.Serbian:   StationDoesNotExistMessageSr,
		language.Croatian:  StationDoesNotExistMessageHr,
		language.Slovak:    StationDoesNotExistMessageSk,
		language.Turkish:   StationDoesNotExistMessageTr,
	}
	StationDoesNotExistMessageSuffixMap = map[language.Tag]string{
		language.Russian:   StationDoesNotExistMessageSuffixRu,
		language.Ukrainian: StationDoesNotExistMessageSuffixUa,
		Belarusian:         StationDoesNotExistMessageSuffixBe,
		language.English:   StationDoesNotExistMessageSuffixEn,
		language.German:    StationDoesNotExistMessageSuffixDe,
		language.Serbian:   StationDoesNotExistMessageSuffixSr,
		language.Croatian:  StationDoesNotExistMessageSuffixHr,
		language.Slovak:    StationDoesNotExistMessageSuffixSk,
		language.Turkish:   StationDoesNotExistMessageSuffixTr,
	}
	MonthNameMap = map[language.Tag]map[time.Month]string{
		language.Russian: {
			time.January:   "Января",
			time.February:  "Февраля",
			time.March:     "Марта",
			time.April:     "Апреля",
			time.May:       "Мая",
			time.June:      "Июня",
			time.July:      "Июля",
			time.August:    "Августа",
			time.September: "Сентября",
			time.October:   "Октября",
			time.November:  "Ноября",
			time.December:  "Декабря",
		},
		language.English: {
			time.January:   "January",
			time.February:  "February",
			time.March:     "March",
			time.April:     "April",
			time.May:       "May",
			time.June:      "June",
			time.July:      "July",
			time.August:    "August",
			time.September: "September",
			time.October:   "October",
			time.November:  "November",
			time.December:  "December",
		},
		language.Serbian: {
			time.January:   "Јануар",
			time.February:  "Фебруар",
			time.March:     "Март",
			time.April:     "Април",
			time.May:       "Мај",
			time.June:      "Јун",
			time.July:      "Јул",
			time.August:    "Август",
			time.September: "Септембар",
			time.October:   "Октобар",
			time.November:  "Новембар",
			time.December:  "Децембар",
		},
		language.Turkish: {
			time.January:   "Ocak",
			time.February:  "Şubat",
			time.March:     "Mart",
			time.April:     "Nisan",
			time.May:       "Mayıs",
			time.June:      "Haziran",
			time.July:      "Temmuz",
			time.August:    "Ağustos",
			time.September: "Eylül",
			time.October:   "Ekim",
			time.November:  "Kasım",
			time.December:  "Aralık",
		},
		language.Ukrainian: {
			time.January:   "Січня",
			time.February:  "Лютого",
			time.March:     "Березня",
			time.April:     "Квітня",
			time.May:       "Травня",
			time.June:      "Червня",
			time.July:      "Липня",
			time.August:    "Серпня",
			time.September: "Вересня",
			time.October:   "Жовтня",
			time.November:  "Листопада",
			time.December:  "Грудня",
		},
		Belarusian: {
			time.January:   "Студзеня",
			time.February:  "Лютага",
			time.March:     "Сакавіка",
			time.April:     "Красавіка",
			time.May:       "Мая",
			time.June:      "Чэрвеня",
			time.July:      "Ліпеня",
			time.August:    "Жніўня",
			time.September: "Верасня",
			time.October:   "Кастрычніка",
			time.November:  "Лістапада",
			time.December:  "Снежня",
		},
		language.German: {
			time.January:   "Januar",
			time.February:  "Februar",
			time.March:     "März",
			time.April:     "April",
			time.May:       "Mai",
			time.June:      "Juni",
			time.July:      "Juli",
			time.August:    "August",
			time.September: "September",
			time.October:   "Oktober",
			time.November:  "November",
			time.December:  "Dezember",
		},
		language.Croatian: {
			time.January:   "Siječanj",
			time.February:  "Veljača",
			time.March:     "Ožujak",
			time.April:     "Travanj",
			time.May:       "Svibanj",
			time.June:      "Lipanj",
			time.July:      "Srpanj",
			time.August:    "Kolovoz",
			time.September: "Rujan",
			time.October:   "Listopad",
			time.November:  "Studeni",
			time.December:  "Prosinac",
		},
		language.Slovak: {
			time.January:   "Január",
			time.February:  "Február",
			time.March:     "Marec",
			time.April:     "Apríl",
			time.May:       "Máj",
			time.June:      "Jún",
			time.July:      "Júl",
			time.August:    "August",
			time.September: "September",
			time.October:   "Október",
			time.November:  "November",
			time.December:  "December",
		},
	}
	ReverseRouteInlineButtonTextMap = map[language.Tag]string{
		language.Russian:   ReverseRouteInlineButtonTextRu,
		language.Ukrainian: ReverseRouteInlineButtonTextUa,
		Belarusian:         ReverseRouteInlineButtonTextBe,
		language.English:   ReverseRouteInlineButtonTextEn,
		language.German:    ReverseRouteInlineButtonTextDe,
		language.Serbian:   ReverseRouteInlineButtonTextSr,
		language.Croatian:  ReverseRouteInlineButtonTextHr,
		language.Slovak:    ReverseRouteInlineButtonTextSk,
		language.Turkish:   ReverseRouteInlineButtonTextTr,
	}
	AlertUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian:   AlertUpdateNotificationTextRu,
		language.Ukrainian: AlertUpdateNotificationTextUa,
		Belarusian:         AlertUpdateNotificationTextBe,
		language.English:   AlertUpdateNotificationTextEn,
		language.German:    AlertUpdateNotificationTextDe,
		language.Serbian:   AlertUpdateNotificationTextSr,
		language.Croatian:  AlertUpdateNotificationTextHr,
		language.Slovak:    AlertUpdateNotificationTextSk,
		language.Turkish:   AlertUpdateNotificationTextTr,
	}
	SimpleUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian:   SimpleUpdateNotificationTextRu,
		language.Ukrainian: SimpleUpdateNotificationTextUa,
		Belarusian:         SimpleUpdateNotificationTextBe,
		language.English:   SimpleUpdateNotificationTextEn,
		language.German:    SimpleUpdateNotificationTextDe,
		language.Serbian:   SimpleUpdateNotificationTextSr,
		language.Croatian:  SimpleUpdateNotificationTextHr,
		language.Slovak:    SimpleUpdateNotificationTextSk,
		language.Turkish:   SimpleUpdateNotificationTextTr,
	}
	OfficialTimetableUrlTextMap = map[language.Tag]string{
		language.Russian:   OfficialTimetableUrlTextRu,
		language.Ukrainian: OfficialTimetableUrlTextUa,
		Belarusian:         OfficialTimetableUrlTextBe,
		language.English:   OfficialTimetableUrlTextEn,
		language.German:    OfficialTimetableUrlTextDe,
		language.Serbian:   OfficialTimetableUrlTextSr,
		language.Croatian:  OfficialTimetableUrlTextHr,
		language.Slovak:    OfficialTimetableUrlTextSk,
		language.Turkish:   OfficialTimetableUrlTextTr,
	}
)

// bot general

type Bot struct {
	Name, Description, ShortDescription string
	CommandNames                        map[BotCommand]string
}

type BotCommand string

var AllCommands = []BotCommand{
	BotCommandStart,
}

const (
	BotCommandStart BotCommand = "/start"
)

const googleMapWithAllStations = "https://www.google.com/maps/d/viewer?mid=1wl76-eu79ECV5OYah2a2MKc2_9ezVME"

// English

const (
	ErrorMessageEn = "" +
		`Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic`

	StartMessageEn = "" +
		`*Montenegro Railways Timetable*

_Made together with @Leti\_deshevle_

Please enter *two stations* separated by *a comma*: 

>*Podgorica, Bijelo Polje*

Or using cyrillic:

>*Подгорица, Бијело Поље*

And I will send you the timetable:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Not sure about the correct spelling of the stations? No problem, just type them, and I will take care of the rest\.

Now it's your turn\!
`
	StationDoesNotExistMessageEn       = "This station does not exist"
	StationDoesNotExistMessageSuffixEn = "Montenegro Railway Map"
	OfficialTimetableUrlTextEn         = "More info"
	ReverseRouteInlineButtonTextEn     = "Reverse"
	AlertUpdateNotificationTextEn      = "" +
		`The timetable is already updated

From June 15th to September 16th a new train Novi Sad - Belgrade - Bar will be added

The rest of the timetable will remain exactly the same`
	SimpleUpdateNotificationTextEn = "Today's timetable is updated"

	// bot description

	BotNameEn        = "🚂 Montenegro: train timetable | Черногория расписание поезд"
	BotDescriptionEn = "" +
		`> Up-to-date timetable
> Knows every station, including Belgrade
> Can show routes between any two station, including transfer

Just type two stations with a comma:

Podgorica, Bar`
	BotShortDescriptionEn = "Up-to-date timetable with all stations and routes, including transfer and international ones, like Belgrade - Bar train"

	// bot commands

	BotCommandNameStartEn = "Start the bot"
)

// Russian

const (
	ErrorMessageRu = "" +
		`Попробуйте снова - две станции, через запятую. Вот так:

Podgorica, Niksic`

	StartMessageRu = "" +
		`*Расписание электричек Черногории*

_Сделан вместе с @Leti\_deshevle_

Пожалуйста, введите *две станции через запятую* на латинице: 

>*Podgorica, Bijelo Polje*

или кириллице: 

>*Подгорица, Бело поле*

И получите расписание:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Не уверены как правильно пишется название станции? Напишите как знаете \- мы догадаемся что вы имели ввиду\.

Теперь ваша очередь\!
`

	StationDoesNotExistMessageRu       = "Такой станции не существует"
	StationDoesNotExistMessageSuffixRu = "Карта ЖД Черногории"
	OfficialTimetableUrlTextRu         = "Подробнее"
	ReverseRouteInlineButtonTextRu     = "Обратно"
	AlertUpdateNotificationTextRu      = "" +
		`Расписание уже обновлено

С 15.06 по 16.09 добавится поезд Нови сад - Бар

В остальном расписание не изменится`
	SimpleUpdateNotificationTextRu = "Расписание на сегодня обновлено"

	// bot description

	BotNameRu        = "🚂 Черногория: расписание поездов и электричек"
	BotDescriptionRu = "" +
		`> Актуальное расписание
> Знает все станции, включая Белград
> Умеет строить маршруты с пересадкой

Просто введите две станции через запятую:

Подгорица, Бар`
	BotShortDescriptionRu = "Актуальное расписание со всеми станциями и маршрутами, включая маршруты с пересадкой и поездом Белград - Бар"

	// bot commands

	BotCommandNameStartRu = "Старт бота"
)

// Serbian

const (
	ErrorMessageSr = "" +
		`Покушајте поново - две станице, раздвојене зарезом. Као овако:

Подгорица, Бар`

	StartMessageSr = "" +
		`*Распоред железнице Црне Горе*

_Направљено уз помоћ @Leti\_deshevle_

Молим вас унесите *две станице* раздвојене *зарезом*: 

>*Podgorica, Bijelo Polje*

Iли користећи ћирилицу:

>*Подгорица, Бијело Поље*

И ја ћу вам послати распоред:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Нисте сигурни у исправном писању станица? Нема проблема, само их унесите, а ја ћу се побринути за остало\.

Сад је ваш ред\!
`

	StationDoesNotExistMessageSr       = "Ова станица не постоји"
	StationDoesNotExistMessageSuffixSr = "Karta železnica Crne Gore"
	OfficialTimetableUrlTextSr         = "Више информација"
	ReverseRouteInlineButtonTextSr     = "Обрнуто"
	AlertUpdateNotificationTextSr      = "Данашњи распоред већ приказан. Користите дугме '" + OfficialTimetableUrlTextSr + "' да бисте видели распоред за друге датуме"
	SimpleUpdateNotificationTextSr     = "Данашњи распоред је приказан"

	// bot description

	BotNameSr        = "🚂 Црна Гора: распоред возова | Crna Gora ZPCG RED VOŽNJE"
	BotDescriptionSr = "" +
		`> Ажуран распоред
> Зна сваку станицу, укључујући Београд
> Може приказати руте између било које две станице, укључујући трансфер

Просто унесите две станице раздвојене зарезом:

Подгорица, Бар`
	BotShortDescriptionSr = "Ажуран распоред са свим станицама и рутама, укључујући руте са трансфером и међународне, као што је воз Београд - Бар"

	// bot commands

	BotCommandNameStartSr = "Покрени бота"
)

// Turkish

const (
	ErrorMessageTr = "" +
		`Tekrar deneyin - iki istasyon, virgülle ayrılmış.Tam olarak şöyle:

Podgorica, Bar`

	StartMessageTr = "" +
		`*Karadağ Demiryolları Tarifesi*

_@Leti\_deshevle ile birlikte yapıldı_

Lütfen *bir virgülle ayrılmış iki istasyon* girin: 

>*Podgorica, Bijelo Polje*

Ve size tarifeyi göndereceğim:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

İstasyonların doğru yazımı konusunda emin değil misiniz? Sorun değil, sadece yazın, gerisini ben halledeceğim\.

Sıra sende\!
`

	StationDoesNotExistMessageTr       = "Bu istasyon mevcut değil"
	StationDoesNotExistMessageSuffixTr = "Karadağ Demiryolu Haritası"
	OfficialTimetableUrlTextTr         = "Daha fazla bilgi"
	ReverseRouteInlineButtonTextTr     = "Tersine"
	AlertUpdateNotificationTextTr      = "" +
		`Tarife zaten güncellendi

15 Haziran'dan 16 Eylül'e kadar Novi Sad - Belgrad - Bar yeni bir tren eklenecek

Tarifenin geri kalanı tam olarak aynı kalacak`
	SimpleUpdateNotificationTextTr = "Bugünün tarifesi güncellendi"

	// bot description

	BotNameTr        = "🚂 Karadağ: tren tarifesi | Montenegro train"
	BotDescriptionTr = "" +
		`> Güncel tarife
> Her istasyonu biliyor, Belgrad dahil
> Herhangi iki istasyon arasında rotaları gösterebilir, aktarma dahil

Sadece bir virgülle ayrılmış iki istasyonu yazın:

Podgorica, Bar`
	BotShortDescriptionTr = "Tüm istasyonlar ve rotalarla güncel tarife, aktarma rotaları ve Belgrad - Bar gibi uluslararası rotalar dahil"

	// bot commands

	BotCommandNameStartTr = "Bota başla"
)

// Belarusian

const (
	ErrorMessageBe = "" +
		`Спробуйце яшчэ - два вакзалы, праз коску. Вось так:

Podgorica, Niksic`

	StartMessageBe = "" +
		`*Расклад электрычакоў Чарнагорыі*

_Зроблена разам з @Leti\_deshevle_

Калі ласка, увядзіце *два вакзалы праз коску* на лацініцы: 

>*Podgorica, Bijelo Polje*

ці кірыліцы: 

>*Подгорица, Бело Поле*

І атрымаеце расклад:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Ня ўпэўненыя як правільна пішацца назва вакзала? Напішыце як ведаеце \- мы дагадаемся што вы мелі на ўвазе\.

Цяпер ваш чарга\!
`

	StationDoesNotExistMessageBe       = "Такога вакзала не існуе"
	StationDoesNotExistMessageSuffixBe = "Карта ЧД Чарнагорыі"
	OfficialTimetableUrlTextBe         = "Падрабязней"
	ReverseRouteInlineButtonTextBe     = "Назад "
	AlertUpdateNotificationTextBe      = "" +
		`Расклад ужо абноўлены

З 15.06 па 16.09 дадасца паезд Нові Сад - Бар

У астатнім расклад не зменіцца`
	SimpleUpdateNotificationTextBe = "Расклад на сёння абноўлены"

	// апісанне бота

	BotNameBe        = "🚂 Чарнагорыя: расклад паездаў і электрычак | Черногория поезд"
	BotDescriptionBe = "" +
		`> Актуальны расклад
> Ведае ўсе вакзалы, уключаючы Белград
> Можа будаваць маршруты з перасадкай

Проста ўвядзіце два вакзалы праз коску:

Подгорица, Бар`
	BotShortDescriptionBe = "Актуальны расклад з усімі вакзаламі і маршрутамі, уключаючы маршруты з перасадкай і паезд Белград - Бар"

	// каманды бота

	BotCommandNameStartBe = "Старт бота"
)

// Ukrainian

const (
	ErrorMessageUa = "" +
		`Спробуйте ще раз - дві станції через кому. Ось так:

Podgorica, Niksic`

	StartMessageUa = "" +
		`*Розклад електричок Чорногорії*

_Зроблено разом з @Leti\_deshevle_

Будь ласка, введіть *дві станції через кому* латиницею: 

>*Podgorica, Bijelo Polje*

або кирилицею: 

>*Подгорица, Бело Поле*

І отримайте розклад:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Не впевнені як правильно пишеться назва станції? Напишіть як знаєте \- ми здогадаємося, що ви мали на увазі\.

Тепер ваша черга\!
`

	StationDoesNotExistMessageUa       = "Такої станції не існує"
	StationDoesNotExistMessageSuffixUa = "Карта Залізниці Чорногорії"
	OfficialTimetableUrlTextUa         = "Детальніше"
	ReverseRouteInlineButtonTextUa     = "Назад"
	AlertUpdateNotificationTextUa      = "" +
		`Розклад вже оновлено

З 15.06 по 16.09 додастся поїзд Нови Сад - Бар

В іншому розклад не зміниться`
	SimpleUpdateNotificationTextUa = "Розклад на сьогодні оновлено"

	// опис бота

	BotNameUa        = "🚂 Чорногорія: розклад поїздів і електричок | Черногория поезд"
	BotDescriptionUa = "" +
		`> Актуальний розклад
> Знає всі станції, включаючи Белград
> Вміє будувати маршрути з пересадкою

Просто введіть дві станції через кому:

Подгорица, Бар`
	BotShortDescriptionUa = "Актуальний розклад з усіма станціями і маршрутами, включаючи маршрути з пересадкою і поїзд Белград - Бар"

	// команди бота

	BotCommandNameStartUa = "Старт бота"
)

// German

const (
	ErrorMessageDe = "" +
		`Versuchen Sie es erneut - zwei Bahnhöfe durch Komma getrennt. Hier ist ein Beispiel:

Podgorica, Niksic`

	StartMessageDe = "" +
		`*Fahrplan der Züge in Montenegro*

_Erstellt in Zusammenarbeit mit @Leti\_deshevle_

Bitte geben Sie *zwei Bahnhöfe durch Komma getrennt* ein, auf Lateinisch: 

>*Podgorica, Bijelo Polje*

und erhalten Sie den Fahrplan:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Sind Sie sich unsicher, wie der Bahnhof richtig geschrieben wird? Schreiben Sie, wie Sie es wissen \- wir werden verstehen, was Sie gemeint haben\.

Jetzt sind Sie dran\!
`

	StationDoesNotExistMessageDe       = "Dieser Bahnhof existiert nicht"
	StationDoesNotExistMessageSuffixDe = "Montenegro Bahnkarte"
	OfficialTimetableUrlTextDe         = "Mehr erfahren"
	ReverseRouteInlineButtonTextDe     = "Zurück"
	AlertUpdateNotificationTextDe      = "" +
		`Der Fahrplan wurde bereits aktualisiert

Vom 15.06 bis 16.09 wird ein Zug von Novi Sad nach Bar hinzugefügt

Ansonsten ändert sich der Fahrplan nicht`
	SimpleUpdateNotificationTextDe = "Der Fahrplan für heute wurde aktualisiert"

	// Bot-Beschreibung

	BotNameDe        = "🚂 Montenegro: Zug- und Zugfahrplan | train timetable"
	BotDescriptionDe = "" +
		`> Aktueller Fahrplan
> Kennt alle Bahnhöfe, einschließlich Belgrad
> Kann Routen mit Umstieg erstellen

Geben Sie einfach zwei Bahnhöfe durch Komma getrennt ein:

Podgorica, Bar`
	BotShortDescriptionDe = "Aktueller Fahrplan mit allen Bahnhöfen und Routen, einschließlich Routen mit Umstieg und dem Zug Belgrad - Bar"

	// Bot-Befehle

	BotCommandNameStartDe = "Start des Bots"
)

// Croatian

const (
	ErrorMessageHr = "" +
		`Pokušajte ponovno - dva kolodvora odvojena zarezom. Evo primjera:

Podgorica, Nikšić`

	StartMessageHr = "" +
		`*Raspored vlakova Crne Gore*

_Izrađeno u suradnji s @Leti\_deshevle_

Molimo unesite *dva kolodvora odvojena zarezom* na latinici: 

>*Podgorica, Bijelo Polje*

ili na ćirilici: 

>*Подгорица, Бијело поле*

i dobit ćete raspored:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Niste sigurni kako se pravilno piše naziv kolodvora? Napišite kako znate \- razumjet ćemo što ste mislili\.

Sada je vaš red\!
`

	StationDoesNotExistMessageHr       = "Taj kolodvor ne postoji"
	StationDoesNotExistMessageSuffixHr = "Karta željeznica Crne Gore"
	OfficialTimetableUrlTextHr         = "Više informacija"
	ReverseRouteInlineButtonTextHr     = "Natrag"
	AlertUpdateNotificationTextHr      = "" +
		`Raspored je već ažuriran

Od 15.06. do 16.09. bit će dodan vlak Novi Sad - Bar

Inače, raspored se ne mijenja`
	SimpleUpdateNotificationTextHr = "Raspored za danas je ažuriran"

	// Opis bota

	BotNameHr        = "🚂 Crna Gora: raspored vlakova | ZPCG RED VOŽNJE"
	BotDescriptionHr = "" +
		`> Trenutni raspored
> Zna sve kolodvore, uključujući Beograd
> Može graditi rute s presjedanjem

Jednostavno unesite dva kolodvora odvojena zarezom:

Podgorica, Bar`
	BotShortDescriptionHr = "Trenutni raspored sa svim kolodvorima i rutama, uključujući rute s presjedanjem i vlak Beograd - Bar"

	// Naredbe bota

	BotCommandNameStartHr = "Pokreni bota"
)

// Slovak

const (
	ErrorMessageSk = "" +
		`Skúste to znova - dva stanice oddelené čiarkou. Tu je príklad:

Podgorica, Nikšić`

	StartMessageSk = "" +
		`*Cestovný poriadok vlakov Čiernej Hory*

_Vytvorené v spolupráci s @Leti\_deshevle_

Prosím, zadajte *dve stanice oddelené čiarkou* v latinčine: 

>*Podgorica, Bijelo Polje*

alebo v cyrilike: 

>*Подгорица, Бијело поле*

a dostanete cestovný poriadok:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Nie ste si istí, ako sa správne píše názov stanice? Napíšte to, ako viete \- pochopíme, čo ste mali na mysli\.

Teraz je na vás\!
`

	StationDoesNotExistMessageSk       = "Taká stanica neexistuje"
	StationDoesNotExistMessageSuffixSk = "Mapa železníc Čiernej Hory"
	OfficialTimetableUrlTextSk         = "Viac informácií"
	ReverseRouteInlineButtonTextSk     = "Späť"
	AlertUpdateNotificationTextSk      = "" +
		`Cestovný poriadok už bol aktualizovaný

Od 15.06 do 16.09 bude pridaný vlak Novi Sad - Bar

Inak sa poradie nezmení`
	SimpleUpdateNotificationTextSk = "Cestovný poriadok pre dnešok bol aktualizovaný"

	// Popis bota

	BotNameSk        = "🚂 Čierna Hora: cestovný poriadok vlakov | Montenegro train"
	BotDescriptionSk = "" +
		`> Aktuálny cestovný poriadok
> Pozná všetky stanice, vrátane Belehradu
> Vie zostaviť trasy s prestupom

Jednoducho zadajte dve stanice oddelené čiarkou:

Podgorica, Bar`
	BotShortDescriptionSk = "Aktuálny cestovný poriadok so všetkými stanicami a trasami, vrátane trás s prestupom a vlaku Belehrad - Bar"

	// Príkazy bota

	BotCommandNameStartSk = "Spustiť bota"
)
