package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"zpcg/internal/parser"
	"zpcg/internal/transfer"
)

const timetableDefaultFilepath = "timetable.gob"

func main() {
	timetableFilepath := flag.String("file", timetableDefaultFilepath, "filepath to export timetable to")
	flag.Parse()

	log.Print("starting to parse timetable...  ")
	detailedTimetable, err := parser.ParseTimetable()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("can not parse timetable: %w", err)).Send()
	}
	log.Trace().Msg("success")

	log.Print("starting to export timetable...  ")
	err = transfer.ExportTimetable(*timetableFilepath, detailedTimetable)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("can not parse timetable: %w", err)).Send()
	}
	log.Trace().Msg("success")
}
