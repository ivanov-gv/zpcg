package callback

import (
	"fmt"
	"strings"

	"github.com/ivanov-gv/zpcg/internal/model/callback"
)

func NewCallbackService() *CallbackService {
	return &CallbackService{}
}

type CallbackService struct{}

const delimiter = "|"

func (s *CallbackService) ParseCallback(callbackData string) callback.Callback {
	splitData := strings.Split(callbackData, delimiter)
	switch callback.Type(splitData[0]) {
	case callback.UpdateType:
		if len(splitData) < 3 {
			return callback.Callback{}
		}
		return callback.Callback{
			Type: callback.UpdateType,
			Data: callback.Data{
				Origin:      splitData[1],
				Destination: splitData[2],
			},
		}
	case callback.ReverseRouteType:
		if len(splitData) < 3 {
			return callback.Callback{}
		}
		return callback.Callback{
			Type: callback.ReverseRouteType,
			Data: callback.Data{
				Origin:      splitData[1],
				Destination: splitData[2],
			},
		}
	default:
		return callback.Callback{}
	}
}

func (s *CallbackService) GenerateUpdateCallbackData(origin, destination string) string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		callback.UpdateType, delimiter,
		origin, delimiter,
		destination, delimiter) // TODO: add current date
}

func (s *CallbackService) GenerateReverseRouteCallbackData(origin, destination string) string {
	return fmt.Sprintf("%s%s%s%s%s%s",
		callback.ReverseRouteType, delimiter,
		origin, delimiter,
		destination, delimiter)
}
