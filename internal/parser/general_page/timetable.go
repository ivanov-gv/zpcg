package general_page

import (
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"io"
	"strconv"
	"strings"
	"zpcg/internal/parser/model"
	parserutils "zpcg/internal/parser/utils"
	"zpcg/internal/utils"
)

func ParseGeneralTimetablePage(reader io.Reader) (map[int]model.GeneralTimetableRow, error) {
	tokenizer := html.NewTokenizer(reader)
	generalTimetableRows := map[int]model.GeneralTimetableRow{}
	for tokenType := tokenizer.Next(); tokenizer.Err() == nil; tokenType = tokenizer.Next() { // until the end of the page is not reached
		if tokenType != html.StartTagToken {
			continue
		}
		// tokenType == html.StartTagToken
		token := tokenizer.Token()
		if IsTableReached(token) {
			// found timetable
			table, err := ParseTable(tokenizer)
			if err != nil {
				return nil, errors.Wrap(err, "ParseTable")
			}
			utils.AddMap(generalTimetableRows, table)
		}
	}
	return generalTimetableRows, nil
}

func ParseTable(tokenizer *html.Tokenizer) (map[int]model.GeneralTimetableRow, error) {
	result := map[int]model.GeneralTimetableRow{}
	for token := tokenizer.Token(); !parserutils.IsTableEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if parserutils.IsRowBeginningReached(token) {
			row, err := ParseRow(tokenizer)
			if err != nil {
				return nil, errors.Wrap(err, "ParseRow")
			}
			result[row.RouteId] = row
		}
	}
	return result, nil
}

func ParseRow(tokenizer *html.Tokenizer) (model.GeneralTimetableRow, error) {
	var (
		cellNumber = -1
		result     model.GeneralTimetableRow
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

		// parse route id
		if cellNumber == 0 && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			routeId, err := strconv.Atoi(token.Data)
			if err != nil {
				return model.GeneralTimetableRow{}, errors.Wrap(err, "can not parse route id with strconv.Atoi")
			}
			result.RouteId = routeId
		}

		// parse detailed timetable link
		if found, url := parserutils.FindAttribute(token.Attr, "", "href"); IsLinkToDetailedTimetabelFound(token) && found {
			result.DetailedTimetableLink = url
		}
	}
	return result, nil
}
