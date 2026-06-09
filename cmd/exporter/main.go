package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/ivanov-gv/zpcg/internal/config/timetable_parser_config"
	"github.com/ivanov-gv/zpcg/internal/service/timetable_export"
	"github.com/ivanov-gv/zpcg/internal/service/timetable_parser"
)

// how to:
// 1. adjust routesUrls slice to include all routes you want to export
// 2. adjust date to export timetable for
// 3. adjust aliases and the station_blacklist in 'internal/model/stations' package, if needed.
// 4. rebuild the timetable
// 5. run tests, add or exclude cases with '// summer period stations'
// 6. deploy to preprod
// 7. check everything
// 8. deploy to prod

const (
	timetableDefaultFilepath = "timetable.gen.go"
	configDefaultFilepath    = "timetable-parser-config.yml"
)

func main() {
	configFilepath := flag.String("config", configDefaultFilepath, "filepath to a config file")
	timetableFilepath := flag.String("file", timetableDefaultFilepath, "filepath to export timetable to")
	flag.Parse()

	log.Print("loading config...")
	config, err := timetable_parser_config.Load(*configFilepath)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("can not load config: %w", err)).Send()
	}

	log.Print("parsing timetable...  ")
	parser := timetable_parser.New(config)
	timetable, err := parser.ParseTimetable()
	if err != nil {
		log.Fatal().Err(fmt.Errorf("can not parse timetable: %w", err)).Send()

	}

	log.Print("exporting timetable...  ")
	err = timetable_export.ExportTimetable(*timetableFilepath, timetable)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("can not export timetable: %w", err)).Send()
	}
	log.Print("success")

	//detailedTimetable, err := timetable_parser.ParseTimetable(routesUrls...)
	//if err != nil {
	//	log.Fatal().Err(fmt.Errorf("can not parse timetable: %w", err)).Send()
	//}
	//log.Trace().Msg("success")
	//
	//log.Print("starting to export timetable...  ")
	//err = timetable_export.ExportTimetable(*timetableFilepath, detailedTimetable)
	//if err != nil {
	//	log.Fatal().Err(fmt.Errorf("can not parse timetable: %w", err)).Send()
	//}
	//log.Trace().Msg("success")
}
