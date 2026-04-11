package callback

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	service := NewCallbackService()
	_, err := service.ParseCallback("1|Podgorica | Nikšić|")
	assert.NoError(t, err)
}
