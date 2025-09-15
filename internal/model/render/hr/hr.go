package hr

import "time"

const (
	ErrorMessage = "" +
		`Pokušajte ponovno - dva kolodvora odvojena zarezom. Evo primjera:

Podgorica, Nikšić`

	StationDoesNotExistMessage   = "Taj kolodvor ne postoji"
	RailwayMapButtonTextMap      = "Karta željeznica Crne Gore"
	OfficialTimetableUrlText     = "Više informacija"
	ReverseRouteInlineButtonText = "Natrag"
	AlertUpdateNotificationText  = "" +
		`Raspored je već ažuriran

Do 14. prosinca 2025 raspored se neće mijenjati`
	SimpleUpdateNotificationText = "Raspored za danas je ažuriran"

	// Opis bota

	BotName        = "🚂 Crna Gora: raspored vlakova | ZPCG RED VOŽNJE"
	BotDescription = "" +
		`> Trenutni raspored
> Zna sve kolodvore, uključujući Beograd
> Može graditi rute s presjedanjem

Jednostavno unesite dva kolodvora odvojena zarezom:

Podgorica, Bar`
	BotShortDescription = "Trenutni raspored sa svim kolodvorima i rutama, uključujući rute s presjedanjem i vlak Beograd - Bar"

	// Naredbe bota

	BotCommandNameStart = "Pokreni bota"
	StartMessage        = "" +
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

	// /help

	BotCommandNameHelp = "Pomoć"
	HelpMessage        = "" +
		`Često postavljana pitanja s čestim odgovorima:

1. Ovo je bot za vozni red vlakova u Crnoj Gori. Vozni red također uključuje vlakove koji idu prema/iz Srbije u Crnu Goru.
2. Nijedna druga zemlja nema željezničke veze s Crnom Gorom, samo Srbija.
3. Provjerite kartu stanice putem /map
4. Samo unesite dvije stanice odvojene zarezima: 'Podgorica, Bar' i dobit ćete raspored.
5. Karte se mogu kupiti samo na kolodvoru ili u vlaku. Samo gotovina, nema online karata, ponekad se na nekim stanicama primaju kartice (da, ponekad).
6. Provjerite cijenu, popuste i ostale detalje klikom na poveznicu 'Više detalja' pri dnu rasporeda.
7. Vozni red ostaje isti tijekom cijele godine, s izuzetkom jednog ljetnog vlaka. Ostatak rasporeda ostaje potpuno isti.
8. Ažurirajte raspored pomoću lijeve tipke "🔄 'datum'"
9. Ponekad vlakovi kasne, posebno tijekom ljetne sezone.
10. Detaljnije informacije o botu putem /about
`

	// /map

	BotCommandNameMap = "Karta svih stanica"
	MapMessage        = "Karta sa svim stanicama"

	// /about

	BotCommandNameAbout = "O ovom botu"
	AboutMessage        = "" +
		`Ovaj bot je dostupan pod BEERWARE licencom.

Sve dok vidite ovu obavijest, možete raditi što god želite s ovim kodom i ovim botom.
Ako se ikada sretnemo i pomisliš,
da je ovaj bot koristan - možete mi kupiti pivo kao zahvalu.

Više volim Nikšićko tamno.

Ja: https://github.com/ivanov-gv
Ovaj projekt: https://github.com/ivanov-gv/zpcg

Napravljeno zajedno s @Leti_deshevle
`
)

var MonthsMap = map[time.Month]string{
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
}
