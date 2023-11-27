package detailed_page

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const detailedTimetableFilepath = "../../../resources/Željeznički prevoz Crne Gore - Broj voza 7108.html"

func TestDetailedTimetablePageParser(t *testing.T) {
	f, err := os.Open(detailedTimetableFilepath)
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseDetailedTimetablePage(7108, "url", f)
	assert.NoError(t, err, "ParseDetailedTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
