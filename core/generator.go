package core

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/ldez/gcg/label"
	"github.com/ldez/gcg/types"
	"golang.org/x/oauth2"
)

const gitHubSearchDateLayout = "2006-01-02T15:04:05Z"

const viewTemplate = `## [{{.CurrentRefName}}](https://github.com/{{.Owner}}/{{.RepositoryName}}/tree/{{.CurrentRefName}}) ({{.CurrentRefDate}})
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

{{- define "LineTemplate" -}}
- {{.FilteredLabelNames}}{{.Issue.Title |html}} ([#{{.Issue.Number}}]({{.Issue.HTMLURL}}) by [{{.Issue.User.Login}}]({{.Issue.User.HTMLURL}}))
{{- end -}}
`

// Generate change log.
func Generate(config *types.Configuration) error {
	ctx := context.Background()

	client := newGitHubClient(ctx, config.GitHubToken)

	// Get previous ref date
	commitPreviousRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.PreviousRef, nil)
	if err != nil {
		return err
	}

	datePreviousRef := commitPreviousRef.Commit.Committer.GetDate().Add(time.Duration(config.ThresholdPreviousRef) * time.Second).Format(gitHubSearchDateLayout)

	// Get current ref version date
	commitCurrentRef, _, err := client.Repositories.GetCommit(ctx, config.Owner, config.RepositoryName, config.CurrentRef, nil)
	if err != nil {
		return err
	}

	dateCurrentRef := commitCurrentRef.Commit.Committer.GetDate().Add(time.Duration(config.ThresholdCurrentRef) * time.Second).Format(gitHubSearchDateLayout)

	// Search PR
	query := fmt.Sprintf("type:pr is:merged repo:%s/%s base:%s merged:%s..%s",
		config.Owner, config.RepositoryName, config.BaseBranch, datePreviousRef, dateCurrentRef)
	if config.Debug {
		log.Println(query)
	}

	searchOptions := &github.SearchOptions{
		Sort:        "created",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 20},
	}

	issues := searchAllIssues(ctx, client, query, searchOptions, config)

	return display(config, issues, commitCurrentRef)
}

func searchAllIssues(ctx context.Context, client *github.Client, query string, searchOptions *github.SearchOptions, config *types.Configuration) []*github.Issue {
	var allIssues []*github.Issue

	for {
		issuesSearchResult, resp, err := client.Search.Issues(ctx, query, searchOptions)
		if err != nil {
			log.Fatal(err)
		}

		for _, issue := range issuesSearchResult.Issues {
			if containsLeastOneLabel(issue.Labels, config.LabelExcludes) {
				if config.Debug {
					log.Println("Exclude:", issue.GetNumber(), issue.GetTitle())
				}

				continue
			}

			allIssues = append(allIssues, issue)
		}

		if resp.NextPage == 0 {
			break
		}

		searchOptions.Page = resp.NextPage
	}

	return allIssues
}

func display(config *types.Configuration, issues []*github.Issue, commitCurrentRef *github.RepositoryCommit) error {
	summary := &types.Summary{
		Owner:          config.Owner,
		RepositoryName: config.RepositoryName,
	}

	for _, issue := range issues {
		switch {
		case containsLabel(issue.Labels, config.LabelDocumentation):
			summary.Documentation = makeAndAppendIssueSummary(summary.Documentation, issue, config)
		case containsLabel(issue.Labels, config.LabelEnhancement):
			summary.Enhancement = makeAndAppendIssueSummary(summary.Enhancement, issue, config)
		case containsLabel(issue.Labels, config.LabelBug):
			summary.Bug = makeAndAppendIssueSummary(summary.Bug, issue, config)
		default:
			summary.Other = makeAndAppendIssueSummary(summary.Other, issue, config)
		}
	}

	sort.Sort(types.ByLabel(summary.Documentation))
	sort.Sort(types.ByLabel(summary.Enhancement))
	sort.Sort(types.ByLabel(summary.Bug))
	sort.Sort(types.ByLabel(summary.Other))

	summary.CurrentRefDate = commitCurrentRef.Commit.Committer.GetDate().Format("2006-01-02")
	if config.FutureCurrentRefName == "" {
		summary.CurrentRefName = config.CurrentRef
	} else {
		summary.CurrentRefName = config.FutureCurrentRefName
	}

	summary.PreviousRefName = config.PreviousRef

	tmplContent := viewTemplate

	if config.TemplateFile != "" {
		raw, err := os.ReadFile(config.TemplateFile)
		if err != nil {
			return err
		}

		tmplContent = string(raw)
	}

	base := template.New("ChangeLog")
	tmplt := template.Must(base.Parse(tmplContent))

	var wr io.Writer

	if config.OutputType == "file" {
		var file *os.File

		file, err := os.Create(config.FileName)
		if err != nil {
			return err
		}

		defer func() { _ = file.Close() }()

		wr = file
	} else {
		wr = os.Stdout
	}

	return tmplt.Execute(wr, summary)
}

func newGitHubClient(ctx context.Context, token string) *github.Client {
	if token == "" {
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	return github.NewClient(oauth2.NewClient(ctx, ts))
}

func containsLabel(labels []*github.Label, str string) bool {
	return slices.ContainsFunc(labels, func(lbl *github.Label) bool {
		return lbl.GetName() == str
	})
}

func containsLeastOneLabel(labels []*github.Label, values []string) bool {
	return slices.ContainsFunc(labels, func(lbl *github.Label) bool {
		return slices.Contains(values, lbl.GetName())
	})
}

func makeAndAppendIssueSummary(summaries []types.IssueSummary, issue *github.Issue, config *types.Configuration) []types.IssueSummary {
	var lbl string

	if config.DisplayLabel {
		labels := label.FilterAndTransform(
			issue.Labels,
			labelFilter(config.DisplayLabelOptions),
			trimAllPrefix(config.DisplayLabelOptions))

		if len(labels) != 0 {
			lbl = fmt.Sprintf("**[%s]** ", strings.Join(labels, ","))
		}
	}

	is := types.IssueSummary{
		FilteredLabelNames: lbl,
		Issue:              issue,
	}

	return append(summaries, is)
}

func trimAllPrefix(options *types.DisplayLabelOptions) label.NameTransform {
	if options != nil {
		return label.TrimAllPrefix(options.TrimmedPrefixes)
	}

	return label.NameIdentity
}

func labelFilter(options *types.DisplayLabelOptions) label.Predicate {
	if options == nil {
		return label.All
	}

	return label.AllMatch(
		label.FilteredBy(label.HasPrefix, options.FilteredPrefixes),
		label.ExcludedBy(label.HasPrefix, options.ExcludedPrefixes),
		label.FilteredBy(label.HasSuffix, options.FilteredSuffixes),
		label.ExcludedBy(label.HasSuffix, options.ExcludedSuffixes))
}
