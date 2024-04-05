package render

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/text/language"

	"zpcg/internal/model"
)

func NewRender(stationsMap map[model.StationId]model.Station,
	trainsMap map[model.TrainId]model.TrainInfo) *Render {
	return &Render{
		stationsMap: stationsMap,
		trainsMap:   trainsMap,
	}
}

type Render struct {
	stationsMap map[model.StationId]model.Station
	trainsMap   map[model.TrainId]model.TrainInfo
}

const (
	timetableLinkAnchor = "#tab3"
	stationsDelimiter   = "\\>"
	timeLayout          = "15:04"
)

func inlineButtonWithOfficialTimetableUrl(languageCode language.Tag, origin, destination string) model.InlineButton {
	var text string
	switch languageCode {
	case language.Russian:
		text = OfficialTimetableUrlTextRu
	default:
		text = OfficialTimetableUrlText
	}
	return model.InlineButton{
		Text: text,
		Url:  getUrlToTimetable(origin, destination, time.Time{}),
	}
}

func (r *Render) DirectRoutes(languageTag language.Tag, paths []model.Path) model.Response {
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
		line := fmt.Sprintf("[%04d](%s%s)` %s %s %s `",
			train.TrainId, train.TimetableUrl, timetableLinkAnchor,
			path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		lines = append(lines, line)
	}
	// add inline keyboard with url to the official website
	inlineKeyboard := [][]model.InlineButton{{inlineButtonWithOfficialTimetableUrl(languageTag, origin.Name, destination.Name)}}
	return model.Response{
		Text:           strings.Join(lines, "\n"),
		ParseMode:      tgbotapi.ModeMarkdownV2,
		InlineKeyboard: inlineKeyboard,
	}
}

func (r *Render) TransferRoutes(languageTag language.Tag, paths []model.Path,
	originId, transferId, destinationId model.StationId) model.Response {
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
			line = fmt.Sprintf("[%04d](%s%s)` %s %s %s `",
				train.TrainId, train.TimetableUrl, timetableLinkAnchor,
				path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		} else {
			// right side of the table - Transfer Stop -> B
			line = fmt.Sprintf("[%04d](%s%s)`         %s %s %s `",
				train.TrainId, train.TimetableUrl, timetableLinkAnchor,
				path.Origin.Departure.Format(timeLayout), stationsDelimiter, path.Destination.Arrival.Format(timeLayout))
		}
		lines = append(lines, line)
	}
	// add inline keyboard with url to the official website
	inlineKeyboard := [][]model.InlineButton{
		{
			inlineButtonWithOfficialTimetableUrl(languageTag, origin.Name, transfer.Name),
			inlineButtonWithOfficialTimetableUrl(languageTag, transfer.Name, destination.Name),
		},
	}
	return model.Response{
		Text:           strings.Join(lines, "\n"),
		ParseMode:      tgbotapi.ModeMarkdownV2,
		InlineKeyboard: inlineKeyboard,
	}
}

func (r *Render) ErrorMessage(languageCode language.Tag) model.Response {
	switch languageCode {
	case language.Russian:
		return model.Response{Text: ErrorMessageRu}
	default:
		return model.Response{Text: ErrorMessageDefault}
	}
}

func (r *Render) StartMessage(languageCode language.Tag) model.Response {
	switch languageCode {
	case language.Russian:
		return model.Response{Text: StartMessageRu, ParseMode: tgbotapi.ModeMarkdownV2}
	default:
		return model.Response{Text: StartMessageDefault, ParseMode: tgbotapi.ModeMarkdownV2}
	}
}

func (r *Render) BlackListedStations(languageCode language.Tag, stations ...model.BlackListedStation) model.Response {
	var lines []string
	for _, station := range stations {
		var line string
		if customMessage, ok := station.LanguageTagToCustomErrorMessageMap[languageCode.String()]; ok {
			line = customMessage
		} else {
			line = fmt.Sprintf("%s: %s", station.Name, StationDoesNotExistMessageMap[languageCode])
		}
		lines = append(lines, line)
	}
	lines = append(lines,
		"", // empty line
		StationDoesNotExistMessageSuffixMap[languageCode])
	return model.Response{
		Text:      strings.Join(lines, "\n"),
		ParseMode: "",
	}
}

func ParseLanguageTag(languageCode string) language.Tag {
	tag, err := language.Parse(languageCode)
	if err != nil {
		return DefaultLanguageTag
	}
	return tag
}
