package de

import "time"

// German

const (
	ErrorMessage = "" +
		`Versuchen Sie es erneut - zwei Bahnhöfe durch Komma getrennt. Hier ist ein Beispiel:

Podgorica, Bar`

	StationDoesNotExistMessage   = "Dieser Bahnhof existiert nicht"
	RailwayMapButtonTextMap      = "Montenegro Bahnkarte"
	OfficialTimetableUrlText     = "Mehr erfahren"
	ReverseRouteInlineButtonText = "Zurück"
	SimpleUpdateNotificationText = "Der Fahrplan für heute wurde aktualisiert"

	// Bot-Beschreibung

	BotName        = "🚂 Montenegro: Zug- und Zugfahrplan | train timetable"
	BotDescription = "" +
		`> Aktueller Fahrplan
> Kennt alle Bahnhöfe, einschließlich Belgrad
> Kann Routen mit Umstieg erstellen

Geben Sie einfach zwei Bahnhöfe durch Komma getrennt ein:

Podgorica, Bar`
	BotShortDescription = "Aktueller Fahrplan mit allen Bahnhöfen und Routen, einschließlich Routen mit Umstieg und dem Zug Belgrad - Bar"

	// Bot-Befehle

	BotCommandNameStart = "Start des Bots"
	StartMessage        = "" +
		`*Fahrplan der Züge in Montenegro*

_Erstellt in Zusammenarbeit mit @Leti\_deshevle_

Bitte geben Sie *zwei Bahnhöfe durch Komma getrennt* ein, auf Lateinisch: 

>*Podgorica, Bar*

`

	// /help

	BotCommandNameHelp = "Hilfe"
	HelpMessage        = "" +
		`Häufig gestellte Fragen, häufig beantwortet:

1. Dies ist ein Bot für Zugfahrpläne in Montenegro. Der Fahrplan umfasst auch Züge von/nach Serbien nach Montenegro.
2. Nur Serbien hat keine anderen Länder mit Montenegro per Bahn verbunden.
3. Überprüfen Sie die Stationskarte über /map
4. Geben Sie einfach zwei durch Kommas getrennte Stationen ein: „Podgorica, Bar“ und Sie erhalten den Fahrplan.
5. Fahrkarten können nur am Bahnhof oder im Zug erworben werden. Nur Bargeld, keine Online-Tickets, manchmal werden an manchen Stationen Karten akzeptiert (ja, manchmal).
6. Überprüfen Sie Preis, Rabatte und andere Details, indem Sie unten im Zeitplan auf den Link „Weitere Details“ klicken.
7. Der Fahrplan bleibt das ganze Jahr über unverändert, mit Ausnahme eines Sommerzuges. Der restliche Zeitplan bleibt genau gleich.
8. Aktualisieren Sie den Zeitplan mit der linken Schaltfläche "🔄 'Datum'"
9. Manchmal haben Züge Verspätung, besonders während der Sommersaison.
10. Nähere Informationen zum Bot über /about
`

	// /map

	BotCommandNameMap = "Karte aller Stationen"
	MapMessage        = "Karte mit allen Stationen"

	// /about

	BotCommandNameAbout = "Über diesen Bot"
	AboutMessage        = "" +
		`Dieser Bot ist unter der BEERWARE-Lizenz verfügbar.

Solange Sie diese Benachrichtigung sehen, können Sie mit diesem Code und diesem Bot tun, was Sie wollen.
Wenn wir uns jemals treffen und du denkst,
dass dieser Bot nützlich ist – als Dankeschön kannst du mir ein Bier ausgeben.

Ich bevorzuge Niksicko Tamno.

https://github.com/ivanov-gv/zpcg

Gemeinsam erstellt mit @Leti_deshevle
`
)

var MonthsMap = map[time.Month]string{
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
}
