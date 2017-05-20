package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

type pullRequestSummary struct {
	enhancement   []string
	documentation []string
	bug           []string
	other         []string
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
	fmt.Printf("[All Commits](https://github.com/%s/%s/compare/%s...%s)\n",
		config.Owner, config.RepositoryName, config.PreviousRef, config.CurrentRef)
	// TODO All PR?
	// TODO Milestone?
	//client.Issues.GetMilestone()
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
