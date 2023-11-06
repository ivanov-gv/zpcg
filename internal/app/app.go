package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"zpcg/internal/model"
	"zpcg/internal/name"
	"zpcg/internal/render"

	"zpcg/internal/pathfinder"
	"zpcg/internal/transfer"
)

func NewApp(timetableFilepath string) (*App, error) {
	timetable, err := transfer.ImportTimetable(timetableFilepath)
	if err != nil {
		return nil, fmt.Errorf("transfer.ImportTimetable: %w", err)
	}
	finder := pathfinder.NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	stationNameResolver := name.NewStationNameResolver(timetable.UnifiedStationNameToStationIdMap, timetable.UnifiedStationNameList)
	_render := render.NewRender(timetable.StationIdToStaionMap, timetable.TrainIdToTrainInfoMap)
	return &App{
		finder:              finder,
		stationNameResolver: stationNameResolver,
		render:              _render,
		transferStationId:   timetable.TransferStationId,
	}, nil
}

type App struct {
	finder              *pathfinder.PathFinder
	stationNameResolver *name.StationNameResolver
	render              *render.Render
	transferStationId   model.StationId
}

func (a *App) HandleUpdate(update tgbotapi.Update) (answer tgbotapi.MessageConfig, isNotEmpty bool) {
	const logFmf = "handleUpdate: "
	if update.Message == nil || update.Message.From == nil {
		return tgbotapi.MessageConfig{}, false
	}
	// process message
	message := update.Message
	log.Println(logFmf, "got new message: ", message.From.FirstName, message.From.UserName, message.Text)
	// generate answer
	var answerText, parseMode string
	switch {
	case strings.HasPrefix(message.Text, "/"):
		// got command - send start message
		answerText, parseMode = a.StartMessage()
	case strings.Contains(message.Text, ","):
		// got message with stations - send a timetable
		originStation, destinationStation, _ := strings.Cut(message.Text, ",")
		answerText, parseMode = a.GenerateRoute(originStation, destinationStation)
	default:
		// error otherwise
		answerText, parseMode = a.ErrorMessage()
	}
	// create message and return
	answer = tgbotapi.NewMessage(message.Chat.ID, answerText)
	answer.ParseMode = parseMode
	return answer, true
}

func (a *App) GenerateRoute(origin, destination string) (message, parseMode string) {
	const logfmt = "GenerateRoute: "
	// find station ids
	originStationId, err := a.stationNameResolver.FindStationIdByApproximateName(origin)
	if err != nil {
		log.Println(logfmt, "err", err, "origin", origin)
		return a.render.ErrorMessage()
	}
	destinationStationId, err := a.stationNameResolver.FindStationIdByApproximateName(destination)
	if err != nil {
		log.Println(logfmt, "err", err, "destination", destination)
		return a.render.ErrorMessage()
	}
	// find route
	routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId)
	// render message
	if isDirect {
		return a.render.DirectRoutes(routes)
	}
	return a.render.TransferRoutes(routes, originStationId, a.transferStationId, destinationStationId)
}

func (a *App) StartMessage() (message, parseMode string) {
	return a.render.StartMessage()
}

func (a *App) ErrorMessage() (message, parseMode string) {
	return a.render.ErrorMessage()
}
