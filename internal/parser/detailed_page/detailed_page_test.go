package detailed_page

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDetailedTimetablePageParser(t *testing.T) {
	f, err := os.Open("../../resources/Željeznički prevoz Crne Gore - Broj voza 7108.html")
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseDetailedTimetablePage(7108, f)
	assert.NoError(t, err, "ParseDetailedTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
