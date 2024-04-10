package app

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/text/language"

	"zpcg/internal/model"
	"zpcg/internal/service/blacklist"
	"zpcg/internal/service/name"
	"zpcg/internal/service/pathfinder"
	"zpcg/internal/service/render"
	"zpcg/internal/service/transfer"
)

func NewApp() (*App, error) {
	// timetable
	timetable := transfer.ImportTimetable()
	// finder
	finder := pathfinder.NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	// name resolver
	stationNameResolver := name.NewStationNameResolver(timetable.UnifiedStationNameToStationIdMap, timetable.UnifiedStationNameList)
	// render
	_render := render.NewRender(timetable.StationIdToStationMap, timetable.TrainIdToTrainInfoMap)
	// blacklist
	blackList := blacklist.NewBlackListService()

	// complete app
	return &App{
		finder:              finder,
		stationNameResolver: stationNameResolver,
		render:              _render,
		blackList:           blackList,
		transferStationId:   timetable.TransferStationId,
	}, nil
}

// App handles requests in the most general form. It has to know nothing about specific messenger, api, libs etc.
// It simply takes message and generates response. Contains whole business logic and business logic only
type App struct {
	finder              *pathfinder.PathFinder
	stationNameResolver *name.StationNameResolver
	render              *render.Render
	blackList           *blacklist.BlackListService
	transferStationId   model.StationId
}

func (a *App) HandleUpdate(update model.Update) (responseWithChatIds []model.ResponseWithChatId, warning error) {
	if !update.Message.IsFilled || !update.Message.From.IsFilled {
		return nil, nil
	}
	// process message
	message := update.Message
	languageTag := render.ParseLanguageTag(update.Message.From.LanguageCode)
	// generate output
	var (
		response model.Response
		err      error
	)
	switch {
	case strings.HasPrefix(message.Text, "/"):
		// got command - send start message
		response = a.render.StartMessage(languageTag)
	case strings.ContainsAny(message.Text, string(lo.SpecialCharset)):
		// got message with stations - send a timetable
		response, err = a.GenerateRoute(languageTag, message.Text)
		if err != nil {
			err = fmt.Errorf("a.GenerateRoute: %w", err)
		}
	default:
		// error otherwise
		err = fmt.Errorf("unknown message type: %s", message.Text)
	}
	// handle error
	if err != nil {
		response = a.render.ErrorMessage(languageTag)
	}
	output := model.ResponseWithChatId{
		Response: response,
		ChatId:   message.ChatId,
	}
	return []model.ResponseWithChatId{output}, err
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

func (a *App) GenerateRoute(languageTag language.Tag, input string) (model.Response, error) {
	origin, destination, err := parseInputStations(input)
	if err != nil {
		return model.Response{}, fmt.Errorf("parseInputStations: %w", err)
	}
	// find station ids
	originStationId, err := a.stationNameResolver.FindStationIdByApproximateName(origin)
	if err != nil {
		return model.Response{}, fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [origin='%s']: %w", origin, err)
	}
	destinationStationId, err := a.stationNameResolver.FindStationIdByApproximateName(destination)
	if err != nil {
		return model.Response{}, fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [destination='%s']: %w", destination, err)
	}
	// check blacklisted stations
	if isBlacklisted, stations := a.blackList.CheckBlackList(originStationId, destinationStationId); isBlacklisted {
		message := a.render.BlackListedStations(languageTag, stations...)
		return message, nil
	}
	// find route
	routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId)
	// render message
	if isDirect {
		message := a.render.DirectRoutes(languageTag, routes)
		return message, nil
	}
	// if !isDirect - transfer route
	message := a.render.TransferRoutes(languageTag, routes, originStationId, a.transferStationId, destinationStationId)
	return message, nil
}
