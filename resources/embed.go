package resources

import (
	"embed"
)

const TimetableGobFileName = "timetable.gob"

//go:embed timetable.gob
var FS embed.FS
