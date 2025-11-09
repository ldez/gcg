package types

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v78/github"
)

// NoOption empty struct.
type NoOption struct{}

// Configuration GCG Configuration.
type Configuration struct {
	ConfigFile           string               `long:"config-file" description:"A configuration file. [optional]" toml:"-"`
	Owner                string               `short:"o" description:"Repository owner."`
	RepositoryName       string               `long:"repo-name" short:"r" description:"Repository name."`
	GitHubToken          string               `long:"token" short:"t" description:"GitHub Token. [optional]"`
	OutputType           string               `long:"output-type" description:"Output destination type. (file|Stdout)"`
	FileName             string               `long:"file-name" description:"Name of the changelog file."`
	CurrentRef           string               `long:"current-ref" short:"c" description:"Current commit reference. Can be a tag, a branch, a SHA."`
	PreviousRef          string               `long:"previous-ref" short:"p" description:"Previous commit reference. Can be a tag, a branch, a SHA."`
	BaseBranch           string               `long:"base-branch" short:"b" description:"Base branch name. PR branch destination."`
	FutureCurrentRefName string               `long:"future-ref-name" short:"f" description:"The future name of the current reference."`
	ThresholdPreviousRef int                  `long:"th-before" description:"Threshold in seconds after the previous ref date."`
	ThresholdCurrentRef  int                  `long:"th-after" description:"Threshold in seconds after the current ref date."`
	Debug                bool                 `long:"debug" description:"Debug mode."`
	DisplayLabel         bool                 `long:"display-label" description:"Display labels"`
	LabelExcludes        []string             `long:"exclude-label" description:"Label to exclude."`
	LabelEnhancement     string               `long:"enhancement-label" description:"Enhancement Label."`
	LabelDocumentation   string               `long:"doc-label" description:"Documentation Label."`
	LabelBug             string               `long:"bug-label" description:"Bug Label."`
	DisplayLabelOptions  *DisplayLabelOptions `long:"dl-options" description:"Label display options."`
	TemplateFile         string               `long:"tmpl-file" description:"A template file. [optional]"`
}

// DisplayLabelOptions the options defining the labeling display.
type DisplayLabelOptions struct {
	FilteredPrefixes []string `long:"prefix-filter" description:"Included label prefixes."`
	ExcludedPrefixes []string `long:"prefix-exclude" description:"Excluded label prefixes."`
	FilteredSuffixes []string `long:"suffix-filter" description:"Included label suffixes."`
	ExcludedSuffixes []string `long:"suffix-exclude" description:"Excluded label suffixes."`
	TrimmedPrefixes  []string `long:"prefix-trim" description:"Trim label with the following prefixes."`
}

// Summary a repository summary.
type Summary struct {
	CurrentRefName  string
	CurrentRefDate  string
	PreviousRefName string
	Owner           string
	RepositoryName  string
	Enhancement     []IssueSummary
	Documentation   []IssueSummary
	Bug             []IssueSummary
	Other           []IssueSummary
}

// IssueSummary an issue summary.
type IssueSummary struct {
	FilteredLabelNames string
	Issue              *github.Issue
}

// LabelDisplayOptionsParser a parser for DisplayLabelOptions.
type LabelDisplayOptionsParser DisplayLabelOptions

// Set a DisplayLabelOptions.
func (c *LabelDisplayOptionsParser) Set(s string) error {
	log.Println("sure?:", s)

	return nil
}

// Get a DisplayLabelOptions.
func (c *LabelDisplayOptionsParser) Get() any { return DisplayLabelOptions(*c) }

// String a string representation of DisplayLabelOptions.
func (c *LabelDisplayOptionsParser) String() string { return fmt.Sprintf("%v", *c) }

// SetValue a DisplayLabelOptions.
func (c *LabelDisplayOptionsParser) SetValue(val any) {
	*c = LabelDisplayOptionsParser(val.(DisplayLabelOptions)) //nolint:forcetypeassert // not useful here.
}

// SliceString type used for flaeg parsing.
type SliceString []string

// Set a SliceString.
func (c *SliceString) Set(rawValue string) error {
	values := strings.Split(rawValue, ",")

	if len(values) == 0 {
		return fmt.Errorf("bad Value format: %s", rawValue)
	}

	for _, value := range values {
		*c = append(*c, value)
	}

	return nil
}

// Get a SliceString.
func (c *SliceString) Get() any { return []string(*c) }

// String a string representation of SliceString.
func (c *SliceString) String() string {
	return strings.Join(*c, ",")
}

// SetValue a SliceString.
func (c *SliceString) SetValue(val any) {
	*c = SliceString(val.([]string)) //nolint:forcetypeassert // not useful here.
}

// ByLabel sort by label.
type ByLabel []IssueSummary

func (a ByLabel) Len() int      { return len(a) }
func (a ByLabel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLabel) Less(i, j int) bool {
	if a[i].FilteredLabelNames == "" && a[j].FilteredLabelNames != "" {
		return false
	}

	if a[j].FilteredLabelNames == "" && a[i].FilteredLabelNames != "" {
		return true
	}

	return a[i].FilteredLabelNames < a[j].FilteredLabelNames
}
