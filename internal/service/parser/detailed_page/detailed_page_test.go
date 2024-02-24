package detailed_page

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/test/resources"
)

func TestDetailedTimetablePageParser(t *testing.T) {
	f, err := resources.TestFS.Open(resources.DetailedTimetableFilepath)
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseDetailedTimetablePage(7108, "url", f)
	assert.NoError(t, err, "ParseDetailedTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
