package config

import (
	"github.com/andygrunwald/go-jira"
	"strings"
)

// JiraClient returns a Jira API client.
func (config *Config) JiraClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(config.Jira.User),
		Password: strings.TrimSpace(config.Jira.Token),
	}
	api, err := jira.NewClient(tp.Client(), config.Jira.Host)
	if err != nil {
		return nil, err
	}
	return api, nil
}
