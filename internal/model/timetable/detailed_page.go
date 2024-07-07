package timetable

import "time"

type DetailedTimetable struct {
	TrainId TrainId
	// TimetableUrl Deprecated: zpcg.me does not have links to exact train timetable anymore
	TimetableUrl  string
	International bool
	ValidFrom     time.Time
	ValidTo       time.Time
	Stops         []Stop
}
