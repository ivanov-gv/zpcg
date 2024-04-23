package app

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/language"

	callback_model "github.com/ivanov-gv/zpcg/internal/model/callback"
	"github.com/ivanov-gv/zpcg/internal/model/message"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"

	"github.com/ivanov-gv/zpcg/internal/service/blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/callback"
	"github.com/ivanov-gv/zpcg/internal/service/name"
	"github.com/ivanov-gv/zpcg/internal/service/pathfinder"
	"github.com/ivanov-gv/zpcg/internal/service/render"
	"github.com/ivanov-gv/zpcg/internal/service/transfer"
)

func NewApp() (*App, error) {
	// timetable
	_timetable := transfer.ImportTimetable()
	// finder
	finder := pathfinder.NewPathFinder(_timetable.StationIdToTrainIdSet, _timetable.TrainIdToStationMap, _timetable.TransferStationId)
	// name resolver
	stationNameResolver := name.NewStationNameResolver(_timetable.UnifiedStationNameToStationIdMap, _timetable.UnifiedStationNameList)
	// render
	_render := render.NewRender(_timetable.StationIdToStationMap, _timetable.TrainIdToTrainInfoMap)
	// blacklist
	blackList := blacklist.NewBlackListService()

	// complete app
	return &App{
		finder:              finder,
		stationNameResolver: stationNameResolver,
		render:              _render,
		blackList:           blackList,
		transferStationId:   _timetable.TransferStationId,
	}, nil
}

// App handles requests in the most general form. It has to know nothing about specific messenger, api, libs etc.
// It simply takes message and generates response. Contains whole business logic and business logic only
type App struct {
	finder              *pathfinder.PathFinder
	stationNameResolver *name.StationNameResolver
	render              *render.Render
	blackList           *blacklist.BlackListService
	callback            *callback.CallbackService
	transferStationId   timetable.StationId
}

func (a *App) HandleUpdate(update message.Update) (responseWithChatIds message.ResponseWithChatId, warning error) {
	switch update.Type {
	case message.MessageUpdateType:
		return a.HandleMessage(update.Message)
	case message.CallbackUpdateType:
		return a.HandleCallback(update.Callback)
	default: // including model.UnsupportedUpdateType
		return message.ResponseWithChatId{}, nil
	}
}

func (a *App) HandleCallback(callbackMessage message.Callback) (responseWithChatIds message.ResponseWithChatId, warning error) {
	languageTag := render.ParseLanguageTag(callbackMessage.From.LanguageCode)
	_callback := a.callback.ParseCallback(callbackMessage.Data)
	switch _callback.Type {
	case callback_model.UpdateType:
		response, err := a.GenerateRouteForStations(languageTag, _callback.Data.Origin, _callback.Data.Destination)
		if err != nil {
			return message.ResponseWithChatId{},
				fmt.Errorf("GenerateRouteForStations [lang=%s, origin=%s, destination=%s] : %w",
					languageTag, _callback.Data.Origin, _callback.Data.Destination, err)
		}
		update := []message.ToUpdate{ // TODO: do not update, if there is nothing to change
			{
				MessageId: callbackMessage.Message.Id,
				Response:  response,
			},
		}
		return message.ResponseWithChatId{ChatId: callbackMessage.ChatId, Update: update}, nil
	case callback_model.ReverseRouteType:
		response, err := a.GenerateRouteForStations(languageTag, _callback.Data.Destination, _callback.Data.Origin)
		if err != nil {
			return message.ResponseWithChatId{},
				fmt.Errorf("GenerateRouteForStations [lang=%s, origin=%s, destination=%s] : %w",
					languageTag, _callback.Data.Destination, _callback.Data.Origin, err)
		}
		update := []message.ToUpdate{
			{
				MessageId: callbackMessage.Message.Id,
				Response:  response, // TODO: updates only once for some reason. seems like callback data generator is buggy
			},
		}
		return message.ResponseWithChatId{ChatId: callbackMessage.ChatId, Update: update}, nil
	default:
		return message.ResponseWithChatId{},
			fmt.Errorf("unknown callback type(%s) [data='%s']", _callback.Type, _callback.Data)
	}
}

func (a *App) HandleMessage(_message message.Message) (responseWithChatIds message.ResponseWithChatId, warning error) {
	languageTag := render.ParseLanguageTag(_message.From.LanguageCode)
	var (
		response message.Response
		err      error
	)
	switch {
	case strings.HasPrefix(_message.Text, "/"):
		// got command - send start message
		response = a.render.StartMessage(languageTag)
	case strings.ContainsAny(_message.Text, string(lo.SpecialCharset)):
		// got message with stations - send a timetable
		response, err = a.GenerateRoute(languageTag, _message.Text)
		if err != nil {
			err = fmt.Errorf("a.GenerateRoute: %w", err)
		}
	default:
		// error otherwise
		err = fmt.Errorf("unknown message type: %s", _message.Text)
	}
	// handle error
	if err != nil {
		response = a.render.ErrorMessage(languageTag)
	}
	send := []message.Response{response}
	return message.ResponseWithChatId{
		Send:   send,
		ChatId: _message.ChatId,
	}, err
}

const stationsDelimiter = ','

func parseInputStations(input string) (originStation, destinationStation string, err error) {
	// convert all special chars to stationsDelimiterComma
	inputWithProperDelimiter := strings.Map(func(r rune) rune {
		if strings.ContainsAny(string(r), string(lo.SpecialCharset)) {
			return stationsDelimiter
		}
		return r
	}, input)

	// parse stations
	stations := strings.Split(inputWithProperDelimiter, string(stationsDelimiter))
	if len(stations) < 2 {
		return "", "", fmt.Errorf("not enough stations provided: %s", inputWithProperDelimiter)
	}
	return stations[0], lo.Must(lo.Last(stations)), nil
}

func (a *App) GenerateRoute(languageTag language.Tag, input string) (message.Response, error) {
	origin, destination, err := parseInputStations(input)
	if err != nil {
		return message.Response{}, fmt.Errorf("parseInputStations: %w", err)
	}
	return a.GenerateRouteForStations(languageTag, origin, destination)
}

func (a *App) GenerateRouteForStations(languageTag language.Tag, origin, destination string) (message.Response, error) {
	// find station ids
	originStationId, err := a.stationNameResolver.FindStationIdByApproximateName(origin)
	if err != nil {
		return message.Response{}, fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [origin='%s']: %w", origin, err)
	}
	destinationStationId, err := a.stationNameResolver.FindStationIdByApproximateName(destination)
	if err != nil {
		return message.Response{}, fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [destination='%s']: %w", destination, err)
	}
	// check blacklisted stations
	if isBlacklisted, stations := a.blackList.CheckBlackList(originStationId, destinationStationId); isBlacklisted {
		_message := a.render.BlackListedStations(languageTag, stations...)
		return _message, nil
	}
	// find route
	routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId)
	// get callbacks
	var (
		updateCallback  = a.callback.GenerateUpdateCallbackData(origin, destination)
		reverseCallback = a.callback.GenerateReverseRouteCallbackData(origin, destination)
	)
	// render message
	var _message message.Response
	if isDirect {
		_message = a.render.DirectRoutes(languageTag, routes, updateCallback, reverseCallback)
	} else {
		_message = a.render.TransferRoutes(languageTag, routes, originStationId, a.transferStationId, destinationStationId, updateCallback, reverseCallback)
	}
	return _message, nil
}
