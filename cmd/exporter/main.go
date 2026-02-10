package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ivanov-gv/zpcg/internal/service/parser"
	"github.com/ivanov-gv/zpcg/internal/service/transfer"
)

// how to:
// 1. adjust routesUrls slice to include all routes you want to export
// 2. adjust date to export timetable for
// 3. adjust aliases and the blacklist in 'internal/model/stations' package, if needed.
// 4. rebuild the timetable
// 5. run tests, add or exclude cases with '// summer period stations'
// 6. deploy to preprod
// 7. check everything
// 8. deploy to prod

const timetableDefaultFilepath = "timetable.gen.go"

const Date = "2025-12-14"

// routesUrls GET response on RoutesApiUrl skips some routes. To have complete information, we need to make some additional requests
var routesUrls = []string{
	"/routes?start=Bar&finish=Subotica&date=" + Date,
	"/routes?start=Subotica&finish=Bar&date=" + Date,
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
	detailedTimetable, err := parser.ParseTimetable(routesUrls...)
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
