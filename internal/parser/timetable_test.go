package parser

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"zpcg/internal/parser/detailed_page"
	"zpcg/internal/parser/general_page"
)

func TestGeneralTimetablePageParser(t *testing.T) {
	f, err := os.Open("../../resources/Željeznički prevoz Crne Gore - Polasci.html")
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := general_page.ParseGeneralTimetablePage(f)
	assert.NoError(t, err, "ParseGeneralTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}

func TestDetailedTimetablePageParser(t *testing.T) {
	f, err := os.Open("../../resources/Željeznički prevoz Crne Gore - Broj voza 7108.html")
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := detailed_page.ParseDetailedTimetablePage(0, f)
	assert.NoError(t, err, "ParseDetailedTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
