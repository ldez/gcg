package core

import (
	"testing"

	"github.com/google/go-github/v78/github"
	"github.com/ldez/gcg/types"
	"github.com/stretchr/testify/assert"
)

func Test_labelFilter(t *testing.T) {
	prefix := "type"
	suffix := "foo"

	options := &types.DisplayLabelOptions{
		FilteredPrefixes: []string{prefix},
		FilteredSuffixes: []string{suffix},
	}

	ghLabelCarotte := &github.Label{Name: github.Ptr("type/carotte/fii")}

	assert.False(t, labelFilter(options)(ghLabelCarotte))
}
