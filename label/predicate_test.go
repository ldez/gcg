package label

import (
	"testing"

	"github.com/google/go-github/v71/github"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.True(t, All(ghLabel))
}

func TestNothing(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.False(t, Nothing(ghLabel))
}

func TestNot(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.False(t, Not(All)(ghLabel))
	assert.True(t, Not(Nothing)(ghLabel))
}

func TestHasPrefix(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.True(t, HasPrefix("courg")(ghLabel))
	assert.False(t, HasPrefix("gette")(ghLabel))
}

func TestHasSuffix(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.True(t, HasSuffix("gette")(ghLabel))
	assert.False(t, HasSuffix("courg")(ghLabel))
}

func TestAllMatch(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.True(t, AllMatch(All, All)(ghLabel))
	assert.False(t, AllMatch(All, Nothing)(ghLabel))
	assert.False(t, AllMatch(Nothing, Nothing)(ghLabel))
}

func TestAnyMatch(t *testing.T) {
	ghLabel := &github.Label{Name: github.Ptr("courgette")}

	assert.True(t, AnyMatch(All, All)(ghLabel))
	assert.True(t, AnyMatch(All, Nothing)(ghLabel))
	assert.False(t, AnyMatch(Nothing, Nothing)(ghLabel))
}
