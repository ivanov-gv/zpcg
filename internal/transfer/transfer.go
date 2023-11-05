package transfer

import (
	"encoding/gob"
	"github.com/pkg/errors"
	"os"
	"zpcg/internal/model"
)

func ExportTimetable(filename string, timetable model.TimetableTransferFormat) error {
	file, err := os.OpenFile(filename+".gob", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return errors.Wrap(err, "can not open file with os.OpenFile")
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(timetable)
	if err != nil {
		return errors.Wrap(err, "can not encode timetable with enc.Encode")
	}
	return nil
}

func ImportTimetable(filename string) (*model.TimetableTransferFormat, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, errors.Wrap(err, "can not open file with os.Open")
	}
	dec := gob.NewDecoder(file)
	result := &model.TimetableTransferFormat{
		//StationIdToTrainIdSet:            map[model.StationId]model.TrainIdSet{},
		//TrainIdToStationMap:              map[model.TrainId]model.StationIdToStationMap{},
		//StationIdToStaionMap:                     map[model.StationId]model.Station{},
		//TrainIdToTrainInfoMap:            map[model.TrainId]model.TrainInfo{},
		//UnifiedStationNameToStationIdMap: map[string]model.StationId{},
		//UnifiedStationNameList:           []string{},
	}
	err = dec.Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "can not decode timetable with enc.Decode")
	}
	return result, nil
}
