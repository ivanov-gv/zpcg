package de

import "time"

// German

const (
	ErrorMessage = "" +
		`Versuchen Sie es erneut - zwei Bahnhöfe durch Komma getrennt. Hier ist ein Beispiel:

Podgorica, Niksic`

	StationDoesNotExistMessage   = "Dieser Bahnhof existiert nicht"
	RailwayMapButtonTextMap      = "Montenegro Bahnkarte"
	OfficialTimetableUrlText     = "Mehr erfahren"
	ReverseRouteInlineButtonText = "Zurück"
	AlertUpdateNotificationText  = "" +
		`Der Fahrplan wurde bereits aktualisiert

Vom 13.06.2025 bis 14.09.2025 wird ein Zug von Subotica nach Bar hinzugefügt

Ansonsten ändert sich der Fahrplan nicht`
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

>*Podgorica, Bijelo Polje*

und erhalten Sie den Fahrplan:

>Podgorica \> Bijelo Polje
>[6100](https://zpcg.me/details?timetable=41)  ` + "`06:20 \\> 08:38`" + `
>\.\.\.

Sind Sie sich unsicher, wie der Bahnhof richtig geschrieben wird? Schreiben Sie, wie Sie es wissen \- wir werden verstehen, was Sie gemeint haben\.

Jetzt sind Sie dran\!
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
7. Der Fahrplan bleibt das ganze Jahr über unverändert, mit Ausnahme eines Sommerzuges. Der Zug wird vom 13. Juni bis 14. September 2025 auf der Strecke Subotica – Belgrad – Bar verkehren. Der restliche Zeitplan bleibt genau gleich.
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

Ich: https://github.com/ivanov-gv
Dieses Projekt: https://github.com/ivanov-gv/zpcg
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
