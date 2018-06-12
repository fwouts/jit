package config

import (
	"github.com/andygrunwald/go-jira"
	"strings"
)

func verify(config Config) bool {
	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(config.Jira.User),
		Password: strings.TrimSpace(config.Jira.Token),
	}
	api, err := jira.NewClient(tp.Client(), config.Jira.Host)
	if err != nil {
		return false
	}
	_, _, err = api.Issue.Search("assignee = currentUser()", nil)
	if err != nil {
		return false
	}
	return true
}
