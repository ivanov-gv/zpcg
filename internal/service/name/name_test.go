package name

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivanov-gv/zpcg/gen/timetable"
)

func TestUnify(t *testing.T) {
	assert.Equal(t, "niksic", Unify("Nikšić"))
	assert.Equal(t, "baresumanovica", Unify("Bare Šumanovića"))
	assert.Equal(t, "вирпазар", Unify("В и Р п а З а Р"))
}

func TestFindBestMatch(t *testing.T) {
	NewStationNameResolver(nil, timetable.Timetable.UnifiedStationNameList)
	// niksic
	match, err := findBestMatch([]rune(Unify("Nikschichsss   ")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "niksic", match)
	// susan
	match, err = findBestMatch([]rune(Unify("shushan")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "susanj", match)
	// никшич
	match, err = findBestMatch([]rune(Unify("никшич")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "никшич", match)
	// подгорица
	match, err = findBestMatch([]rune(Unify("подгорица")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "подгорица", match)
	// сутоморе
	match, err = findBestMatch([]rune(Unify("сутоморе")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "сутоморе", match)
	// бар
	match, err = findBestMatch([]rune(Unify("бар")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "бар", match)
	// шушань
	match, err = findBestMatch([]rune(Unify("шушань")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "шушань", match)
	// вирпазар
	match, err = findBestMatch([]rune(Unify("вирпазар")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "вирпазар", match)
	// аэропорт
	match, err = findBestMatch([]rune(Unify("аэропорт")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "аэродром", match)
	// мойковац
	match, err = findBestMatch([]rune(Unify("мойковац")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "моиковац", match)
	// бело поле
	match, err = findBestMatch([]rune(Unify("бело поле")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "белополе", match)
	// даниловград
	match, err = findBestMatch([]rune(Unify("даниловград")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "даниловград", match)
	// белград
	match, err = findBestMatch([]rune(Unify("белград")), timetable.Timetable.UnifiedStationNameList)
	assert.NoError(t, err)
	assert.Equal(t, "белградцентр", match)
}
