package label

import (
	"testing"

	"github.com/google/go-github/v71/github"
	"github.com/stretchr/testify/assert"
)

func TestFilterAndTransform(t *testing.T) {
	testCases := []struct {
		desc      string
		labels    []*github.Label
		filter    Predicate
		transform NameTransform
		expected  []string
	}{
		{
			desc: "identity",
			labels: []*github.Label{
				{Name: github.Ptr("carotte")},
				{Name: github.Ptr("courgette")},
				{Name: github.Ptr("tomate")},
			},
			filter:    All,
			transform: NameIdentity,
			expected:  []string{"carotte", "courgette", "tomate"},
		},
		{
			desc: "filter by prefix",
			labels: []*github.Label{
				{Name: github.Ptr("type/carotte")},
				{Name: github.Ptr("courgette")},
				{Name: github.Ptr("tomate")},
			},
			filter:    HasPrefix("type/"),
			transform: NameIdentity,
			expected:  []string{"type/carotte"},
		},
		{
			desc: "transform names",
			labels: []*github.Label{
				{Name: github.Ptr("type/carotte")},
				{Name: github.Ptr("type/courgette")},
				{Name: github.Ptr("type/tomate")},
			},
			filter:    All,
			transform: TrimPrefix("type/"),
			expected:  []string{"carotte", "courgette", "tomate"},
		},
		{
			desc: "filter by prefix and transform names",
			labels: []*github.Label{
				{Name: github.Ptr("type/carotte")},
				{Name: github.Ptr("type/courgette")},
				{Name: github.Ptr("tomate")},
			},
			filter:    HasPrefix("type/"),
			transform: TrimPrefix("type/"),
			expected:  []string{"carotte", "courgette"},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			names := FilterAndTransform(test.labels, test.filter, test.transform)

			assert.Equal(t, test.expected, names)
		})
	}
}

func TestFilteredBy(t *testing.T) {
	testCases := []struct {
		desc          string
		label         *github.Label
		namePredicate func(string) Predicate
		values        []string
		assert        assert.BoolAssertionFunc
	}{
		{
			desc:          "by existing suffix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasSuffix,
			values:        []string{"otte", "cour"},
			assert:        assert.True,
		},
		{
			desc:          "by non-existing suffix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasSuffix,
			values:        []string{"to", "cour"},
			assert:        assert.False,
		},
		{
			desc:          "by existing prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{"car", "cour"},
			assert:        assert.True,
		},
		{
			desc:          "by non-existing prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{"to", "cour"},
			assert:        assert.False,
		},
		{
			desc:          "by empty prefixes list",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        nil,
			assert:        assert.True,
		},
		{
			desc:          "by empty prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{""},
			assert:        assert.True,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, FilteredBy(test.namePredicate, test.values)(test.label))
		})
	}
}

func TestExcludedBy(t *testing.T) {
	testCases := []struct {
		desc          string
		label         *github.Label
		namePredicate func(string) Predicate
		values        []string
		assert        assert.BoolAssertionFunc
	}{
		{
			desc:          "by existing prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{"car", "cour"},
			assert:        assert.False,
		},
		{
			desc:          "by non-existing prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{"to", "cour"},
			assert:        assert.True,
		},
		{
			desc:          "by empty prefixes list",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        nil,
			assert:        assert.True,
		},
		{
			desc:          "by empty prefix",
			label:         &github.Label{Name: github.Ptr("carotte")},
			namePredicate: HasPrefix,
			values:        []string{""},
			assert:        assert.True,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, ExcludedBy(test.namePredicate, test.values)(test.label))
		})
	}
}

func TestTrimAllPrefix(t *testing.T) {
	testCases := []struct {
		desc     string
		prefixes []string
		name     string
		expected string
	}{
		{
			desc:     "existing prefix",
			prefixes: []string{"type/", "legume/"},
			name:     "type/legume/potiron",
			expected: "potiron",
		},
		{
			desc:     "non-existing prefix",
			prefixes: []string{"value/", "fruit/"},
			name:     "type/legume/potiron",
			expected: "type/legume/potiron",
		},
		{
			desc:     "empty prefix",
			prefixes: nil,
			name:     "type/legume/potiron",
			expected: "type/legume/potiron",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.expected, TrimAllPrefix(test.prefixes)("type/legume/potiron"))
		})
	}
}
