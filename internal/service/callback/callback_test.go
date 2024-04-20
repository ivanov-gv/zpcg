package callback

import (
	"testing"
)

func TestConstants(t *testing.T) {
	service := NewCallbackService()
	service.ParseCallback("1|Podgorica | Nikšić|")
}
