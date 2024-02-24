package transfer

import (
	"encoding/gob"
	"io"
	"os"

	"github.com/pkg/errors"

	"zpcg/internal/model"
)

func ExportTimetable(filename string, timetable model.TimetableTransferFormat) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
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
	result, err := ImportTimetableFromReader(file)
	if err != nil {
		return nil, errors.Wrap(err, "ImportTimetableFromReader")
	}
	return result, nil
}

func ImportTimetableFromReader(reader io.Reader) (*model.TimetableTransferFormat, error) {
	dec := gob.NewDecoder(reader)
	result := &model.TimetableTransferFormat{}
	err := dec.Decode(result)
	if err != nil {
		return nil, errors.Wrap(err, "can not decode timetable with enc.Decode")
	}
	return result, nil
}
