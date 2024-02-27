package render

import (
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
)

// Default

var DefaultLanguageTag = language.English

const (
	StartMessageDefault = StartMessageEn
	ErrorMessageDefault = ErrorMessageEn
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
)

// Russian

const (
	ErrorMessageRu = `Попробуйте снова - две станции на латинице, через запятую. Вот так:

Podgorica, Niksic`

	StartMessageRu = "" +
		"*Расписание электричек Черногории*\n" +
		"\n" +
		"_Сделан вместе с @Leti\\_deshevle_\n" +
		"\n" +
		"Пожалуйста, введите *две станции через запятую на латинице*: \n" +
		"\n" +
		">*Podgorica, Bijelo Polje*\n" +
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
	StationDoesNotExistMessageSuffixRu = "" // TODO: add /info "Хотите узнать где в Черногории есть жд сообщение? Используйте команду /info"
)
