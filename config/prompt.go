package config

import (
	"errors"
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

// Prompt asks the user to fill in the configuration and saves it in `configFilePath`.
func Prompt(configFilePath string, existingConfig *Config) (*Config, error) {
	questions := []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "Enter Jira base URL:",
				Default: (func() string {
					if existingConfig != nil {
						return existingConfig.Jira.Host
					}
					return ""
				})(),
				Help: "Open Jira and copy-paste the hostname. For example, https://issues.apache.org is the base URL for Apache's Jira issue tracker.",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("The Jira base URL must be a string.")
				}
				if len(str) == 0 {
					return errors.New("The Jira base URL must be set.")
				}
				if !strings.HasPrefix(str, "http://") && !strings.HasPrefix(str, "https://") {
					return errors.New("The Jira base URL must start with http(s)://.")
				}
				return nil
			},
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "What username do you sign into Jira with?",
				Default: (func() string {
					if existingConfig != nil {
						return existingConfig.Jira.User
					}
					return ""
				})(),
			},
			Validate: survey.Required,
		},
		{
			Name: "token",
			Prompt: &survey.Input{
				Message: "Enter a Jira API token:",
				Default: (func() string {
					if existingConfig != nil {
						return existingConfig.Jira.Token
					}
					return ""
				})(),
			},
			Validate: survey.Required,
		},
	}
	answers := struct {
		Host  string
		User  string
		Token string
	}{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}
	config := Config{
		Jira: answers,
	}
	valid := verify(config)
	if !valid {
		fmt.Println("Uh-oh! We were unable to verify your configuration. Let's start over.")
		return Prompt(configFilePath, &config)
	}
	fmt.Println("Configuration looks valid!")
	configData, err := yaml.Marshal(&config)
	if err != nil {
		return nil, err
	}
	// Save the config, read/write only allowed for current user.
	ioutil.WriteFile(configFilePath, configData, 0600)
	return &config, nil
}
