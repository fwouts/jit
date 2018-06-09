package main

import (
	"flag"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/deckarep/golang-set"
	"github.com/fatih/color"
	"github.com/libgit2/git2go"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	repoPath := flag.String("repo", cwd, "Path to the git repository")
	jiraHost := flag.String("jira-host", "", "URL of the Jira Cloud host")
	jiraUsername := flag.String("jira-user", "", "Jira username")
	jiraToken := flag.String("jira-token", "", "Jira API token")
	flag.Parse()

	if *jiraHost == "" {
		log.Fatal("Please set a Jira host with --jira-host.")
	}
	if *jiraUsername == "" {
		log.Fatal("Please set a Jira username with --jira-user.")
	}
	if *jiraToken == "" {
		log.Fatal("Please set a Jira API token with --jira-token.")
	}

	repo, err := git.OpenRepository(*repoPath)
	if err != nil {
		log.Fatal(err)
	}
	branchPrefixes := mapset.NewSet()
	i, err := repo.NewBranchIterator(git.BranchLocal)
	if err != nil {
		log.Fatal(err)
	}
	jiraIssueFormat, err := regexp.Compile("^([A-Z]+-[0-9]+)")
	if err != nil {
		log.Fatal(err)
	}
	err = i.ForEach(func(b *git.Branch, t git.BranchType) error {
		name, err := b.Name()
		if err != nil {
			return err
		}
		res := jiraIssueFormat.FindStringSubmatch(name)
		if len(res) > 0 {
			branchPrefixes.Add(res[1])
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(*jiraUsername),
		Password: strings.TrimSpace(*jiraToken),
	}
	api, err := jira.NewClient(tp.Client(), *jiraHost)
	if err != nil {
		log.Fatal(err)
	}
	issues, _, err := api.Issue.Search("assignee = currentUser()", nil)
	if err != nil {
		log.Fatal(err)
	}
	choices := make([]string, len(issues))
	for i, issue := range issues {
		var colouredPriority string
		switch issue.Fields.Priority.Name {
		case "Low":
			colouredPriority = color.GreenString(issue.Fields.Priority.Name)
		case "Medium":
			colouredPriority = color.BlueString(issue.Fields.Priority.Name)
		case "High":
			colouredPriority = color.MagentaString(issue.Fields.Priority.Name)
		case "Critical":
		case "Blocker":
			colouredPriority = color.RedString(issue.Fields.Priority.Name)
		default:
			colouredPriority = color.GreenString(issue.Fields.Priority.Name)
		}
		var summary string
		if branchPrefixes.Contains(issue.Key) {
			summary = fmt.Sprintf("*%s %v (%s): %+v", issue.Key, colouredPriority, issue.Fields.Type.Name, issue.Fields.Summary)
		} else {
			summary = fmt.Sprintf("%s %v (%s): %+v", issue.Key, colouredPriority, issue.Fields.Type.Name, issue.Fields.Summary)
		}
		choices[i] = summary
	}
	questions := []*survey.Question{
		{
			Name: "issue",
			Prompt: &survey.Select{
				Message: "Choose an issue:",
				Options: choices,
			},
		},
	}
	answers := struct {
		Issue string
	}{}
	err = survey.Ask(questions, &answers)
	if err != nil {
		log.Fatal(err)
	}
}
