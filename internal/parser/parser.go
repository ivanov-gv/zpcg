package parser

import (
	"encoding/gob"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"os"
	"zpcg/internal/parser/detailed_page"
	"zpcg/internal/parser/general_page"
	"zpcg/internal/parser/model"
)

const (
	BaseUrl                 = "https://zpcg.me"
	GeneralTimetablePageUrl = "https://zpcg.me/search"
)

func ParseTimetable() (map[int]model.DetailedTimetable, error) {
	generalTimetableResponse, err := retryablehttp.Get(GeneralTimetablePageUrl)
	if err != nil {
		return nil, errors.Wrap(err, "can not get general timetable page with retryablehttp.Get")
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "general_page.ParseGeneralTimetablePage")
	}

	detailedTimetableMap := make(map[int]model.DetailedTimetable, len(generalTimetableMap))
	// do not rewrite this loop with concurrency because zpcg.me do not have enough resources to handle all those requests
	// concurrency version is in the commit f5a2f983ce73fcc74f271d3bc4db51c2c56fe89f
	for routeId, generalTimetable := range generalTimetableMap {
		routeId, generalTimetable := routeId, generalTimetable
		response, err := retryablehttp.Get(BaseUrl + generalTimetable.DetailedTimetableLink)
		if err != nil {
			return nil, errors.Wrapf(err, "can not get route info with route id = %d, link = %s using retryablehttp.Get",
				routeId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(routeId, response.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "routeId = %d, link = %s detailed_page.ParseDetailedTimetablePage",
				routeId, generalTimetable.DetailedTimetableLink)
		}
		detailedTimetableMap[detailedTimetable.RouteId] = detailedTimetable
	}
	return detailedTimetableMap, nil
}

func ExportTimetable(timetable map[int]model.DetailedTimetable, filename string) error {
	file, err := os.OpenFile(filename+".gob", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "can not open file with os.OpenFile")
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(timetable)
	if err != nil {
		return errors.Wrap(err, "can not encode timetable with enc.Encode")
	}
	return nil
}

func ImportTimetable(filename string) (map[int]model.DetailedTimetable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "can not open file with os.Open")
	}
	enc := gob.NewDecoder(file)
	result := &map[int]model.DetailedTimetable{}
	err = enc.Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "can not decode timetable with enc.Decode")
	}
	return *result, nil
}
