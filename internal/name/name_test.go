package name

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnify(t *testing.T) {
	assert.Equal(t, "niksic", Unify("Nikšić"))
	assert.Equal(t, "baresumanovica", Unify("Bare Šumanovića"))
}
