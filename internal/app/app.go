package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"
	"golang.org/x/text/language"

	callback_model "github.com/ivanov-gv/zpcg/internal/model/callback"
	message_model "github.com/ivanov-gv/zpcg/internal/model/message"
	"github.com/ivanov-gv/zpcg/internal/pkg/timetable_export"
	"github.com/ivanov-gv/zpcg/internal/service/callback"
	"github.com/ivanov-gv/zpcg/internal/service/date"
	"github.com/ivanov-gv/zpcg/internal/service/message_render"
	"github.com/ivanov-gv/zpcg/internal/service/pathfinder"
	"github.com/ivanov-gv/zpcg/internal/service/station_blacklist"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
	"github.com/ivanov-gv/zpcg/internal/service/timetable"
)

func NewApp(opts ...date.Option) (*App, error) {
	// timetable
	_timetable := timetable_export.ImportTimetable()
	// name resolver
	stationNameResolver := station_name_resolver.NewStationNameResolver(_timetable.UnifiedStationNameToStationIdMap,
		_timetable.UnifiedStationNameList, _timetable.StationIdToStationMap)
	// render
	render := message_render.NewRender(_timetable.StationIdToStationMap)
	// station_blacklist
	blackList := station_blacklist.NewBlackListService()
	// date
	dateService := date.NewDateService(context.TODO(), opts...)
	// timetable
	timetableService, err := timetable.New(_timetable, dateService)
	if err != nil {
		return nil, fmt.Errorf("timetable.New: %w", err)
	}
	// finder
	finder := pathfinder.NewPathFinder(timetableService)

	// complete app
	return &App{
		finder:              finder,
		stationNameResolver: stationNameResolver,
		render:              render,
		blackList:           blackList,
		dateService:         dateService,
		timetableService:    timetableService,
	}, nil
}

// App handles requests in the most general form. It has to know nothing about specific messenger, api, libs etc.
// It simply takes message and generates response. Contains whole business logic and business logic only
type App struct {
	finder              *pathfinder.PathFinder
	stationNameResolver *station_name_resolver.StationNameResolver
	render              *message_render.Render
	blackList           *station_blacklist.BlackListService
	callback            *callback.CallbackService
	dateService         *date.DateService
	timetableService    *timetable.TimetableService
}

func (a *App) HandleUpdate(update message_model.Update) (responseWithChatIds message_model.ResponseWithChatId, warning error) {
	switch update.Type {
	case message_model.MessageUpdateType:
		return a.HandleMessage(update.Message)
	case message_model.CallbackUpdateType:
		return a.HandleCallback(update.Callback)
	default: // including model.UnsupportedUpdateType
		return message_model.ResponseWithChatId{}, nil
	}
}

func (a *App) HandleCallback(callbackMessage message_model.Callback) (responseWithChatIds message_model.ResponseWithChatId, warning error) {
	var languageTag = message_render.ParseLanguageTag(callbackMessage.From.LanguageCode)
	defer func() { // we have to answer the callback anyway
		responseWithChatIds.ChatId = callbackMessage.ChatId
		responseWithChatIds.AnswerCallback.CallbackQueryId = callbackMessage.Id
	}()
	_callback, err := a.callback.ParseCallback(callbackMessage.Data)
	if err != nil {
		return message_model.ResponseWithChatId{}, fmt.Errorf("failed to parse callback [data='%s']: %w", callbackMessage.Data, err)
	}
	switch _callback.Type {
	case callback_model.UpdateType:
		var data = _callback.UpdateData
		// check date
		if data.Date.Equal(a.dateService.CurrentDateAsTime()) {
			return message_model.ResponseWithChatId{
				AnswerCallback: message_model.ToAnswerCallbackQuery{
					Text:      a.render.AlertUpdateNotificationText(languageTag, a.timetableService.Season().UpdateButtonAlertText),
					ShowAlert: true,
				},
			}, nil
		}
		// update outdated
		response, err := a.GenerateRouteForStations(languageTag, data.Origin, data.Destination)
		if err != nil {
			return message_model.ResponseWithChatId{},
				fmt.Errorf("GenerateRouteForStations [lang=%s, origin=%s, destination=%s] : %w",
					languageTag, data.Origin, data.Destination, err)
		}
		update := []message_model.ToUpdate{
			{
				MessageId:       callbackMessage.Message.Id,
				InlineMessageId: callbackMessage.InlineMessageId,
				Response:        response,
			},
		}
		return message_model.ResponseWithChatId{
			ChatId:         callbackMessage.ChatId,
			Update:         update,
			AnswerCallback: message_model.ToAnswerCallbackQuery{Text: a.render.SimpleUpdateNotificationText(languageTag)},
		}, nil
	case callback_model.ReverseRouteType:
		var data = _callback.ReverseRouteData
		response, err := a.GenerateRouteForStations(languageTag, data.Destination, data.Origin)
		if err != nil {
			return message_model.ResponseWithChatId{},
				fmt.Errorf("GenerateRouteForStations [lang=%s, origin=%s, destination=%s] : %w",
					languageTag, data.Destination, data.Origin, err)
		}
		send := []message_model.ToSend{
			{
				Response: response,
			},
		}
		return message_model.ResponseWithChatId{ChatId: callbackMessage.ChatId, Send: send}, nil
	default:
		return message_model.ResponseWithChatId{},
			fmt.Errorf("unknown callback type(%s) [data='%s']", _callback.Type, callbackMessage.Data)
	}
}

func (a *App) HandleMessage(_message message_model.Message) (responseWithChatIds message_model.ResponseWithChatId, warning error) {
	languageTag := message_render.ParseLanguageTag(_message.From.LanguageCode)
	var (
		response message_model.Response
		err      error
	)
	switch {
	case strings.HasPrefix(_message.Text, "/"):
		// got command
		response = a.render.Command(languageTag, _message.Text)
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
		switch {
		case errors.Is(err, station_name_resolver.ErrNoMatchesFound):
			var stationInput string
			if err, ok := errors.AsType[oops.OopsError](err); ok {
				stationInput = err.Public()
			}
			response = a.render.StationNotFoundMessage(languageTag, stationInput)
		case errors.Is(err, pathfinder.ErrNoRoutesFound):
			response = a.render.NoRoutesInThisSeasonMessage(languageTag, a.timetableService.Season().NoTrainsWarning)
		default:
			response = a.render.ErrorMessage(languageTag)
		}
	}
	send := []message_model.ToSend{
		{
			Response: response,
		},
	}
	return message_model.ResponseWithChatId{
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
	if len(stations) < 2 { //nolint:mnd // need at least origin and destination
		return "", "", fmt.Errorf("not enough stations provided: %s", inputWithProperDelimiter)
	}
	return stations[0], lo.Must(lo.Last(stations)), nil
}

func (a *App) GenerateRoute(languageTag language.Tag, input string) (message_model.Response, error) {
	origin, destination, err := parseInputStations(input)
	if err != nil {
		return message_model.Response{}, fmt.Errorf("parseInputStations: %w", err)
	}
	return a.GenerateRouteForStations(languageTag, origin, destination)
}

func (a *App) GenerateRouteForStations(languageTag language.Tag, originInput, destinationInput string) (message_model.Response, error) {
	// find station ids
	originStationId, originName, err := a.stationNameResolver.FindStationIdByApproximateName(originInput)
	if err != nil {
		return message_model.Response{}, oops.Public(originInput).
			Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
				"can't find station name [origin='%s']: %w", originInput, err)
	}
	destinationStationId, destinationName, err := a.stationNameResolver.FindStationIdByApproximateName(destinationInput)
	if err != nil {
		return message_model.Response{}, oops.Public(destinationInput).
			Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
				"can't find station name [destination='%s']: %w", destinationInput, err)
	}
	// check blacklisted stations
	if isBlacklisted, stations := a.blackList.CheckBlackList(originStationId, destinationStationId); isBlacklisted {
		_message := a.render.BlackListedStations(languageTag, stations...)
		return _message, nil
	}
	// find route
	routes, isDirect, err := a.finder.FindRoutes(originStationId, destinationStationId)
	if err != nil {
		return message_model.Response{}, fmt.Errorf("a.finder.FindRoutes [origin='%s', destination='%s']: %w",
			originInput, destinationInput, err)
	}
	// get callbacks
	var (
		updateCallback  = a.callback.GenerateUpdateCallbackData(originName, destinationName, a.dateService.CurrentDateAsShortString())
		reverseCallback = a.callback.GenerateReverseRouteCallbackData(originName, destinationName)
	)
	// render message
	var _message message_model.Response
	if isDirect {
		_message = a.render.DirectRoutes(languageTag, a.timetableService.Season(), routes, a.dateService.CurrentDateAsTime(), updateCallback, reverseCallback)
	} else {
		_message = a.render.TransferRoutes(languageTag, a.timetableService.Season(), routes, a.dateService.CurrentDateAsTime(),
			originStationId, a.timetableService.TransferStationId(), destinationStationId,
			updateCallback, reverseCallback)
	}
	return _message, nil
}
