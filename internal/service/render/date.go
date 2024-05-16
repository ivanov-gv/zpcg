package render

import (
	"fmt"
	"time"

	"golang.org/x/text/language"
)

func Date(tag language.Tag, currentDate time.Time) string {
	_, month, day := currentDate.Date()
	return fmt.Sprintf("%d %s", day, localizeMonth(tag, month))
}

func localizeMonth(tag language.Tag, month time.Month) string {
	return GetMessage(MonthNameMap, tag)[month]
}
