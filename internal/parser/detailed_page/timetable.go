package detailed_page

import (
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"zpcg/internal/model"
	parser_model "zpcg/internal/parser/model"
	parserutils "zpcg/internal/parser/utils"
)

func ParseDetailedTimetablePage(routeNumber model.TrainId, detailedTimetableUrl string, reader io.Reader) (parser_model.DetailedTimetable, error) {
	tokenizer := html.NewTokenizer(reader)
	var timetable parser_model.DetailedTimetable
	for tokenType := tokenizer.Next(); tokenizer.Err() == nil; tokenType = tokenizer.Next() { // until the end of the page is not reached
		if tokenType != html.StartTagToken {
			continue
		}
		// tokenType == html.StartTagToken
		token := tokenizer.Token()
		if IsTimetableReached(token) {
			// check table type
			if tableType := GetTableType(tokenizer); tableType != DetailedTableRoute {
				continue
			}
			// found timetable with detailed route
			parsedTimetable, err := ParseRouteTable(tokenizer)
			if err != nil {
				return parser_model.DetailedTimetable{}, errors.Wrap(err, "ParseRouteTable")
			}
			timetable = parsedTimetable
		}
	}
	timetable.TrainId = routeNumber
	timetable.TimetableUrl = detailedTimetableUrl
	return timetable, nil
}

func IsTimetableReached(token html.Token) bool {
	//<div id="detail-stop-grid" class="grid-view">
	return token.Type == html.StartTagToken && token.Data == "div" &&
		parserutils.HasAttribute(token.Attr, "", "id", "detail-stop-grid")
}

type DetailedTableType int

const (
	NothingFound DetailedTableType = iota
	DetailedTablePrice
	DetailedTableRoute
)

func GetTableType(tokenizer *html.Tokenizer) DetailedTableType {
	// we are trying to find the following construction:
	//<table class="items">
	//<thead>
	//<tr>
	//<th id="detail-stop-grid_c0">Stanica</th><th id="detail-stop-grid_c1">Dolazak</th><th id="detail-stop-grid_c2">Polazak</th></tr>
	//</thead>
	var (
		// <th id="detail-stop-grid_c0">Stanica</th>
		foundTagThStanica bool
		// <th id="detail-stop-grid_c1">Dolazak</th>
		foundTagThDolazak bool
		// <th id="detail-stop-grid_c2">Polazak</th></tr>
		foundTagThPolazak bool
		// <th id="detail-stop-grid_c1">Drugi razred</th>
		foundTagThDrugiRazred bool
		// <th id="detail-stop-grid_c2">Prvi razred</th>
		foundTagThPrviRazred bool
	)

	for token := tokenizer.Token(); !IsTableHeadEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		var thTag, textTag html.Token
		if !(token.Type == html.StartTagToken && token.Data == "th") { // not a '<th ...>' tag
			continue
		}
		thTag = token
		tokenizer.Next()
		token = tokenizer.Token()
		if !(token.Type == html.TextToken) {
			continue
		}
		textTag = token

		// we found two subsequent tags - thTag and textTag. check them
		switch {
		case parserutils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c0") &&
			textTag.Data == "Stanica":
			foundTagThStanica = true
		case parserutils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c1") &&
			textTag.Data == "Dolazak":
			foundTagThDolazak = true
		case parserutils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c2") &&
			textTag.Data == "Polazak":
			foundTagThPolazak = true
		case parserutils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c1") &&
			textTag.Data == "Drugi razred":
			foundTagThDrugiRazred = true
		case parserutils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c2") &&
			textTag.Data == "Prvi razred":
			foundTagThPrviRazred = true
		}
	}
	switch {
	case foundTagThPolazak && foundTagThDolazak && foundTagThStanica:
		return DetailedTableRoute
	case foundTagThStanica && foundTagThDrugiRazred && foundTagThPrviRazred:
		return DetailedTablePrice
	default:
		return NothingFound
	}
}

func IsTableHeadEndReached(token html.Token) bool {
	//</thead>
	return token.Type == html.EndTagToken && token.Data == "thead"
}

func ParseRouteTable(tokenizer *html.Tokenizer) (parser_model.DetailedTimetable, error) {
	var result parser_model.DetailedTimetable
	for token := tokenizer.Token(); !parserutils.IsTableEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if parserutils.IsRowBeginningReached(token) {
			station, err := ParseRow(tokenizer)
			if err != nil {
				return parser_model.DetailedTimetable{}, errors.Wrap(err, "ParseRow")
			}
			result.Stations = append(result.Stations, station)
		}
	}
	return result, nil
}

func ParseRow(tokenizer *html.Tokenizer) (model.Stop, error) {
	var (
		cellNumber         = -1
		stationName        string
		arrival, departure time.Time
	)
	for token := tokenizer.Token(); !parserutils.IsRowEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if parserutils.IsCellBeginningReached(token) {
			cellNumber++
			continue
		}
		if cellNumber == -1 {
			continue
		}
		if parserutils.IsCellEndReached(token) {
			continue
		}

		// parse station name
		if cellNumber == 0 && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			stationName = token.Data
		}

		// parse arrival time
		if cellNumber == 1 && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			if token.Data == "\u00a0" {
				continue
			}
			var err error
			arrival, err = time.Parse("15:04", token.Data)
			if err != nil {
				return model.Stop{}, errors.Wrap(err, "can not parse arrival with time.Parse")
			}
		}

		// parse departure time
		if cellNumber == 2 && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			if token.Data == "\u00a0" {
				continue
			}
			var err error
			departure, err = time.Parse("15:04", token.Data)
			if err != nil {
				return model.Stop{}, errors.Wrap(err, "can not parse departure with time.Parse")
			}
		}
	}

	// for the first and the last station of the route departure or arrival is not present in timetable
	if departure.IsZero() {
		departure = arrival
	}
	if arrival.IsZero() {
		arrival = departure
	}
	return model.Stop{
		Id:        generateStationId(stationName),
		Arrival:   arrival,
		Departure: departure,
	}, nil
}
