package parser

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"sync"
	"zpcg/internal/parser/detailed_page"
	"zpcg/internal/parser/general_page"
	"zpcg/internal/parser/model"
	"zpcg/internal/utils"
)

const (
	BaseUrl                 = "https://zpcg.me"
	GeneralTimetablePageUrl = "https://zpcg.me/search"
)

func ParseTimetable() (map[int]model.DetailedTimetable, error) {
	generalTimetableResponse, err := http.Get(GeneralTimetablePageUrl)
	if err != nil {
		return nil, errors.Wrap(err, "can not get general timetable page with http.Get")
	}
	generalTimetableMap, err := general_page.ParseGeneralTimetablePage(generalTimetableResponse.Body)
	if err != nil {
		return nil, errors.Wrap(err, "general_page.ParseGeneralTimetablePage")
	}

	detailedTimetableMap := sync.Map{}
	errGroup := errgroup.Group{}
	for routeId, generalTimetable := range generalTimetableMap {
		routeId, generalTimetable := routeId, generalTimetable
		errGroup.Go(func() error {
			response, err := retryablehttp.Get(BaseUrl + generalTimetable.DetailedTimetableLink)
			if err != nil {
				return errors.Wrapf(err, "can not get route info with route id = %d, link = %s using http.Get",
					routeId, generalTimetable.DetailedTimetableLink)
			}
			detailedTimetable, err := detailed_page.ParseDetailedTimetablePage(routeId, response.Body)
			if err != nil {
				return errors.Wrapf(err, "routeId = %d, link = %s detailed_page.ParseDetailedTimetablePage",
					routeId, generalTimetable.DetailedTimetableLink)
			}
			detailedTimetableMap.Store(detailedTimetable.RouteId, detailedTimetable)
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return nil, errors.Wrap(err, "errGroup.Wait")
	}
	result, err := utils.SyncMapToMap[int, model.DetailedTimetable](&detailedTimetableMap, len(generalTimetableMap))
	if err != nil {
		return nil, errors.Wrap(err, "utils.SyncMapToMap")
	}
	return result, nil
}
