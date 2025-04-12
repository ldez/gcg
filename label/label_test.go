package label

import (
	"testing"

	"github.com/google/go-github/v71/github"
	"github.com/stretchr/testify/assert"
)

func TestFilterAndTransform(t *testing.T) {
	carotte := "carotte"
	ghLabelCarotte := &github.Label{
		Name: github.Ptr(carotte),
	}

	courgette := "courgette"
	ghLabelCourgette := &github.Label{
		Name: github.Ptr(courgette),
	}

	tomate := "tomate"
	ghLabelTomate := &github.Label{
		Name: github.Ptr(tomate),
	}

	legumes := []*github.Label{ghLabelCarotte, ghLabelCourgette, ghLabelTomate}

	names := FilterAndTransform(legumes, All, NameIdentity)

	expected := []string{carotte, courgette, tomate}

	assert.Equal(t, expected, names)
}

func TestFilterByPrefixAndTransform(t *testing.T) {
	carotte := "type/carotte"

	ghLabelCarotte := &github.Label{
		Name: github.Ptr(carotte),
	}

	ghLabelCourgette := &github.Label{
		Name: github.Ptr("courgette"),
	}

	ghLabelTomate := &github.Label{
		Name: github.Ptr("tomate"),
	}

	legumes := []*github.Label{ghLabelCarotte, ghLabelCourgette, ghLabelTomate}

	names := FilterAndTransform(legumes, HasPrefix("type/"), NameIdentity)

	expected := []string{carotte}

	assert.Equal(t, expected, names)
}

func TestFilterAndTransformName(t *testing.T) {
	ghLabelCarotte := &github.Label{
		Name: github.Ptr("type/carotte"),
	}

	ghLabelCourgette := &github.Label{
		Name: github.Ptr("type/courgette"),
	}

	ghLabelTomate := &github.Label{
		Name: github.Ptr("type/tomate"),
	}

	legumes := []*github.Label{ghLabelCarotte, ghLabelCourgette, ghLabelTomate}

	names := FilterAndTransform(legumes, All, TrimPrefix("type/"))

	expected := []string{"carotte", "courgette", "tomate"}

	assert.Equal(t, expected, names)
}

func TestFilterByPrefixAndTransformName(t *testing.T) {
	ghLabelCarotte := &github.Label{
		Name: github.Ptr("type/carotte"),
	}

	ghLabelCourgette := &github.Label{
		Name: github.Ptr("type/courgette"),
	}

	ghLabelTomate := &github.Label{
		Name: github.Ptr("tomate"),
	}

	legumes := []*github.Label{ghLabelCarotte, ghLabelCourgette, ghLabelTomate}

	names := FilterAndTransform(legumes, HasPrefix("type/"), TrimPrefix("type/"))

	expected := []string{"carotte", "courgette"}

	assert.Equal(t, expected, names)
}

func TestFilteredByExistingSuffix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"otte", "cour"}

	assert.True(t, FilteredBy(HasSuffix, prefixes)(ghLabel))
}

func TestFilteredByNonExistingSuffix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"to", "cour"}

	assert.False(t, FilteredBy(HasSuffix, prefixes)(ghLabel))
}

func TestFilteredByExistingPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"car", "cour"}

	assert.True(t, FilteredBy(HasPrefix, prefixes)(ghLabel))
}

func TestFilteredByNonExistingPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"to", "cour"}

	assert.False(t, FilteredBy(HasPrefix, prefixes)(ghLabel))
}

func TestFilteredByEmptyPrefixesList(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	var prefixes []string

	assert.True(t, FilteredBy(HasPrefix, prefixes)(ghLabel))
}

func TestFilteredByEmptyPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{""}

	assert.True(t, FilteredBy(HasPrefix, prefixes)(ghLabel))
}

func TestExcludedByExistingPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"car", "cour"}

	assert.False(t, ExcludedBy(HasPrefix, prefixes)(ghLabel))
}

func TestExcludedByNonExistingPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{"to", "cour"}

	assert.True(t, ExcludedBy(HasPrefix, prefixes)(ghLabel))
}

func TestExcludedByEmptyPrefixesList(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	var prefixes []string

	assert.True(t, ExcludedBy(HasPrefix, prefixes)(ghLabel))
}

func TestExcludedByEmptyPrefix(t *testing.T) {
	ghLabel := &github.Label{
		Name: github.Ptr("carotte"),
	}

	prefixes := []string{""}

	assert.True(t, ExcludedBy(HasPrefix, prefixes)(ghLabel))
}

func TestTrimAllPrefixExistingPrefix(t *testing.T) {
	prefixes := []string{"type/", "legume/"}

	assert.Equal(t, "potiron", TrimAllPrefix(prefixes)("type/legume/potiron"))
}

func TestTrimAllPrefixNonExistingPrefix(t *testing.T) {
	prefixes := []string{"value/", "fruit/"}

	assert.Equal(t, "type/legume/potiron", TrimAllPrefix(prefixes)("type/legume/potiron"))
}

func TestTrimAllPrefixEmptyPrefix(t *testing.T) {
	var prefixes []string

	assert.Equal(t, "type/legume/potiron", TrimAllPrefix(prefixes)("type/legume/potiron"))
}
