package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ivanov-gv/zpcg/gen/timetable"
)

func TestUnify(t *testing.T) {
	assert.Equal(t, "niksic", Unify("Nikšić"))
	assert.Equal(t, "baresumanovica", Unify("Bare Šumanovića"))
	assert.Equal(t, "вирпазар", Unify("В и Р п а З а Р"))
}

func TestFindStationByApproximateName(t *testing.T) {
	resolver := NewStationNameResolver(
		timetable.Timetable.UnifiedStationNameToStationIdMap,
		timetable.Timetable.UnifiedStationNameList,
		timetable.Timetable.StationIdToStationMap,
	)

	tests := []struct {
		input        string
		expectedName string
	}{
		// typos and transliterations
		{"Nikschichsss   ", "Nikšić"},
		{"shushan", "Šušanj"},
		{"padgareeka", "Podgorica"},
		// Cyrillic inputs
		{"никшич", "Nikšić"},
		{"подгорица", "Podgorica"},
		{"сутоморе", "Sutomore"},
		{"бар", "Bar"},
		{"шушань", "Šušanj"},
		{"вирпазар", "Virpazar"},
		{"мойковац", "Mojkovac"},
		{"даниловград", "Danilovgrad"},
		{"белград", "Beograd Centar"},
		// approximate Cyrillic
		{"бело поле", "Bijelo Polje"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, stationName, err := resolver.FindStationIdByApproximateName(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, stationName)
		})
	}
}

func TestFindStationByApproximateName_NoMatch(t *testing.T) {
	resolver := NewStationNameResolver(
		timetable.Timetable.UnifiedStationNameToStationIdMap,
		timetable.Timetable.UnifiedStationNameList,
		timetable.Timetable.StationIdToStationMap,
	)

	noMatchInputs := []string{"London", "Copenhagen", "Berlin"}

	for _, input := range noMatchInputs {
		t.Run(input, func(t *testing.T) {
			_, _, err := resolver.FindStationIdByApproximateName(input)
			assert.ErrorIs(t, err, ErrNoMatchesFound)
		})
	}
}
