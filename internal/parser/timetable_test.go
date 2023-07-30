package parser

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"zpcg/internal/parser/general_page"
)

func TestTimetableParser(t *testing.T) {
	f, err := os.Open("../../resources/Željeznički prevoz Crne Gore - Polasci.html")
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := general_page.ParseGeneralTimetablePage(f)
	assert.NoError(t, err, "ParseTimetable")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
