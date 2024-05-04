package render

import (
	"time"

	"golang.org/x/text/language"
)

var SupportedLanguages = []language.Tag{
	language.Russian,
	language.English,
}

var (
	ErrorMessageMap = map[language.Tag]string{
		language.Russian: ErrorMessageRu,
		language.English: ErrorMessageEn,
	}
	StartMessageMap = map[language.Tag]string{
		language.Russian: StartMessageRu,
		language.English: StartMessageEn,
	}
	StationDoesNotExistMessageMap = map[language.Tag]string{
		language.Russian: StationDoesNotExistMessageRu,
		language.English: StationDoesNotExistMessageEn,
	}
	StationDoesNotExistMessageSuffixMap = map[language.Tag]string{
		language.Russian: StationDoesNotExistMessageSuffixRu,
		language.English: StationDoesNotExistMessageSuffixEn,
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
	}
	ReverseRouteInlineButtonTextMap = map[language.Tag]string{
		language.Russian: ReverseRouteInlineButtonTextRu,
		language.English: ReverseRouteInlineButtonTextEn,
	}
	AlertUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian: AlertUpdateNotificationTextRu,
		language.English: AlertUpdateNotificationTextEn,
	}
	SimpleUpdateNotificationTextMap = map[language.Tag]string{
		language.Russian: SimpleUpdateNotificationTextRu,
		language.English: SimpleUpdateNotificationTextEn,
	}
)

// bot general
const (
	BotCommandStart = "/start"
)

// Default

var DefaultLanguageTag = language.English

const (
	StartMessageDefault      = StartMessageEn
	ErrorMessageDefault      = ErrorMessageEn
	OfficialTimetableUrlText = OfficialTimetableUrlTextEn
)

// English

const (
	ErrorMessageEn = `Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic`

	StartMessageEn = "" +
		"*Montenegro Railways Timetable*\n" +
		"\n" +
		"_Made together with @Leti\\_deshevle_\n" +
		"\n" +
		"Please enter *two stations* separated by *a comma*: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
		"\n" +
		"Or using cyrillic:\n" +
		"\n" +
		">*Подгорица, Бијело Поље*\n" +
		"\n" +
		"And I will send you the timetable:\n" +
		"\n" +
		">Podgorica \\> Bijelo Polje\n" +
		">[6100](https://zpcg.me/details?timetable=41)  `06:20 \\> 08:38` \n" +
		">\\.\\.\\.\n" +
		"\n" +
		"Not sure about the correct spelling of the stations? No problem, just type them, and I will take care of the rest\\.\n" +
		"\n" +
		"Now it's your turn\\!"

	StationDoesNotExistMessageEn       = "This station does not exist"
	StationDoesNotExistMessageSuffixEn = " " // TODO: add /info "Would you like to know more about available train stations in Montenegro? Check the /info command"
	OfficialTimetableUrlTextEn         = "More info"
	ReverseRouteInlineButtonTextEn     = "Reverse"
	AlertUpdateNotificationTextEn      = "Today's timetable is already shown. Use '" + OfficialTimetableUrlTextEn + "' button to see timetable for other dates"
	SimpleUpdateNotificationTextEn     = "Today's timetable is shown"

	// bot description

	BotNameEn             = "🚂 Montenegro: train timetable bot"
	BotDescriptionEn      = ""
	BotShortDescriptionEn = ""

	// bot commands

	BotCommandNameStartEn = "Start the bot"
)

// Russian

const (
	ErrorMessageRu = `Попробуйте снова - две станции, через запятую. Вот так:

Podgorica, Niksic`

	StartMessageRu = "" +
		"*Расписание электричек Черногории*\n" +
		"\n" +
		"_Сделан вместе с @Leti\\_deshevle_\n" +
		"\n" +
		"Пожалуйста, введите *две станции через запятую* на латинице: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
		"\n" +
		"или кириллице: \n" +
		"\n" +
		">*Подгорица, Бело поле*\n" +
		"\n" +
		"И получите расписание:\n" +
		"\n" +
		">Podgorica \\> Bijelo Polje\n" +
		">[6100](https://zpcg.me/details?timetable=41)  `06:20 \\> 08:38` \n" +
		">\\.\\.\\.\n" +
		"\n" +
		"Не уверены как правильно пишется название станции? Напишите как знаете \\- мы догадаемся что вы имели ввиду\\.\n" +
		"\n" +
		"Теперь ваша очередь\\!"

	StationDoesNotExistMessageRu       = "Такой станции не существует"
	StationDoesNotExistMessageSuffixRu = "  " // TODO: add /info "Хотите узнать где в Черногории есть жд сообщение? Используйте команду /info"
	OfficialTimetableUrlTextRu         = "Подробнее"
	ReverseRouteInlineButtonTextRu     = "Обратно"
	AlertUpdateNotificationTextRu      = "Расписание на сегодня обновлено. Нажмите '" + OfficialTimetableUrlTextRu + "', чтобы увидеть расписание на другие дни"
	SimpleUpdateNotificationTextRu     = "Расписание на сегодня показано"

	// bot description

	BotNameRu             = "🚂 Montenegro: train timetable bot"
	BotDescriptionRu      = ""
	BotShortDescriptionRu = ""

	// bot commands

	BotCommandNameStartRu = "Start the bot"
)

// Montenegrin

const (
	ErrorMessageMn = `Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic`

	StartMessageMn = "" +
		"*Montenegro Railways Timetable*\n" +
		"\n" +
		"_Made together with @Leti\\_deshevle_\n" +
		"\n" +
		"Please enter *two stations* separated by *a comma*: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
		"\n" +
		"Or using cyrillic:\n" +
		"\n" +
		">*Подгорица, Бијело Поље*\n" +
		"\n" +
		"And I will send you the timetable:\n" +
		"\n" +
		">Podgorica \\> Bijelo Polje\n" +
		">[6100](https://zpcg.me/details?timetable=41)  `06:20 \\> 08:38` \n" +
		">\\.\\.\\.\n" +
		"\n" +
		"Not sure about the correct spelling of the stations? No problem, just type them, and I will take care of the rest\\.\n" +
		"\n" +
		"Now it's your turn\\!"

	StationDoesNotExistMessageMn       = "This station does not exist"
	StationDoesNotExistMessageSuffixMn = " " // TODO: add /info "Would you like to know more about available train stations in Montenegro? Check the /info command"
	OfficialTimetableUrlTextMn         = "More info"
	ReverseRouteInlineButtonTextMn     = "Reverse"
	AlertUpdateNotificationTextMn      = "Today's timetable is already shown. Use '" + OfficialTimetableUrlTextEn + "' button to see timetable for other dates"
	SimpleUpdateNotificationTextMn     = "Today's timetable is shown"

	// bot description

	BotNameMn             = "🚂 Montenegro: train timetable bot"
	BotDescriptionMn      = ""
	BotShortDescriptionMn = ""

	// bot commands

	BotCommandNameStartMn = "Start the bot"
)

// Serbian

const (
	ErrorMessageSr = `Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic`

	StartMessageSr = "" +
		"*Montenegro Railways Timetable*\n" +
		"\n" +
		"_Made together with @Leti\\_deshevle_\n" +
		"\n" +
		"Please enter *two stations* separated by *a comma*: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
		"\n" +
		"Or using cyrillic:\n" +
		"\n" +
		">*Подгорица, Бијело Поље*\n" +
		"\n" +
		"And I will send you the timetable:\n" +
		"\n" +
		">Podgorica \\> Bijelo Polje\n" +
		">[6100](https://zpcg.me/details?timetable=41)  `06:20 \\> 08:38` \n" +
		">\\.\\.\\.\n" +
		"\n" +
		"Not sure about the correct spelling of the stations? No problem, just type them, and I will take care of the rest\\.\n" +
		"\n" +
		"Now it's your turn\\!"

	StationDoesNotExistMessageSr       = "This station does not exist"
	StationDoesNotExistMessageSuffixSr = " " // TODO: add /info "Would you like to know more about available train stations in Montenegro? Check the /info command"
	OfficialTimetableUrlTextSr         = "More info"
	ReverseRouteInlineButtonTextSr     = "Reverse"
	AlertUpdateNotificationTextSr      = "Today's timetable is already shown. Use '" + OfficialTimetableUrlTextEn + "' button to see timetable for other dates"
	SimpleUpdateNotificationTextSr     = "Today's timetable is shown"

	// bot description

	BotNameSr             = "🚂 Montenegro: train timetable bot"
	BotDescriptionSr      = ""
	BotShortDescriptionSr = ""

	// bot commands

	BotCommandNameStartSr = "Start the bot"
)

// Turkish

const (
	ErrorMessageTr = `Try again - two stations, separated by a comma. Just like that:

Podgorica, Niksic`

	StartMessageTr = "" +
		"*Montenegro Railways Timetable*\n" +
		"\n" +
		"_Made together with @Leti\\_deshevle_\n" +
		"\n" +
		"Please enter *two stations* separated by *a comma*: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
		"\n" +
		"Or using cyrillic:\n" +
		"\n" +
		">*Подгорица, Бијело Поље*\n" +
		"\n" +
		"And I will send you the timetable:\n" +
		"\n" +
		">Podgorica \\> Bijelo Polje\n" +
		">[6100](https://zpcg.me/details?timetable=41)  `06:20 \\> 08:38` \n" +
		">\\.\\.\\.\n" +
		"\n" +
		"Not sure about the correct spelling of the stations? No problem, just type them, and I will take care of the rest\\.\n" +
		"\n" +
		"Now it's your turn\\!"

	StationDoesNotExistMessageTr       = "This station does not exist"
	StationDoesNotExistMessageSuffixTr = " " // TODO: add /info "Would you like to know more about available train stations in Montenegro? Check the /info command"
	OfficialTimetableUrlTextTr         = "More info"
	ReverseRouteInlineButtonTextTr     = "Reverse"
	AlertUpdateNotificationTextTr      = "Today's timetable is already shown. Use '" + OfficialTimetableUrlTextEn + "' button to see timetable for other dates"
	SimpleUpdateNotificationTextTr     = "Today's timetable is shown"

	// bot description

	BotNameTr             = "🚂 Montenegro: train timetable bot"
	BotDescriptionTr      = ""
	BotShortDescriptionTr = ""

	// bot commands

	BotCommandNameStartTr = "Start the bot"
)
