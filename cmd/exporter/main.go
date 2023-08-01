package main

import (
	"log"
	"zpcg/internal/parser"
)

func main() {
	log.Print("starting to parse timetable...  ")
	timetable, err := parser.ParseTimetable()
	if err != nil {
		log.Fatalf("can not parse timetable: %s", err.Error())
	}
	log.Println("success")

	log.Print("starting to export timetable...  ")
	err = parser.ExportTimetable(timetable, "timetable")
	if err != nil {
		log.Fatalf("can not parse timetable: %s", err.Error())
	}
	log.Println("success")
}
