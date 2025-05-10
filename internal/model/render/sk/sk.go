package sk

import "time"

const (
	ErrorMessage = "" +
		`Skúste to znova - dva stanice oddelené čiarkou. Tu je príklad:

Podgorica, Nikšić`

	StationDoesNotExistMessage   = "Taká stanica neexistuje"
	RailwayMapButtonTextMap      = "Mapa železníc Čiernej Hory"
	OfficialTimetableUrlText     = "Viac informácií"
	ReverseRouteInlineButtonText = "Späť"
	AlertUpdateNotificationText  = "" +
		`Cestovný poriadok už bol aktualizovaný

Od 13.06.2025 do 14.09.2025 bude pridaný vlak Subotica - Bar

Inak sa poradie nezmení`
	SimpleUpdateNotificationText = "Cestovný poriadok pre dnešok bol aktualizovaný"

	// Popis bota

	BotName        = "🚂 Čierna Hora: cestovný poriadok vlakov | Montenegro train"
	BotDescription = "" +
		`> Aktuálny cestovný poriadok
> Pozná všetky stanice, vrátane Belehradu
> Vie zostaviť trasy s prestupom

Jednoducho zadajte dve stanice oddelené čiarkou:

Podgorica, Bar`
	BotShortDescription = "Aktuálny cestovný poriadok so všetkými stanicami a trasami, vrátane trás s prestupom a vlaku Belehrad - Bar"

	// Príkazy bota

	BotCommandNameStart = "Spustiť bota"
	StartMessage        = "" +
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

	// /help

	BotCommandNameHelp = "Pomoč"
	HelpMessage        = "" +
		`Často kladené otázky, na ktoré sa často odpovedá:

1. Toto je bot pre cestovné poriadky vlakov v Čiernej Hore. Cestovný poriadok zahŕňa aj vlaky smerujúce do/zo Srbska do Čiernej Hory.
2. Žiadne iné krajiny nemajú železničné spojenie s Čiernou Horou, iba Srbsko.
3. Skontrolujte mapu stanice cez /map
4. Stačí zadať dve stanice oddelené čiarkami: „Podgorica, Bar“ a zobrazí sa vám cestovný poriadok.
5. Lístky je možné zakúpiť iba na stanici alebo vo vlaku. Iba hotovosť, žiadne online lístky, niekedy sa na niektorých staniciach akceptujú karty (áno, niekedy).
6. Cenu, zľavy a ďalšie podrobnosti si overte kliknutím na odkaz „Viac informácií“ v dolnej časti rozvrhu.
7. Cestovný poriadok zostáva počas celého roka rovnaký, s výnimkou jedného letného vlaku. Vlak bude premávať od 13. júna do 14. septembra 2025 na trase Subotica - Belehrad - Bar. Zvyšok harmonogramu zostáva úplne rovnaký.
8. Aktualizujte rozvrh pomocou ľavého tlačidla „🔄 'dátum'“
9. Vlaky niekedy meškajú, najmä počas letnej sezóny.
10. Podrobnejšie informácie o bote nájdete na /about
`

	// /map

	BotCommandNameMap = "Mapa všetkých staníc"
	MapMessage        = "Mapa so všetkými stanicami"

	// /about

	BotCommandNameAbout = "O tem botu"
	AboutMessage        = "" +
		`Tento bot je dostupný pod licenciou BEERWARE.

Pokiaľ vidíte toto upozornenie, môžete s týmto kódom a týmto botom robiť čokoľvek chcete.
Ak sa niekedy stretneme a pomyslíš si,
že tento bot je užitočný - ako poďakovanie mi môžete kúpiť pivo.

Ja mám radšej Niksicko tamno.

Ja: https://github.com/ivanov-gv
Tento projekt: https://github.com/ivanov-gv/zpcg
`
)

var MonthsMap = map[time.Month]string{
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
}
