package parser

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
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
