package app

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"

	"zpcg/internal/config"
	"zpcg/internal/model"
	"zpcg/internal/service/name"
	"zpcg/internal/service/pathfinder"
	"zpcg/internal/service/render"
	"zpcg/internal/service/transfer"
	"zpcg/resources"
)

func NewApp(_config config.Config) (*App, error) {
	// timetable reader
	timetableReader, err := resources.FS.Open(_config.TimetableGobFileName)
	if err != nil {
		return nil, fmt.Errorf("fs.Open: %w", err)
	}
	// timetable
	timetable, err := transfer.ImportTimetableFromReader(timetableReader)
	if err != nil {
		return nil, fmt.Errorf("transfer.ImportTimetable: %w", err)
	}
	// finder
	finder := pathfinder.NewPathFinder(timetable.StationIdToTrainIdSet, timetable.TrainIdToStationMap, timetable.TransferStationId)
	// name resolver
	stationNameResolver := name.NewStationNameResolver(timetable.UnifiedStationNameToStationIdMap, timetable.UnifiedStationNameList)
	// render
	_render := render.NewRender(timetable.StationIdToStaionMap, timetable.TrainIdToTrainInfoMap)

	// complete app
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

const (
	stationsDelimitersSet = stationsDelimiterComma +
		stationsDelimiterArrow

	stationsDelimiterComma = ","
	stationsDelimiterArrow = ">"
)

func (a *App) HandleUpdate(update tgbotapi.Update) (answer tgbotapi.MessageConfig, isNotEmpty bool) {
	if update.Message == nil || update.Message.From == nil {
		return tgbotapi.MessageConfig{}, false
	}
	// process message
	message := update.Message
	languageTag := parseLanguageTag(update.SentFrom().LanguageCode)
	// generate answer
	var (
		answerText, parseMode string
		err                   error
	)
	switch {
	case strings.HasPrefix(message.Text, "/"):
		// got command - send start message
		answerText, parseMode = a.StartMessage(languageTag)
	case strings.ContainsAny(message.Text, stationsDelimitersSet):
		// got message with stations - send a timetable
		var (
			originStation, destinationStation string
		)
		// parse stations
		if _origin, _destination, found := strings.Cut(message.Text, stationsDelimiterComma); found {
			originStation, destinationStation = _origin, _destination
		} else if _origin, _destination, found := strings.Cut(message.Text, stationsDelimiterArrow); found {
			originStation, destinationStation = _origin, _destination
		}
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
		answerText, parseMode = a.ErrorMessage(languageTag)
	}
	// create message and return
	answer = tgbotapi.NewMessage(message.Chat.ID, answerText)
	answer.ParseMode = parseMode
	answer.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false) // TODO: removes keyboard for the users who has it from the older versions of the @monterails_bot
	// log answer
	logTrace(update, answer, languageTag, err)
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

func logTrace(update tgbotapi.Update, answer tgbotapi.MessageConfig, languageTag language.Tag, err error) {
	const logFmt = "handleUpdate: %s"
	var (
		message  = update.Message
		logEvent *zerolog.Event
	)
	// set level
	if err != nil {
		logEvent = log.Warn().Err(err)
	} else {
		logEvent = log.Trace()
	}
	// get first 2 lines of the answer
	answerLines := strings.Split(answer.Text, "\n")
	if len(answerLines) > 2 {
		answerLines = answerLines[:2]
	}
	answerShort := strings.Join(answerLines, "\n")
	// log
	logEvent.
		Int64("chatId", message.Chat.ID).
		Str("languageCode", update.SentFrom().LanguageCode).
		Str("languageTag", languageTag.String()).
		Str("messageText", message.Text).
		Str("answerShort", answerShort).
		Msgf(logFmt, "new message handled")
}

func (a *App) StartMessage(languageTag language.Tag) (message, parseMode string) {
	return a.render.StartMessage(languageTag)
}

func (a *App) ErrorMessage(languageTag language.Tag) (message, parseMode string) {
	return a.render.ErrorMessage(languageTag)
}
