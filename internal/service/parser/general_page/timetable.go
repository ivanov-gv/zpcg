package general_page

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
	parser_utils "github.com/ivanov-gv/zpcg/internal/service/parser/utils"
	"github.com/ivanov-gv/zpcg/internal/utils"
)

func ParseGeneralTimetablePage(reader io.Reader) (map[timetable.TrainId]timetable.GeneralTimetableRow, error) {
	tokenizer := html.NewTokenizer(reader)
	generalTimetableRows := map[timetable.TrainId]timetable.GeneralTimetableRow{}
	for tokenType := tokenizer.Next(); tokenizer.Err() == nil; tokenType = tokenizer.Next() { // until the end of the page is not reached
		if tokenType != html.StartTagToken {
			continue
		}
		token := tokenizer.Token()
		if !IsTableReached(token) {
			continue
		}
		// reached timetable
		table, err := ParseTable(tokenizer)
		if err != nil {
			return nil, fmt.Errorf("ParseTable: %w", err)
		}
		utils.AddMap(generalTimetableRows, table)
	}
	return generalTimetableRows, nil
}

func ParseTable(tokenizer *html.Tokenizer) (map[timetable.TrainId]timetable.GeneralTimetableRow, error) {
	result := map[timetable.TrainId]timetable.GeneralTimetableRow{}
	for token := tokenizer.Token(); !parser_utils.IsTableEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if !parser_utils.IsRowBeginningReached(token) {
			continue
		}
		// row reached
		row, err := ParseRow(tokenizer)
		if err != nil {
			return nil, fmt.Errorf("ParseRow: %w", err)
		}
		result[row.TrainId] = row
	}
	return result, nil
}

const (
	linkAttributeKey       = "href"
	cellNumberForRouteId   = 0
	cellNumberForTrainType = 6
)

func ParseRow(tokenizer *html.Tokenizer) (timetable.GeneralTimetableRow, error) {
	var (
		cellNumber = -1
		result     timetable.GeneralTimetableRow
	)
	for token := tokenizer.Token(); !parser_utils.IsRowEndReached(token); _, token = tokenizer.Next(), tokenizer.Token() {
		if parser_utils.IsCellBeginningReached(token) {
			cellNumber++ // cell beginning reached - +1 to cell count
			continue
		}
		if cellNumber == -1 { // skip all tags until first cell is reached
			continue
		}
		if parser_utils.IsCellEndReached(token) {
			continue
		}

		// parse route id
		if cellNumber == cellNumberForRouteId && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			trainId, err := strconv.Atoi(token.Data)
			if err != nil {
				return timetable.GeneralTimetableRow{}, fmt.Errorf("can not parse route id with strconv.Atoi: %w", err)
			}
			result.TrainId = timetable.TrainId(trainId)
			continue
		}

		// parse train type
		if cellNumber == cellNumberForTrainType && token.Type == html.TextToken && !strings.Contains(token.Data, "\n") {
			trainType, err := ParseTrainType(token.Data)
			if err != nil {
				return timetable.GeneralTimetableRow{}, fmt.Errorf("ParseTrainType: %w", err)
			}
			result.TrainType = trainType
			continue
		}

		// parse detailed timetable link
		if found, url := parser_utils.FindAttribute(token.Attr, "", linkAttributeKey); IsLinkToDetailedTimetableFound(token) && found {
			result.DetailedTimetableLink = url
			continue
		}
	}
	return result, nil
}

const (
	TokenForLocalTrainType = "lokalni"
	TokenForFastTrainType  = "brzi"
)

func ParseTrainType(tokenData string) (timetable.TrainType, error) {
	switch tokenData {
	case TokenForLocalTrainType:
		return timetable.LocalTrain, nil
	case TokenForFastTrainType:
		return timetable.FastTrain, nil
	default:
		return 0, fmt.Errorf("can't parse train type from token: %s", tokenData)
	}
}
