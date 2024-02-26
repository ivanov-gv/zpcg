package render

import (
	"fmt"
	"strings"

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

func (r *Render) DirectRoutes(paths []model.Path) (message, parseMode string) {
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
			path.Origin.Arrival.Format(timeLayout), stationsDelimiter, path.Destination.Departure.Format(timeLayout))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), tgbotapi.ModeMarkdownV2
}

func (r *Render) TransferRoutes(paths []model.Path, originId, transferId, destinationId model.StationId) (message, parseMode string) {
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
				path.Origin.Arrival.Format(timeLayout), stationsDelimiter, path.Destination.Departure.Format(timeLayout))
		} else {
			// right side of the table - Transfer Stop -> B
			line = fmt.Sprintf("[%04d](%s%s)`         %s %s %s `",
				train.TrainId, train.TimetableUrl, timetableLinkAnchor,
				path.Origin.Arrival.Format(timeLayout), stationsDelimiter, path.Destination.Departure.Format(timeLayout))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), tgbotapi.ModeMarkdownV2
}

func (r *Render) ErrorMessage(languageCode language.Tag) (message, parseMode string) {
	switch languageCode {
	case language.Russian:
		return ErrorMessageRu, ""
	default:
		return ErrorMessageDefault, ""
	}
}

func (r *Render) StartMessage(languageCode language.Tag) (message, parseMode string) {
	switch languageCode {
	case language.Russian:
		return StartMessageRu, tgbotapi.ModeMarkdownV2
	default:
		return StartMessageDefault, tgbotapi.ModeMarkdownV2
	}
}

//func (r *Render) BlackListedStations(languageCode language.Tag, station1, station2 model.BlackListedStation) (message, parseMode string) {
//	switch languageCode { // TODO: render
//	case language.Russian:
//	default:
//
//	}
//}
