package callback

import (
	"fmt"
	"strings"

	"github.com/ivanov-gv/zpcg/internal/model/callback"
	"github.com/ivanov-gv/zpcg/internal/service/date"
)

func NewCallbackService() *CallbackService {
	return &CallbackService{}
}

type CallbackService struct{}

const (
	delimiter = "|"
)

const (
	typeIndex = iota
	dataStartIndex
)

func (s *CallbackService) ParseCallback(callbackData string) (callback.Callback, error) {
	splitData := strings.Split(callbackData, delimiter)
	switch callback.Type(splitData[typeIndex]) {
	case callback.UpdateType:
		return parseUpdateCallbackData(splitData)
	case callback.ReverseRouteType:
		return parseReverseRouteCallbackData(splitData)
	default:
		return callback.Callback{}, ErrInvalidCallbackType
	}
}

// update callback

func (s *CallbackService) GenerateUpdateCallbackData(origin, destination, currentDate string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s",
		callback.UpdateType, delimiter,
		origin, delimiter,
		destination, delimiter,
		currentDate)
}

const (
	updateCallback_OriginIndex = dataStartIndex + iota
	updateCallback_DestinationIndex
	updateCallback_DateIndex

	updateCallback_RequiredLen
)

func parseUpdateCallbackData(data []string) (callback.Callback, error) {
	if len(data) < updateCallback_RequiredLen {
		return callback.Callback{}, fmt.Errorf("len(data) is less than required [data='%v', required len='%d']: %w",
			data, updateCallback_RequiredLen, ErrInvalidData)
	}
	currentDate, err := date.CurrentDateFromShortString(data[updateCallback_DateIndex])
	if err != nil {
		return callback.Callback{}, fmt.Errorf("date.CurrentDateFromShortString: %w", err)
	}
	return callback.Callback{
		Type: callback.UpdateType,
		UpdateData: callback.UpdateData{
			Origin:      data[updateCallback_OriginIndex],
			Destination: data[updateCallback_DestinationIndex],
			Date:        currentDate,
		},
	}, nil
}

// reverse route callback

func (s *CallbackService) GenerateReverseRouteCallbackData(origin, destination string) string {
	return fmt.Sprintf("%s%s%s%s%s",
		callback.ReverseRouteType, delimiter,
		origin, delimiter,
		destination)
}

const (
	reverseRouteCallback_OriginIndex = dataStartIndex + iota
	reverseRouteCallback_DestinationIndex

	reverseRouteCallback_RequiredLen
)

func parseReverseRouteCallbackData(data []string) (callback.Callback, error) {
	if len(data) < reverseRouteCallback_RequiredLen {
		return callback.Callback{}, fmt.Errorf("len(data) is less than required [data='%v', required len ='%d']: %w",
			data, reverseRouteCallback_RequiredLen, ErrInvalidData)
	}
	return callback.Callback{
		Type: callback.ReverseRouteType,
		ReverseRouteData: callback.ReverseRouteData{
			Origin:      data[reverseRouteCallback_OriginIndex],
			Destination: data[reverseRouteCallback_DestinationIndex],
		},
	}, nil
}
