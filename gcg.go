package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/containous/flaeg"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	DefaultBaseBranch         = "master"
	DefaultEnhancementLabel   = "enhancement"
	DefaultDocumentationLabel = "documentation"
	DefaultBugLabel           = "bug"
	DefaultOutputDestination  = "file"
	DefaultFileName           = "CHANGELOG.md"
	GitHubSearchDateLayout    = "2006-01-02T15:04:05Z"
)

type Configuration struct {
	CurrentRef           string `long:"current-ref" short:"c" description:"Current commit reference. Can be a tag, a branch, a SHA."`
	PreviousRef          string `long:"previous-ref" short:"p" description:"Previous commit reference. Can be a tag, a branch, a SHA."`
	BaseBranch           string `long:"base-branch" short:"b" description:"Base branch name. PR branch destination."`
	FutureCurrentRefName string `long:"future-ref-name" short:"f" description:"The future name of the current reference."`
	Owner                string `short:"o" description:"Repository owner."`
	RepositoryName       string `long:"repo-name" short:"r" description:"Repository name."`
	GitHubToken          string `long:"token" short:"t" description:"GitHub Token"`
	LabelExclude         string `long:"exclude-label" description:"Label to exclude."`
	LabelEnhancement     string `long:"enhancement-label" description:"Enhancement Label."`
	LabelDocumentation   string `long:"doc-label" description:"Documentation Label."`
	LabelBug             string `long:"bug-label" description:"Bug Label."`
	OutputType           string `long:"output-type" description:"Output destination type. (file|Stdout)"`
	FileName             string `long:"file-name" description:"Name of the changelog file."`
	Debug                bool   `long:"debug" description:"Debug mode."`
}

type Summary struct {
	CurrentRefName  string
	CurrentRefDate  string
	PreviousRefName string
	Owner           string
	RepositoryName  string
	Enhancement     []github.Issue
	Documentation   []github.Issue
	Bug             []github.Issue
	Other           []github.Issue
}

func main() {

	config := &Configuration{
		BaseBranch:         DefaultBaseBranch,
		LabelEnhancement:   DefaultEnhancementLabel,
		LabelDocumentation: DefaultDocumentationLabel,
		LabelBug:           DefaultBugLabel,
		OutputType:         DefaultOutputDestination,
		FileName:           DefaultFileName,
	}

	rootCmd := &flaeg.Command{
		Name: "gcg",
		Description: `GCG is a GitHub Changelog Generator.
The generator use only Pull Requests.
		`,
		Config:                config,
		DefaultPointersConfig: &Configuration{},
		Run: func() error {
			if config.Debug {
				log.Printf("Run GCG command with config : %+v\n", config)
			}
			required(config.CurrentRef, "current-ref")
			required(config.PreviousRef, "previous-ref")
			required(config.Owner, "owner")
			required(config.RepositoryName, "repo-name")

			generate(config)
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.Run()
}

func generate(config *Configuration) {
	ctx := context.Background()

	client := newGitHubClient(ctx, config.GitHubToken)

	// Get previous ref date
	commitPreviousRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.PreviousRef)
	check(err)

	datePreviousRef := commitPreviousRef.Commit.Committer.Date.Add(1 * time.Second).Format(GitHubSearchDateLayout)

	// Get current ref version date
	commitCurrentRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.CurrentRef)
	check(err)

	dateCurrentRef := commitCurrentRef.Commit.Committer.Date.Format(GitHubSearchDateLayout)

	// Search PR
	query := fmt.Sprintf("type:pr is:merged repo:%s/%s base:%s merged:%s..%s",
		config.Owner, config.RepositoryName, config.BaseBranch, datePreviousRef, dateCurrentRef)
	if config.Debug {
		log.Println(query)
	}

	searchOptions := &github.SearchOptions{
		Sort:        "created",
		Order:       "asc",
		ListOptions: github.ListOptions{PerPage: 20},
	}

	allSearchResult := searchAllIssues(ctx, client, query, searchOptions, config)
	display(config, allSearchResult, commitCurrentRef)
}

func searchAllIssues(ctx context.Context, client *github.Client, query string, searchOptions *github.SearchOptions, config *Configuration) []github.Issue {
	var allSearchResult []github.Issue
	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			log.Fatal(err)
		}
		for _, issueResult := range issuesSearchResult.Issues {
			if contains(issueResult.Labels, config.LabelExclude) {
				if config.Debug {
					log.Println("Exclude:", *issueResult.Number, *issueResult.Title)
				}
			} else {
				allSearchResult = append(allSearchResult, issueResult)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		searchOptions.Page = resp.NextPage
	}
	return allSearchResult
}

func display(config *Configuration, allSearchResult []github.Issue, commitCurrentRef *github.RepositoryCommit) {
	summary := &Summary{
		Owner:          config.Owner,
		RepositoryName: config.RepositoryName,
	}

	for _, issueResult := range allSearchResult {
		if contains(issueResult.Labels, config.LabelDocumentation) {
			summary.Documentation = append(summary.Documentation, issueResult)
		} else if contains(issueResult.Labels, config.LabelEnhancement) {
			summary.Enhancement = append(summary.Enhancement, issueResult)
		} else if contains(issueResult.Labels, config.LabelBug) {
			summary.Bug = append(summary.Bug, issueResult)
		} else {
			summary.Other = append(summary.Other, issueResult)
		}
	}

	summary.CurrentRefDate = commitCurrentRef.Commit.Committer.Date.Format("2006-01-02")
	if len(config.FutureCurrentRefName) == 0 {
		summary.CurrentRefName = config.CurrentRef
	} else {
		summary.CurrentRefName = config.FutureCurrentRefName
	}

	summary.PreviousRefName = config.PreviousRef

	//// TODO Milestone?

	viewTemplate := `{{define "LineTemplate"}}- [#{{.Number}}]({{.URL}}) {{.Title |html}} ([{{.User.Login}}]({{.User.URL}})){{end}}
## [{{.CurrentRefName}}](https://github.com/{{.Owner}}/{{.RepositoryName}}/tree/{{.CurrentRefName}}) ({{.CurrentRefDate}})
[All Commits](https://github.com/{{.Owner}}/{{.RepositoryName}}/compare/{{.PreviousRefName}}...{{.CurrentRefName}})

{{if .Enhancement -}}
**Enhancements:**
{{range .Enhancement -}}
{{template "LineTemplate" .}}
{{end}}
{{- end}}
{{if .Bug -}}
**Bug fixes:**
{{range .Bug -}}
{{template "LineTemplate" .}}
{{end}}
{{- end}}
{{if .Documentation -}}
**Documentation:**
{{range .Documentation -}}
{{template "LineTemplate" .}}
{{end}}
{{- end}}
{{if .Other -}}
**Misc:**
{{range .Other -}}
{{template "LineTemplate" .}}
{{end}}
{{- end}}
	`

	tmplt := template.Must(template.New("ChangeLog").Parse(viewTemplate))

	var wr io.Writer
	if config.OutputType == "file" {
		file, err := os.Create(config.FileName)
		defer file.Close()
		check(err)
		wr = file
	} else {
		wr = os.Stdout
	}

	err := tmplt.Execute(wr, summary)
	check(err)
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	var client *github.Client
	if len(token) == 0 {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}
	return client
}

func contains(labels []github.Label, str string) bool {
	for _, label := range labels {
		if *label.Name == str {
			return true
		}
	}
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func required(field string, fieldName string) error {
	if len(field) == 0 {
		log.Fatalf("%s is mandatory.", fieldName)
	}
	return nil
}
