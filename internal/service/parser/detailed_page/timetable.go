package detailed_page

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/net/html"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	parser_utils "github.com/ivanov-gv/zpcg/internal/service/parser/utils"
	"github.com/ivanov-gv/zpcg/internal/utils"
)

func ParseDetailedTimetablePage(routeNumber timetable.TrainId, detailedTimetableUrl string, reader io.Reader) (timetable.DetailedTimetable, error) {
	tokenizer := html.NewTokenizer(reader)
	var _timetable timetable.DetailedTimetable
	for tokenType := tokenizer.Next(); tokenizer.Err() == nil; tokenType = tokenizer.Next() { // until the end of the page is not reached
		if tokenType != html.StartTagToken {
			continue
		}
		// tokenType == html.StartTagToken
		token := tokenizer.Token()
		if !IsTimetableReached(token) {
			continue
		}
		// check table type
		if tableType := GetTableType(tokenizer); tableType != DetailedTableRoute {
			continue
		}
		// found timetable with detailed route
		parsedTimetable, err := ParseRouteTable(tokenizer)
		if err != nil {
			return timetable.DetailedTimetable{}, fmt.Errorf("ParseRouteTable: %w", err)
		}
		_timetable = parsedTimetable
	}
	_timetable.TrainId = routeNumber
	_timetable.TimetableUrl = detailedTimetableUrl
	return _timetable, nil
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
		case parser_utils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c0") &&
			textTag.Data == "Stanica":
			foundTagThStanica = true
		case parser_utils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c1") &&
			textTag.Data == "Dolazak":
			foundTagThDolazak = true
		case parser_utils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c2") &&
			textTag.Data == "Polazak":
			foundTagThPolazak = true
		case parser_utils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c1") &&
			textTag.Data == "Drugi razred":
			foundTagThDrugiRazred = true
		case parser_utils.HasAttribute(thTag.Attr, "", "id", "detail-stop-grid_c2") &&
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

func ParseRouteTable(tokenizer *html.Tokenizer) (timetable.DetailedTimetable, error) {
	var result timetable.DetailedTimetable
	for token := tokenizer.Token(); !parser_utils.IsTableEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if !parser_utils.IsRowBeginningReached(token) {
			continue
		}
		// row beginning reached
		// if the time is not present in the timetable - just get the last station departure time or empty one
		var fallbackTime time.Time
		if prevStop, found := lo.Last(result.Stops); found {
			fallbackTime = prevStop.Departure
		}
		station, err := ParseRow(tokenizer, fallbackTime)
		if err != nil {
			return timetable.DetailedTimetable{}, fmt.Errorf("ParseRow: %w", err)
		}
		// there might be a train with stops after midnight
		// in this case we need to add 24h to the arrival/departure time
		if prevStop, found := lo.Last(result.Stops); found &&
			(prevStop.Departure.After(utils.Midnight) || // previous stop departure is after midnight
				prevStop.Departure.After(station.Arrival) || // previous stop departure is before midnight, but current stop arrival is after midnight. like 23:58 -> 00:02
				station.Arrival.After(station.Departure)) { // arrival at a current stop is after departure: 23:58 -> 00:02
			if !station.Arrival.After(station.Departure) { // if both arrival and departure are after midnight - add 24h to both. example: 00:02 -> 00:05
				station.Arrival = station.Arrival.Add(utils.Day)
			}
			station.Departure = station.Departure.Add(utils.Day)
		}
		result.Stops = append(result.Stops, station)
	}
	return result, nil
}

func ParseRow(tokenizer *html.Tokenizer, fallbackTime time.Time) (timetable.Stop, error) {
	var (
		cellNumber         = -1
		stationName        string
		arrival, departure time.Time
	)
	for token := tokenizer.Token(); !parser_utils.IsRowEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if parser_utils.IsCellBeginningReached(token) {
			cellNumber++
			continue
		}
		if cellNumber == -1 {
			continue
		}
		if parser_utils.IsCellEndReached(token) {
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
				return timetable.Stop{}, fmt.Errorf("can not parse arrival with time.Parse: %w", err)
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
				return timetable.Stop{}, fmt.Errorf("can not parse departure with time.Parse: %w", err)
			}
		}
	}

	// timetable has empty lines sometimes
	if departure.IsZero() && arrival.IsZero() {
		departure = fallbackTime
		arrival = fallbackTime
	}

	// for the first and the last station of the route departure or arrival is not present in the timetable
	if departure.IsZero() {
		departure = arrival
	}
	if arrival.IsZero() {
		arrival = departure
	}
	return timetable.Stop{
		Id:        generateStationId(stationName),
		Arrival:   arrival,
		Departure: departure,
	}, nil
}
