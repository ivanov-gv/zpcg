package en

import "time"

const (
	ErrorMessage = "" +
		`Try again - two stations, separated by a comma. Just like that:

Podgorica, Bar`
	StationDoesNotExistMessage   = "This station does not exist"
	RailwayMapButtonTextMap      = "Montenegro Railway Map"
	OfficialTimetableUrlText     = "More details"
	ReverseRouteInlineButtonText = "Reverse"
	SimpleUpdateNotificationText = "Today's timetable is updated"

	// bot description

	BotName        = "🚂 Montenegro: train timetable | Черногория расписание поезд"
	BotDescription = "" +
		`> Up-to-date timetable
> Knows every station, including Belgrade
> Can show routes between any two station, including transfer

Just type two stations with a comma:

Podgorica, Bar`
	BotShortDescription = "Up-to-date timetable with all stations and routes, including transfer and international ones, like Belgrade - Bar train"

	// bot commands

	// /start

	BotCommandNameStart = "Start the bot"
	StartMessage        = "" +
		`*Montenegro Railways Timetable*

_Made together with @Leti\_deshevle_

Please enter *two stations* separated by *a comma*: 

>*Podgorica, Bar*

Or using cyrillic:

>*Подгорица, Бар*

`

	// /help

	BotCommandNameHelp = "Help"
	HelpMessage        = "" +
		`Frequently given answers (for the frequently asked questions):

1. This is a bot for trains in Montenegro. You can find a train to or from Serbia to Montenegro also. 
2. No other countries have a train connection with Montenegro, only Serbia.
3. Check stations /map
4. Just type two stations with a comma: 'Podgorica, Bar', and you'll get the timetable.
5. You can buy tickets only on a station or in a train. Only cash, no online tickets, cards are acceptable on some station, sometimes (yes, sometimes).
6. Check the price, discounts and other details with the 'More details' link on the bottom of the timetable
7. Timetable remains the same for the whole year, except for the one summer train. The rest of the timetable remains exactly the same.
8. Update timetable with the left "🔄 'date'" button 	
9. Sometimes trains are running late, especially in summer season.
10. More information about the bot with /about
`

	// /map

	BotCommandNameMap = "A map with all stations"
	MapMessage        = "A map with all stations"

	// /about

	BotCommandNameAbout = "About this bot"
	AboutMessage        = "" +
		`This bot is accessible under BEERWARE license.

As long as you retain this notice you
can do whatever you want with this stuff. If we meet some day, and you think
this stuff is worth it, you can buy me a beer in return. 

I prefer Nikšićko tamno.

https://github.com/ivanov-gv/zpcg

Made together with @Leti_deshevle
`
)

var MonthsMap = map[time.Month]string{
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
}
