package app

import (
	"fmt"
	"io"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/language"

	"zpcg/internal/model"
	"zpcg/internal/name"
	"zpcg/internal/pathfinder"
	"zpcg/internal/render"
	"zpcg/internal/transfer"
)

func NewApp(timetableReader io.Reader) (*App, error) {
	timetable, err := transfer.ImportTimetableFromReader(timetableReader)
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
	const logFmt = "handleUpdate: "
	if update.Message == nil || update.Message.From == nil {
		return tgbotapi.MessageConfig{}, false
	}
	// process message
	message := update.Message
	languageTag := parseLanguageTag(update.SentFrom().LanguageCode)
	log.Println(logFmt, "got new message: ", message.From.FirstName, message.From.UserName,
		update.SentFrom().LanguageCode, languageTag.String(),
		message.Text)
	// generate answer
	var (
		answerText, parseMode string
		err                   error
	)
	switch {
	case strings.HasPrefix(message.Text, "/"):
		// got command - send start message
		answerText, parseMode = a.StartMessage(languageTag)
	case strings.Contains(message.Text, ","):
		// got message with stations - send a timetable
		originStation, destinationStation, _ := strings.Cut(message.Text, ",")
		answerText, parseMode, err = a.GenerateRoute(originStation, destinationStation)
		if err != nil {
			err = fmt.Errorf("a.GenerateRoute: %w", err)
		}
	default:
		// error otherwise
		err = fmt.Errorf("unknown message type: %s", message.Text)
	}
	// handle error
	if err != nil {
		log.Println(logFmt, "err", err)
		answerText, parseMode = a.ErrorMessage(languageTag)
	}
	// create message and return
	answer = tgbotapi.NewMessage(message.Chat.ID, answerText)
	answer.ParseMode = parseMode
	return answer, true
}

func (a *App) GenerateRoute(origin, destination string) (message, parseMode string, err error) {
	// find station ids
	originStationId, err := a.stationNameResolver.FindStationIdByApproximateName(origin)
	if err != nil {
		return "", "", fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [origin='%s']: %w", origin, err)
	}
	destinationStationId, err := a.stationNameResolver.FindStationIdByApproximateName(destination)
	if err != nil {
		return "", "", fmt.Errorf("a.stationNameResolver.FindStationIdByApproximateName: "+
			"can't find station name [destination='%s']: %w", destination, err)
	}
	// find route
	routes, isDirect := a.finder.FindRoutes(originStationId, destinationStationId)
	// render message
	if isDirect {
		message, parseMode = a.render.DirectRoutes(routes)
		return message, parseMode, nil
	}
	// if !isDirect - transfer route
	message, parseMode = a.render.TransferRoutes(routes, originStationId, a.transferStationId, destinationStationId)
	return message, parseMode, nil
}

func parseLanguageTag(languageCode string) language.Tag {
	tag, err := language.Parse(languageCode)
	if err != nil {
		return render.DefaultLanguageTag
	}
	return tag
}

func (a *App) StartMessage(languageTag language.Tag) (message, parseMode string) {
	return a.render.StartMessage(languageTag)
}

func (a *App) ErrorMessage(languageTag language.Tag) (message, parseMode string) {
	return a.render.ErrorMessage(languageTag)
}
