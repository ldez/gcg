package label

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameIdentity(t *testing.T) {
	name := "tomate"

	assert.Equal(t, name, NameIdentity(name))
}

func TestTrimPrefix(t *testing.T) {
	name := "tomate"

	assert.Equal(t, "ate", TrimPrefix("tom")(name))
	assert.Equal(t, name, TrimPrefix("ate")(name))
}

func TestChain(t *testing.T) {
	name := "tomate"

	assert.Equal(t, "ate", Chain(NameIdentity, TrimPrefix("tom"))(name))
	assert.Equal(t, name, Chain(NameIdentity, TrimPrefix("ate"))(name))
}
