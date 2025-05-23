package be

import "time"

const (
	ErrorMessage = "" +
		`Спробуйце яшчэ - два вакзалы, праз коску. Вось так:

Podgorica, Niksic`

	StationDoesNotExistMessage   = "Такога вакзала не існуе"
	RailwayMapButtonTextMap      = "Карта ЧД Чарнагорыі"
	OfficialTimetableUrlText     = "Падрабязней"
	ReverseRouteInlineButtonText = "Назад "
	AlertUpdateNotificationText  = "" +
		`Расклад ужо абноўлены

З 13.06.2025 по 14.09.2025 дадасца паезд Суботица - Бар

У астатнім расклад не зменіцца`
	SimpleUpdateNotificationText = "Расклад на сёння абноўлены"

	// апісанне бота

	BotName        = "🚂 Чарнагорыя: расклад паездаў і электрычак | Черногория поезд"
	BotDescription = "" +
		`> Актуальны расклад
> Ведае ўсе вакзалы, уключаючы Белград
> Можа будаваць маршруты з перасадкай

Проста ўвядзіце два вакзалы праз коску:

Подгорица, Бар`
	BotShortDescription = "Актуальны расклад з усімі вакзаламі і маршрутамі, уключаючы маршруты з перасадкай і паезд Белград - Бар"

	// каманды бота

	BotCommandNameStart = "Старт бота"
	StartMessage        = "" +
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

	// /help

	BotCommandNameHelp = "Дапамога"
	HelpMessage        = "" +
		`Частыя адказы на пытанні, якія часта задаюць:

1. Гэта бот для раскладу цягнікоў у Чарнагорыі. У раскладзе таксама ёсць цягнікі, якія ідуць у / з Сербіі ў Чарнагорыю.
2. Ніякія іншыя краіны не маюць чыгуначных зносін з Чарнагорыяй, толькі Сербія.
3. Праверце карту станцый праз /map
4. Проста увядзіце дзве станцыі праз коску: 'Падгорыца, Бар', і вы атрымаеце расклад.
5. Білеты можна купіць толькі на станцыі або ў цягніку. Толькі наяўныя, анлайн-білетаў няма, часам карты прымаюцца на некаторых станцыях (так, часам).
6. Праверце цану, скідкі і іншыя дэталі па спасылцы 'Падрабязней' унізе раскладу.
7. Расклад застаецца нязменным на ўвесь год, за выключэннем аднаго летняга цягніка. Цягнік будзе курсіраваць з 13 чэрвеня па 14 верасня 2025 года па маршруце Субоціца - Бялград - Бар. Астатняе расклад застаецца такім жа.
8. Абнавіце расклад з дапамогай левай кнопкі "🔄 'дата'"
9. Часам цягнікі спазняюцца, асабліва ў летні сезон.
10. Больш падрабязная інфармацыя аб боце праз /about
`

	// /map

	BotCommandNameMap = "Карта ўсіх станцый"
	MapMessage        = "Карта з усімі станцыямі"

	// /about

	BotCommandNameAbout = "Пра гэтага бота"
	AboutMessage        = "" +
		`Гэты бот даступны па ліцэнзіі BEERWARE.

Пакуль вы бачыце гэта апавяшчэнне, вы можаце рабіць з гэтым кодам і гэтым ботам усё, што захочаце.
Калі мы калі-небудзь сустрэнемся, і вы лічыце,
што гэты бот карысны - вы можаце купіць мне піва ў якасці падзякі.

Я аддаю перавагу Никшичко тамака.

Я: https://github.com/ivanov-gv
Гэты праект: https://github.com/ivanov-gv/zpcg

Зроблена разам з @Leti_deshevle
`
)

var MonthsMap = map[time.Month]string{
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
}
