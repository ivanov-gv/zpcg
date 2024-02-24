package general_page

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"zpcg/test/resources"
)

func TestGeneralTimetablePageParser(t *testing.T) {
	f, err := resources.TestFS.Open(resources.GeneralTimetableHtmlFilepath)
	assert.NoError(t, err, "os.Open")
	parsedLinks, err := ParseGeneralTimetablePage(f)
	assert.NoError(t, err, "ParseGeneralTimetablePage")
	assert.NotEmpty(t, parsedLinks)
	t.Log(parsedLinks)
}
