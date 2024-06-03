package render

import (
	"strings"
	"time"
)

const (
	urlBase    = "https://zpcg.me/red-voznje?"
	fromParam  = "start="
	toParam    = "finish="
	dateParam  = "date="
	dateLayout = "2006-01-02"
)

func getUrlToTimetable(origin, destination string, date time.Time) string {
	var params []string
	if len(origin) > 0 {
		params = append(params, fromParam+trimSpaces(origin))
	}
	if len(destination) > 0 {
		params = append(params, toParam+trimSpaces(destination))
	}
	if !date.IsZero() {
		params = append(params, dateParam+date.Format(dateLayout))
	}
	return urlBase + strings.Join(params, "&")
}

func trimSpaces(in string) string {
	return strings.Replace(in, " ", "+", -1)
}
