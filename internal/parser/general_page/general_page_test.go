package general_page

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGeneralTimetablePageParser(t *testing.T) {
	f, err := os.Open("../../resources/Željeznički prevoz Crne Gore - Polasci.html")
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseGeneralTimetablePage(f)
	assert.NoError(t, err, "ParseGeneralTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
