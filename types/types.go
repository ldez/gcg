package types

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
)

type Configuration struct {
	CurrentRef           string               `long:"current-ref" short:"c" description:"Current commit reference. Can be a tag, a branch, a SHA."`
	PreviousRef          string               `long:"previous-ref" short:"p" description:"Previous commit reference. Can be a tag, a branch, a SHA."`
	BaseBranch           string               `long:"base-branch" short:"b" description:"Base branch name. PR branch destination."`
	FutureCurrentRefName string               `long:"future-ref-name" short:"f" description:"The future name of the current reference."`
	Owner                string               `short:"o" description:"Repository owner."`
	RepositoryName       string               `long:"repo-name" short:"r" description:"Repository name."`
	GitHubToken          string               `long:"token" short:"t" description:"GitHub Token."`
	LabelExcludes        []string             `long:"exclude-label" description:"Label to exclude."`
	LabelEnhancement     string               `long:"enhancement-label" description:"Enhancement Label."`
	LabelDocumentation   string               `long:"doc-label" description:"Documentation Label."`
	LabelBug             string               `long:"bug-label" description:"Bug Label."`
	DisplayLabel         bool                 `long:"display-label" description:"Display labels"`
	DisplayLabelOptions  *DisplayLabelOptions `long:"dl-options" description:"Label display options."`
	OutputType           string               `long:"output-type" description:"Output destination type. (file|Stdout)"`
	FileName             string               `long:"file-name" description:"Name of the changelog file."`
	Debug                bool                 `long:"debug" description:"Debug mode."`
}

type DisplayLabelOptions struct {
	FilteredPrefixes []string `long:"prefix-filter" description:"Included label prefixes."`
	ExcludedPrefixes []string `long:"prefix-exclude" description:"Excluded label prefixes."`
	FilteredSuffixes []string `long:"suffix-filter" description:"Included label suffixes."`
	ExcludedSuffixes []string `long:"suffix-exclude" description:"Excluded label suffixes."`
	TrimmedPrefixes  []string `long:"prefix-trim" description:"Trim label with the following prefixes."`
}

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

type IssueSummary struct {
	FilteredLabelNames string
	Issue              github.Issue
}

type LabelDisplayOptionsParser DisplayLabelOptions

func (c *LabelDisplayOptionsParser) Set(s string) error {
	log.Println("sure?:", s)
	return nil
}

func (c *LabelDisplayOptionsParser) Get() interface{} { return DisplayLabelOptions(*c) }

func (c *LabelDisplayOptionsParser) String() string { return fmt.Sprintf("%v", *c) }

func (c *LabelDisplayOptionsParser) SetValue(val interface{}) {
	*c = LabelDisplayOptionsParser(val.(DisplayLabelOptions))
}

type SliceString []string

func (c *SliceString) Set(value string) error {
	values := strings.Split(value, ",")
	if len(values) == 0 {
		return errors.New("Bad Value format: " + value)
	}
	for _, val := range values {
		*c = append(*c, val)
	}
	return nil
}

func (c *SliceString) Get() interface{} { return []string(*c) }

func (c *SliceString) String() string {
	return strings.Join(*c, ",")
}

func (c *SliceString) SetValue(val interface{}) {
	*c = SliceString(val.([]string))
}
