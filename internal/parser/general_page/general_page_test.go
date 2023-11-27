package general_page

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const generalTimetableHtmlFilepath = "../../../resources/Željeznički prevoz Crne Gore - Polasci.html"

func TestGeneralTimetablePageParser(t *testing.T) {
	f, err := os.Open(generalTimetableHtmlFilepath)
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseGeneralTimetablePage(f)
	assert.NoError(t, err, "ParseGeneralTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
