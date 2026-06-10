package pathfinder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivanov-gv/zpcg/internal/pkg/timetable_export"
	"github.com/ivanov-gv/zpcg/internal/service/date"
	"github.com/ivanov-gv/zpcg/internal/service/station_name_resolver"
	"github.com/ivanov-gv/zpcg/internal/service/timetable"
)

const (
	NiksicStationName      = "Nikšić"
	DanilovgradStationName = "Danilovgrad"
	BarStationName         = "Bar"
)

func TestFindDirectPaths(t *testing.T) {
	_timetable := timetable_export.ImportTimetable()

	for _, season := range _timetable.Seasons {
		t.Run(season.Name, func(t *testing.T) {
			dateService := date.NewDateService(t.Context(), date.FixedDate(season.Start))
			timetableService, err := timetable.New(_timetable, dateService)
			assert.NoError(t, err)

			pathFinder := NewPathFinder(timetableService)
			paths := pathFinder.findDirectPaths(season,
				_timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(NiksicStationName)],
				_timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(DanilovgradStationName)])
			assert.NotNil(t, paths)
			assert.NotEmpty(t, paths)
		})
	}
}

func TestFindPaths(t *testing.T) {
	_timetable := timetable_export.ImportTimetable()

	for _, season := range _timetable.Seasons {
		t.Run(season.Name, func(t *testing.T) {
			dateService := date.NewDateService(t.Context(), date.FixedDate(season.Start))
			timetableService, err := timetable.New(_timetable, dateService)
			assert.NoError(t, err)

			pathFinder := NewPathFinder(timetableService)
			paths := pathFinder.findPathsWithTransfer(season,
				_timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(NiksicStationName)],
				_timetable.UnifiedStationNameToStationIdMap[station_name_resolver.Unify(BarStationName)])
			assert.NotNil(t, paths)
			assert.NotEmpty(t, paths)
		})
	}
}
