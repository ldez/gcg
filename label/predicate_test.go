package label

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.True(t, All(ghLabel))
}

func TestNothing(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.False(t, Nothing(ghLabel))
}

func TestNot(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.False(t, Not(All)(ghLabel))
	assert.True(t, Not(Nothing)(ghLabel))
}

func TestHasPrefix(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.True(t, HasPrefix("courg")(ghLabel))
	assert.False(t, HasPrefix("gette")(ghLabel))
}

func TestHasSuffix(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.True(t, HasSuffix("gette")(ghLabel))
	assert.False(t, HasSuffix("courg")(ghLabel))
}

func TestAllMatch(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.True(t, AllMatch(All, All)(ghLabel))
	assert.False(t, AllMatch(All, Nothing)(ghLabel))
	assert.False(t, AllMatch(Nothing, Nothing)(ghLabel))
}

func TestAnyMatch(t *testing.T) {
	name := "courgette"
	ghLabel := github.Label{
		Name: &name,
	}
	assert.True(t, AnyMatch(All, All)(ghLabel))
	assert.True(t, AnyMatch(All, Nothing)(ghLabel))
	assert.False(t, AnyMatch(Nothing, Nothing)(ghLabel))
}
