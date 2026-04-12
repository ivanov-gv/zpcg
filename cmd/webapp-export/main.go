// Package main exports the compiled-in timetable as JSON files for the Telegram Mini App prototype.
//
// It reads gen/timetable.Timetable (the same struct the bot itself uses) and writes:
//
//	<outDir>/stations.json — stations indexed by id with type + cyr/lat names
//	<outDir>/trains.json   — trains with their ordered stop list + times as "HH:MM"
//	<outDir>/meta.json     — small metadata blob (transfer station id, generated-at, version)
//
// Run from repo root:
//
//	go run ./cmd/webapp-export -out=webapp/data
//
// The output is consumed by webapp/js/data.js in the prototype.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/ivanov-gv/zpcg/gen/timetable"
	tt "github.com/ivanov-gv/zpcg/internal/model/timetable"
)

type stationJSON struct {
	Id         int    `json:"id"`
	ZpcgStopId int    `json:"zpcgStopId"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	NameEn     string `json:"nameEn"`
	NameCyr    string `json:"nameCyr"`
}

type stopJSON struct {
	StationId int    `json:"stationId"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
}

type trainJSON struct {
	Id           int        `json:"id"`
	TimetableUrl string     `json:"timetableUrl"`
	Stops        []stopJSON `json:"stops"`
}

type stationTypeJSON struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	NameEn string `json:"nameEn"`
}

type metaJSON struct {
	GeneratedAt       string            `json:"generatedAt"`
	TransferStationId int               `json:"transferStationId"`
	StationTypes      []stationTypeJSON `json:"stationTypes"`
}

func main() {
	outDir := flag.String("out", "webapp/data", "output directory for generated JSON files")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", *outDir, err)
		os.Exit(1)
	}

	t := timetable.Timetable

	// stations
	stations := make([]stationJSON, 0, len(t.StationIdToStationMap))
	for _, s := range t.StationIdToStationMap {
		stations = append(stations, stationJSON{
			Id:         int(s.Id),
			ZpcgStopId: s.ZpcgStopId,
			Type:       int(s.Type),
			Name:       s.Name,
			NameEn:     s.NameEn,
			NameCyr:    s.NameCyr,
		})
	}
	sort.Slice(stations, func(i, j int) bool { return stations[i].Id < stations[j].Id })
	writeJSON(filepath.Join(*outDir, "stations.json"), stations)

	// trains
	trains := make([]trainJSON, 0, len(t.TrainIdToStationMap))
	for trainId, stops := range t.TrainIdToStationMap {
		tj := trainJSON{
			Id:    int(trainId),
			Stops: make([]stopJSON, 0, len(stops)),
		}
		if info, ok := t.TrainIdToTrainInfoMap[trainId]; ok {
			tj.TimetableUrl = info.TimetableUrl
		}
		for _, stop := range stops {
			tj.Stops = append(tj.Stops, stopJSON{
				StationId: int(stop.Id),
				Arrival:   hhmm(stop.Arrival),
				Departure: hhmm(stop.Departure),
			})
		}
		// Sort stops along the train's route by departure time (or arrival if equal).
		// Cross-midnight trains don't exist in this dataset so naive sort is fine.
		sort.Slice(tj.Stops, func(i, j int) bool {
			if tj.Stops[i].Departure == tj.Stops[j].Departure {
				return tj.Stops[i].Arrival < tj.Stops[j].Arrival
			}
			return tj.Stops[i].Departure < tj.Stops[j].Departure
		})
		trains = append(trains, tj)
	}
	sort.Slice(trains, func(i, j int) bool { return trains[i].Id < trains[j].Id })
	writeJSON(filepath.Join(*outDir, "trains.json"), trains)

	// meta
	types := make([]stationTypeJSON, 0, len(t.StationTypes))
	for _, st := range t.StationTypes {
		types = append(types, stationTypeJSON{
			Id:     int(st.Id),
			Name:   st.Name,
			NameEn: st.NameEn,
		})
	}
	sort.Slice(types, func(i, j int) bool { return types[i].Id < types[j].Id })
	meta := metaJSON{
		GeneratedAt:       time.Now().UTC().Format(time.RFC3339),
		TransferStationId: int(t.TransferStationId),
		StationTypes:      types,
	}
	writeJSON(filepath.Join(*outDir, "meta.json"), meta)

	fmt.Printf("exported %d stations, %d trains, %d station types to %s\n",
		len(stations), len(trains), len(types), *outDir)
	_ = tt.LocalTrain // silence unused import if model constants aren't referenced
}

func hhmm(t time.Time) string {
	return t.Format("15:04")
}

func writeJSON(path string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal %s: %v\n", path, err)
		os.Exit(1)
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("  wrote %s (%d bytes)\n", path, len(b))
}
