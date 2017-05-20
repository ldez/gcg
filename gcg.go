package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/containous/flaeg"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	DefaultBaseBranch         = "master"
	DefaultEnhancementLabel   = "enhancement"
	DefaultDocumentationLabel = "documentation"
	DefaultBugLabel           = "bug"
)

type Configuration struct {
	CurrentRef           string `long:"current-ref" short:"c" description:"Current commit reference. Can be a tag, a branch, a SHA."`
	PreviousRef          string `long:"previous-ref" short:"p" description:"Previous commit reference. Can be a tag, a branch, a SHA."`
	BaseBranch           string `long:"base-branch" short:"b" description:"Base branch name. PR branch destination."`
	FutureCurrentRefName string `long:"future-ref-name" short:"f" description:"TODO"`
	Owner                string `short:"o" description:"Repository owner."`
	RepositoryName       string `long:"repo-name" short:"r" description:"Repository name."`
	GithubToken          string `long:"token" short:"t" description:"GitHub Token"`
	Milestone            string `short:"m" description:""`
	LabelExclude         string `long:"exclude-label" short:"ex" description:"Label to exclude."`
	LabelEnhancement     string `long:"enhancement-label" short:"el" description:"Enhancement Label."`
	LabelDocumentation   string `long:"doc-label" short:"dl" description:"Documentation Label."`
	LabelBug             string `long:"bug-label" short:"bl" description:"Bug Label."`
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
	}

	rootCmd := &flaeg.Command{
		Name:                  "gcg",
		Description:           `GCG is a GitHub Changelog Generator.`,
		Config:                config,
		DefaultPointersConfig: &Configuration{},
		Run: func() error {
			fmt.Printf("Run GCG command with config : %+v\n", config)
			run(config)
			return nil
		},
	}

	flag := flaeg.New(rootCmd, os.Args[1:])
	flag.Run()
}

func run(config *Configuration) {
	ctx := context.Background()

	var client *github.Client
	if len(config.GithubToken) == 0 {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.GithubToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	}

	// Get previous ref date
	commitPreviousRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.PreviousRef)
	if err != nil {
		log.Fatal(err)
	}
	datePreviousRef := commitPreviousRef.Commit.Author.Date.Format("2006-01-02T15:04:05Z")

	// Get current ref version date
	commitCurrentRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.CurrentRef)
	if err != nil {
		log.Fatal(err)
	}
	dateCurrentRef := commitCurrentRef.Commit.Author.Date.Format("2006-01-02T15:04:05Z")

	// Search PR
	query := fmt.Sprintf("type:pr is:merged repo:%s/%s base:%s merged:%s..%s",
		config.Owner, config.RepositoryName, config.BaseBranch, datePreviousRef, dateCurrentRef)
	log.Println(query)

	searchOptions := &github.SearchOptions{
		Sort:        "created",
		Order:       "asc",
		ListOptions: github.ListOptions{PerPage: 20},
	}

	var allSearchResult []github.Issue
	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			log.Fatal(err)
		}
		for _, issueResult := range issuesSearchResult.Issues {
			if contains(issueResult.Labels, config.LabelExclude) {
				//log.Println("Exclude:", *issueResult.Number, *issueResult.Title)
			} else {
				allSearchResult = append(allSearchResult, issueResult)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		searchOptions.Page = resp.NextPage
	}
	display(config, allSearchResult, commitCurrentRef)
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

	summary.CurrentRefDate = commitCurrentRef.Commit.Author.Date.Format("2006-01-02")
	if len(config.FutureCurrentRefName) == 0 {
		summary.CurrentRefName = config.CurrentRef
	} else {
		summary.CurrentRefName = config.FutureCurrentRefName
	}

	summary.PreviousRefName = config.PreviousRef

	//// TODO All PR?
	//// TODO Milestone?

	viewTemplate := `
## [{{.CurrentRefName}}](https://github.com/{{.Owner}}/{{.RepositoryName}}/tree/{{.CurrentRefName}}) ({{.CurrentRefDate}})

[All Commits](https://github.com/{{.Owner}}/{{.RepositoryName}}/compare/{{.PreviousRefName}}...{{.CurrentRefName}})

{{if .Enhancement -}}
Enhancements:
{{range .Enhancement -}}
- {{.Title}} [#{{.Number}}]({{.URL}}) ([{{.User.Login}}]({{.User.URL}}))
{{end}}
{{- end}}
{{if .Bug -}}
Bug fixes:
{{range .Bug -}}
- {{.Title}} [#{{.Number}}]({{.URL}}) ([{{.User.Login}}]({{.User.URL}}))
{{end}}
{{- end}}
{{if .Documentation -}}
Documentation:
{{range .Documentation -}}
- {{.Title}} [#{{.Number}}]({{.URL}}) ([{{.User.Login}}]({{.User.URL}}))
{{end}}
{{- end}}
{{if .Other -}}
Misc:
{{range .Other -}}
- {{.Title}} [#{{.Number}}]({{.URL}}) ([{{.User.Login}}]({{.User.URL}}))
{{end}}
{{- end}}
	`
	tmplt, err := template.New("ChangeLog").Parse(viewTemplate)
	if err != nil {
		log.Fatal(err)
	}
	err = tmplt.Execute(os.Stdout, summary)
	if err != nil {
		log.Fatal(err)
	}
}

func contains(labels []github.Label, str string) bool {
	for _, label := range labels {
		if *label.Name == str {
			return true
		}
	}
	return false
}
