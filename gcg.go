package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	CurrentRef         string `description:"Current commit reference. Can be a tag, a branch, a SHA."`
	PreviousRef        string `description:"Previous commit reference. Can be a tag, a branch, a SHA."`
	BaseBranch         string `description:"Base branch name. PR branch destination."`
	Owner              string `description:"Repository owner."`
	RepositoryName     string `description:"Repository name."`
	GithubToken        string `description:"GitHub Token"`
	Milestone          string `description:""`
	LabelExclude       string `description:""`
	LabelEnhancement   string `description:""`
	LabelDocumentation string `description:""`
	LabelBug           string `description:""`
}

type pullRequestSummary struct {
	enhancement   []string
	documentation []string
	bug           []string
	other         []string
}

func main() {
	config := &Configuration{
		CurrentRef:         "v1.3.0-rc1",
		PreviousRef:        "v1.2.0-rc1",
		BaseBranch:         "master",
		Owner:              "containous",
		RepositoryName:     "traefik",
		GithubToken:        "",
		LabelExclude:       "area/infrastructure",
		LabelEnhancement:   "kind/enhancement",
		LabelDocumentation: "area/documentation",
		LabelBug:           "kind/bug/fix",
	}

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
	fmt.Println(query)

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
		allSearchResult = append(allSearchResult, issuesSearchResult.Issues...)
		if resp.NextPage == 0 {
			break
		}
		searchOptions.Page = resp.NextPage
	}

	summary := &pullRequestSummary{}

	// Get PR information
	for _, issueResult := range allSearchResult {
		line := fmt.Sprintf("- %s [#%v](%s) ([%s](%s))",
			strings.TrimSpace(*issueResult.Title), *issueResult.Number, *issueResult.URL, *issueResult.User.Login, *issueResult.User.URL)

		if contains(issueResult.Labels, config.LabelExclude) {
			// skip
			fmt.Println("Exclude:", *issueResult.Number, *issueResult.Title)
		} else if contains(issueResult.Labels, config.LabelDocumentation) {
			summary.documentation = append(summary.documentation, line)
		} else if contains(issueResult.Labels, config.LabelEnhancement) {
			summary.enhancement = append(summary.enhancement, line)
		} else if contains(issueResult.Labels, config.LabelBug) {
			summary.bug = append(summary.bug, line)
		} else {
			summary.other = append(summary.other, line)
		}
	}

	fmt.Printf("## [%s](TODO) (%s)\n\n", config.CurrentRef, commitCurrentRef.Commit.Author.Date.Format("2006-01-02"))

	// All commits
	fmt.Printf("[All Commits](https://github.com/%s/%s/compare/%s...%s\n)",
		config.Owner, config.RepositoryName, config.PreviousRef, config.CurrentRef)
	// TODO All PR?
	// TODO Milestone?
	// https://github.com/containous/traefik/milestone/5 <- ID not text
	fmt.Printf("Milestone [%s](TODO)\n", config.Milestone)
	fmt.Println("")

	if len(summary.enhancement) != 0 {
		fmt.Println("Enhancements:")
		for _, p := range summary.enhancement {
			fmt.Println(p)
		}
		fmt.Println("")
	}

	if len(summary.bug) != 0 {
		fmt.Println("Bug fixes:")
		for _, p := range summary.bug {
			fmt.Println(p)
		}
		fmt.Println("")
	}

	if len(summary.documentation) != 0 {
		fmt.Println("Documentation:")
		for _, p := range summary.documentation {
			fmt.Println(p)
		}
		fmt.Println("")
	}

	if len(summary.other) != 0 {
		fmt.Println("Others:")
		for _, p := range summary.other {
			fmt.Println(p)
		}
		fmt.Println("")
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
