package transfer

import (
	"fmt"
	"os"

	"zpcg/gen/timetable"
	"zpcg/internal/model"
)

const timetableGoFileFormat = `
package timetable

import (
	"time"

	"zpcg/internal/model"
)

var Timetable = %#v
`

func ExportTimetable(filename string, timetable model.TimetableTransferFormat) error {
	fileContent := fmt.Sprintf(timetableGoFileFormat, timetable)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("can not open file with os.OpenFile: %w", err)
	}
	_, err = file.WriteString(fileContent)
	if err != nil {
		return fmt.Errorf("can not encode timetable with enc.Encode: %w", err)
	}
	return nil
}

func ImportTimetable() model.TimetableTransferFormat {
	return timetable.Timetable
}
