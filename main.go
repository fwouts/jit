package main

import (
	"flag"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/fatih/color"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"strings"
)

func main() {
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
		choices[i] = fmt.Sprintf("%s %v (%s): %+v", issue.Key, colouredPriority, issue.Fields.Type.Name, issue.Fields.Summary)
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
