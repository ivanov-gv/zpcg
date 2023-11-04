package transfer

import (
	"encoding/gob"
	"github.com/pkg/errors"
	"os"
	"zpcg/internal/model"
	"zpcg/internal/parser/detailed_page"
)

type TimetableExportFormat struct {
	StationIdToTrainIdSet map[model.StationId]model.TrainIdSet
	TrainIdToStationMap   map[model.TrainId]model.StationIdToStationMap
	StationIdMap          map[model.StationId]string
}

func ExportTimetable(
	stationIdToTrainIdSet map[model.StationId]model.TrainIdSet,
	trainIdToStationMap map[model.TrainId]model.StationIdToStationMap,
	filename string,
) error {
	file, err := os.OpenFile(filename+".gob", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "can not open file with os.OpenFile")
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(TimetableExportFormat{
		StationIdToTrainIdSet: stationIdToTrainIdSet,
		TrainIdToStationMap:   trainIdToStationMap,
		StationIdMap:          detailed_page.GetStationIdToNameMap(),
	})
	if err != nil {
		return errors.Wrap(err, "can not encode timetable with enc.Encode")
	}
	return nil
}

func ImportTimetable(filename string) (
	map[model.StationId]model.TrainIdSet,
	map[model.TrainId]model.StationIdToStationMap,
	map[model.StationId]string,
	error,
) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not open file with os.Open")
	}
	dec := gob.NewDecoder(file)
	result := &TimetableExportFormat{
		StationIdToTrainIdSet: map[model.StationId]model.TrainIdSet{},
		TrainIdToStationMap:   map[model.TrainId]model.StationIdToStationMap{},
		StationIdMap:          map[model.StationId]string{},
	}
	err = dec.Decode(result)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "can not decode timetable with enc.Decode")
	}
	return result.StationIdToTrainIdSet, result.TrainIdToStationMap, result.StationIdMap, nil
}
