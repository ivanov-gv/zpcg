package render

import (
	"strings"
)

const (
	urlBase   = "https://zpcg.me/red-voznje?"
	fromParam = "start="
	toParam   = "finish="
)

func getUrlToTimetable(origin, destination string) string {
	var params []string
	if len(origin) > 0 {
		params = append(params, fromParam+trimSpaces(origin))
	}
	if len(destination) > 0 {
		params = append(params, toParam+trimSpaces(destination))
	}
	return urlBase + strings.Join(params, "&")
}

func trimSpaces(in string) string {
	return strings.ReplaceAll(in, " ", "+")
}
