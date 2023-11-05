package render

import (
	"fmt"
	"strings"
	"unicode/utf8"
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

func (r *Render) DirectRoutes(paths []model.Path) string {
	if len(paths) == 0 {
		return "Маршрут не найден"
	}
	// render each line for the result message
	var lines []string
	// render header
	origin := r.stationsMap[paths[0].Origin.Id]
	destination := r.stationsMap[paths[0].Destination.Id]
	header := fmt.Sprintf("Поезд | %s -> %s ", origin.Name, destination.Name)
	lines = append(lines, header)
	// TODO: calculate suffix for departure time layout to align
	// render the rest of the message
	for _, path := range paths {
		train := r.trainsMap[path.TrainId]
		line := fmt.Sprintf("[% 5d](%s) | %s -> %s ",
			train.TrainId, train.TimetableUrl,
			path.Origin.Departure.Format("15:04"),
			path.Destination.Arrival.Format("15:04"))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (r *Render) TransferRoutes(paths []model.Path, originId, transferId, destinationId model.StationId) string {
	if len(paths) == 0 {
		return "Маршрут не найден"
	}
	// render each line for the result message
	var lines []string
	// render header
	origin := r.stationsMap[originId]
	transfer := r.stationsMap[transferId]
	destination := r.stationsMap[destinationId]
	header := fmt.Sprintf("Поезд | %s -> %s | %s -> %s", origin.Name, transfer.Name, transfer.Name, destination.Name)
	// calculate prefix to align text
	stopsNamesLen := utf8.RuneCountInString(origin.Name) + utf8.RuneCountInString(transfer.Name)
	prefix := strings.Repeat(" ", stopsNamesLen+4) //" -> " = 4 letters
	suffix := ""
	if stopsNamesLen > 10 {
		suffix = strings.Repeat(" ", stopsNamesLen-10) // "15:04" + "15:04" = 10 letters
	}
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
			line = fmt.Sprintf("[% 5d](%s) | %s -> %s %s|",
				train.TrainId, train.TimetableUrl,
				path.Origin.Departure.Format("15:04"),
				path.Destination.Arrival.Format("15:04"),
				suffix)
		} else {
			// right side of the table - Transfer Stop -> B
			line = fmt.Sprintf("[% 5d](%s) | %s | %s -> %s",
				train.TrainId, train.TimetableUrl,
				prefix,
				path.Origin.Departure.Format("15:04"),
				path.Destination.Arrival.Format("15:04"))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (r *Render) StationResolveError(name string) string {
	return fmt.Sprintf("Не смог понять, какую станцию вы имели ввиду: %s . Попробуйте еще раз", name)
}
