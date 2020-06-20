package core

import (
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/ldez/gcg/types"
	"github.com/stretchr/testify/assert"
)

func TestLabelFilter(t *testing.T) {
	prefix := "type"
	suffix := "foo"

	options := &types.DisplayLabelOptions{
		FilteredPrefixes: []string{prefix},
		FilteredSuffixes: []string{suffix},
	}

	carotte := "type/carotte/fii"
	ghLabelCarotte := &github.Label{
		Name: &carotte,
	}

	assert.False(t, labelFilter(options)(ghLabelCarotte))
}
