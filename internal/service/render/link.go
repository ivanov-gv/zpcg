package render

import (
	"strings"
	"time"
)

const (
	urlBase    = "https://ivanov-gv.github.io/zpcg/official-site-redirect.html?"
	fromParam  = "from="
	toParam    = "to="
	dateParam  = "date="
	dateLayout = "02.01.2006."
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
