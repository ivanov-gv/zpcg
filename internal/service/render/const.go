package render

import (
	"time"

	"golang.org/x/text/language"
)

var DefaultLanguageTag = language.English

var SupportedLanguages = []language.Tag{
	language.Russian,
	language.English,
}

var (
	CommandMap = map[language.Tag]Bot{
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
	}
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
	OfficialTimetableUrlTextMap = map[language.Tag]string{
		language.Russian: OfficialTimetableUrlTextRu,
		language.English: OfficialTimetableUrlTextEn,
	}
)

// bot general

type Bot struct {
	Name, Description, ShortDescription string
	CommandNames                        map[BotCommand]string
}

type BotCommand string

var allCommands = []BotCommand{
	BotCommandStart,
}

const (
	BotCommandStart BotCommand = "/start"
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

	BotNameEn             = "🚂 Montenegro: train timetable | Черногория расписание поезд"
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

	BotNameRu             = "🚂 Черногория: расписание поездов и электричек"
	BotDescriptionRu      = ""
	BotShortDescriptionRu = ""

	// bot commands

	BotCommandNameStartRu = "Start the bot"
)
