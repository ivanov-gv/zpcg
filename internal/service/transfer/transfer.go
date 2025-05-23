package transfer

import (
	"fmt"
	"os"

	"github.com/ivanov-gv/zpcg/gen/timetable"
	timetable2 "github.com/ivanov-gv/zpcg/internal/model/timetable"
)

const timetableGoFileFormat = `
// Package timetable defines consts for tg-server bot
// Code generated by cmd/exporter. DO NOT EDIT.
package timetable

import (
	"time"

	"github.com/ivanov-gv/zpcg/internal/model/timetable"
)

var Timetable = %#v
`

func ExportTimetable(filename string, timetable timetable2.ExportFormat) error {
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

func ImportTimetable() timetable2.ExportFormat {
	return timetable.Timetable
}
