package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ivanov-gv/zpcg/internal/service/parser"
	"github.com/ivanov-gv/zpcg/internal/service/transfer"
)

const timetableDefaultFilepath = "timetable.gen.go"

const Date = "2024-12-15"

// additionalRoutesUrls GET response on RoutesApiUrl skips some routes. In order to have complete information we need to make some additional requests
var additionalRoutesUrls = []string{
	"/routes?start=Bar&finish=Novi+Sad&date=" + Date,
	"/routes?start=Novi+Sad&finish=Bar&date=" + Date,
	"/routes?start=Bar&finish=Zemun&date=" + Date,
	"/routes?start=Zemun&finish=Bar&date=" + Date,
	"/routes?start=Podgorica&finish=Bar&date=" + Date,
	"/routes?start=Bar&finish=Podgorica&date=" + Date,
	"/routes?start=Podgorica&finish=Nikšić&date=" + Date,
	"/routes?start=Nikšić&finish=Podgorica&date=" + Date,
	"/routes?start=Bar&finish=Bijelo+Polje&date=" + Date,
	"/routes?start=Bijelo+Polje&finish=Bar&date=" + Date,
}

func main() {
	timetableFilepath := flag.String("file", timetableDefaultFilepath, "filepath to export timetable to")
	flag.Parse()

	log.Print("starting to parse timetable...  ")
	detailedTimetable, err := parser.ParseTimetable(additionalRoutesUrls...)
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
