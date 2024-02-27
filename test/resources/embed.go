package resources

import "embed"

const (
	DetailedTimetableFilepath    = "Željeznički prevoz Crne Gore - Broj voza 7108.html"
	GeneralTimetableHtmlFilepath = "Željeznički prevoz Crne Gore - Polasci.html"
)

//go:embed "Željeznički prevoz Crne Gore - Broj voza 7108.html"
//go:embed "Željeznički prevoz Crne Gore - Polasci.html"
var TestFS embed.FS
