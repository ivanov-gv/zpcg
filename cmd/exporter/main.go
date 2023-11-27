package main

import (
	"log"

	"zpcg/internal/parser"
	"zpcg/internal/transfer"
)

func main() {
	log.Print("starting to parse timetable...  ")
	detailedTimetable, err := parser.ParseTimetable()
	if err != nil {
		log.Fatalf("can not parse timetable: %s", err.Error())
	}
	log.Println("success")

	log.Print("starting to export timetable...  ")
	err = transfer.ExportTimetable("timetable", detailedTimetable)
	if err != nil {
		log.Fatalf("can not parse timetable: %s", err.Error())
	}
	log.Println("success")
}
