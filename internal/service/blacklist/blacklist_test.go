package blacklist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	for _, station := range BlackListedStations {
		t.Run(station.Name, func(t *testing.T) {
			// name is not empty and the map is not nil
			assert.NotEmptyf(t, station.Name, "name is not empty")
			assert.NotNilf(t, station.LanguageTagToCustomErrorMessageMap, "the map is not nil")
		})
	}
}
