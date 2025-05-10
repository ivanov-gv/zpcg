package render

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/ivanov-gv/zpcg/internal/model/message"
	model_render "github.com/ivanov-gv/zpcg/internal/model/render"
	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

func NewRender(stationsMap map[timetable.StationId]timetable.Station,
	trainsMap map[timetable.TrainId]timetable.TrainInfo) *Render {
	return &Render{
		stationsMap: stationsMap,
		trainsMap:   trainsMap,
	}
}

type Render struct {
	stationsMap map[timetable.StationId]timetable.Station
	trainsMap   map[timetable.TrainId]timetable.TrainInfo
}

const (
	timetableLinkAnchor = "#tab3"
	stationsDelimiter   = "\\>"
	timeLayout          = "15:04"
)

func inlineButtonWithUpdateCallback(languageCode language.Tag, currentDate time.Time, updateCallback string) message.InlineButton {
	return message.InlineButton{
		Type: message.CallbackInlineButtonType,
		Text: fmt.Sprintf("🔄 %s", Date(languageCode, currentDate)),
		Callback: message.CallbackButton{
			Data: updateCallback,
		},
	}
}

func inlineButtonWithReverseCallback(languageTag language.Tag, reverseCallback string) message.InlineButton {
	return message.InlineButton{
		Type: message.CallbackInlineButtonType,
		Text: fmt.Sprintf("↪️ %s", GetMessage(model_render.ReverseRouteInlineButtonTextMap, languageTag)),
		Callback: message.CallbackButton{
			Data: reverseCallback,
		},
	}
}

func (r *Render) DirectRoutes(languageTag language.Tag, paths []timetable.Path, currentDate time.Time,
	updateCallback, reverseCallback string) message.Response {
	// render each line for the result message
	var lines []string
	// render header
	origin := r.stationsMap[paths[0].Origin.Id]
	destination := r.stationsMap[paths[0].Destination.Id]
	header := fmt.Sprintf("`%10s %s %s`", origin.Name, stationsDelimiter, destination.Name)
	// add prefix to align header with table content
	lines = append(lines, header)
	// render the rest of the message
	for _, path := range paths {
		train := r.trainsMap[path.TrainId]
		line := fmt.Sprintf("`%04d`` %s %s %s `",
			train.TrainId,
			path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		lines = append(lines, line)
	}
	// render footer
	footer := moreDetailsLink(languageTag, origin.Name, destination.Name)
	lines = append(lines, "", footer)
	// add inline keyboard with url to the official website
	inlineKeyboard := [][]message.InlineButton{
		{inlineButtonWithUpdateCallback(languageTag, currentDate, updateCallback), inlineButtonWithReverseCallback(languageTag, reverseCallback)},
	}
	return message.Response{
		Text:           strings.Join(lines, "\n"),
		ParseMode:      message.ModeMarkdownV2,
		InlineKeyboard: inlineKeyboard,
	}
}

func (r *Render) TransferRoutes(languageTag language.Tag, paths []timetable.Path, currentDate time.Time,
	originId, transferId, destinationId timetable.StationId, updateCallback, reverseCallback string) message.Response {
	// render each line for the result message
	var lines []string
	// render header
	origin := r.stationsMap[originId]
	transfer := r.stationsMap[transferId]
	destination := r.stationsMap[destinationId]
	header := fmt.Sprintf("`%s %s %s %s %s`",
		origin.Name, stationsDelimiter, transfer.Name, stationsDelimiter, destination.Name)
	// add header
	lines = append(lines, header)
	// add other lines
	for _, path := range paths {
		var (
			train = r.trainsMap[path.TrainId]
			line  string
		)
		if path.Origin.Id == originId && path.Destination.Id == transferId {
			// left side of the table - A -> Transfer Stop
			line = fmt.Sprintf("`%04d`` %s %s %s `",
				train.TrainId,
				path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		} else {
			// right side of the table - Transfer Stop -> B
			line = fmt.Sprintf("`%04d``         %s %s %s `",
				train.TrainId,
				path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		}
		lines = append(lines, line)
	}
	// footer
	footer := moreDetailsLink(languageTag, origin.Name, destination.Name)
	lines = append(lines, "", footer)
	// add inline keyboard with url to the official website
	inlineKeyboard := [][]message.InlineButton{
		{inlineButtonWithUpdateCallback(languageTag, currentDate, updateCallback), inlineButtonWithReverseCallback(languageTag, reverseCallback)},
	}
	return message.Response{
		Text:           strings.Join(lines, "\n"),
		ParseMode:      message.ModeMarkdownV2,
		InlineKeyboard: inlineKeyboard,
	}
}

func (r *Render) ErrorMessage(languageTag language.Tag) message.Response {
	return message.Response{
		Text:      GetMessage(model_render.ErrorMessageMap, languageTag),
		ParseMode: message.ModeNone,
	}
}

func (r *Render) Command(languageTag language.Tag, command string) message.Response {
	switch model_render.BotCommand(strings.TrimSpace(command)) {
	default:
		fallthrough // default option - /start
	case model_render.BotCommandStart:
		return r.startMessage(languageTag)
	case model_render.BotCommandHelp:
		return r.helpMessage(languageTag)
	case model_render.BotCommandMap:
		return r.mapMessage(languageTag)
	case model_render.BotCommandAbout:
		return r.aboutMessage(languageTag)
	}
}

func (r *Render) startMessage(languageTag language.Tag) message.Response {
	return message.Response{
		Text:      GetMessage(model_render.StartMessageMap, languageTag),
		ParseMode: message.ModeMarkdownV2,
	}
}

func (r *Render) helpMessage(languageTag language.Tag) message.Response {
	return message.Response{
		Text:      GetMessage(model_render.HelpMessageMap, languageTag),
		ParseMode: message.ModeNone,
	}
}

func (r *Render) aboutMessage(languageTag language.Tag) message.Response {
	return message.Response{
		Text:      GetMessage(model_render.AboutMessageMap, languageTag),
		ParseMode: message.ModeNone,
	}
}

func (r *Render) mapMessage(languageTag language.Tag) message.Response {
	return message.Response{
		Text: GetMessage(model_render.MapMessageMap, languageTag) + "\n" +
			model_render.GoogleMapWithAllStations,
		ParseMode: message.ModeNone,
	}
}

func (r *Render) BlackListedStations(languageTag language.Tag, stations ...timetable.BlackListedStation) message.Response {
	var lines []string
	for _, station := range stations {
		var line string
		if customMessage, ok := station.LanguageTagToCustomErrorMessageMap[languageTag]; ok {
			line = customMessage
		} else if defaultMessage, ok := station.LanguageTagToCustomErrorMessageMap[model_render.DefaultLanguageTag]; ok {
			line = defaultMessage
		} else {
			line = fmt.Sprintf("%s: %s", station.Name, GetMessage(model_render.StationDoesNotExistMessageMap, languageTag))
		}
		lines = append(lines, line)
	}

	return message.Response{
		Text:      strings.Join(lines, "\n"),
		ParseMode: message.ModeNone,
		InlineKeyboard: [][]message.InlineButton{
			{
				{
					Type: message.UrlInlineButtonType,
					Text: GetMessage(model_render.RailwayMapButtonTextMap, languageTag),
					Url:  message.UrlButton{Url: model_render.GoogleMapWithAllStations},
				},
			},
		},
	}
}

func (r *Render) AlertUpdateNotificationText(languageTag language.Tag) string {
	return GetMessage(model_render.AlertUpdateNotificationTextMap, languageTag)
}

func (r *Render) SimpleUpdateNotificationText(languageTag language.Tag) string {
	return GetMessage(model_render.SimpleUpdateNotificationTextMap, languageTag)
}

func moreDetailsLink(languageTag language.Tag, origin, destination string) string {
	return fmt.Sprintf("[%s](%s)",
		GetMessage(model_render.OfficialTimetableUrlTextMap, languageTag),
		getUrlToTimetable(origin, destination, time.Time{}))
}

func ParseLanguageTag(languageCode string) language.Tag {
	tag, err := language.Parse(languageCode)
	if err != nil {
		return model_render.DefaultLanguageTag
	}
	return tag
}

func GetMessage[T any](_map map[language.Tag]T, tag language.Tag) T {
	if _message, ok := _map[tag]; ok {
		return _message
	}
	return _map[model_render.DefaultLanguageTag]
}
