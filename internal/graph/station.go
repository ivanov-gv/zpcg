package graph

import "zpcg/internal/parser/model"

const MaxNumberOfStations = 75

func GetStationsFromRoutes(routes map[int]model.DetailedTimetable) (stationNameToIdMap map[string]int, numOfStations int) {
	stationNameToIdMap = make(map[string]int, MaxNumberOfStations)

	var i int
	for _, timetable := range routes {
		for _, station := range timetable.Stations {
			if _, found := stationNameToIdMap[station.Name]; found {
				continue
			}
			// station is not added yet
			stationNameToIdMap[station.Name] = i
			i++
		}
	}
	return stationNameToIdMap, len(stationNameToIdMap)
}
