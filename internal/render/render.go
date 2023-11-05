package render

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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

func (r *Render) DirectRoutes(paths []model.Path) (message, parseMode string) {
	// render each line for the result message
	var lines []string
	// render header
	origin := r.stationsMap[paths[0].Origin.Id]
	destination := r.stationsMap[paths[0].Destination.Id]
	header := fmt.Sprintf("`%6s \\-\\> %s`", origin.Name, destination.Name)
	// add prefix to align header with table content
	lines = append(lines, header)
	// render the rest of the message
	for _, path := range paths {
		train := r.trainsMap[path.TrainId]
		line := fmt.Sprintf("[%04d](%s#tab3) ` %s \\-\\> %s `",
			train.TrainId, train.TimetableUrl,
			path.Origin.Departure.Format("15:04"),
			path.Destination.Arrival.Format("15:04"))
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
	header := fmt.Sprintf("` %s \\-\\> %s \\-\\> %s `", origin.Name, transfer.Name, destination.Name)
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
			line = fmt.Sprintf("[%04d](%s#tab3) ` %s \\-\\> %s `",
				train.TrainId, train.TimetableUrl,
				path.Origin.Departure.Format("15:04"),
				path.Destination.Arrival.Format("15:04"))
		} else {
			// right side of the table - Transfer Stop -> B
			line = fmt.Sprintf("[%04d](%s#tab3) `          %s \\-\\> %s `",
				train.TrainId, train.TimetableUrl,
				path.Origin.Departure.Format("15:04"),
				path.Destination.Arrival.Format("15:04"))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), tgbotapi.ModeMarkdownV2
}

const ErrorMessage = `Try again - two stations, separated by comma. Just like that:
Podgorica, Niksic`

func (r *Render) ErrorMessage() (message, parseMode string) {
	return ErrorMessage, ""
}

const StartMessage = "" +
	"*Montenegro Railways Timetable*\n" +
	"\n" +
	"Type two stations, separated by comma: \n" +
	"\n" +
	"*Podgorica, Bijelo Polje*\n" +
	"\n" +
	"And I will send the timetable for you:\n" +
	"\n" +
	"Podgorica \\-\\> Bijelo Polje\n" +
	"[6100](https://zpcg.me/details?timetable=41)  `06:20 \\-\\> 08:38` \n" +
	"\\.\\.\\.\n" +
	"\n" +
	"Don't know how to write the station name the right way? Don't worry, just type, I will do the rest\\.\n" +
	"\n" +
	"*PADAGORNICCCCA , BELO POLLLLLE*\n" +
	"\n" +
	"Podgorica \\-\\> Bijelo Polje\n" +
	"[6100](https://zpcg.me/details?timetable=41)  `06:20 \\-\\> 08:38` \n" +
	"\\.\\.\\.\n" +
	"\n" +
	"Now it's your turn\\!"

func (r *Render) StartMessage() (message, parseMode string) {
	return StartMessage, tgbotapi.ModeMarkdownV2
}
