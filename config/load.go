package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// Load reads the Jit configuration from a Git repository.
func Load(gitRepoPath string) (*Config, error) {
	configDirPath := path.Join(gitRepoPath, ".jit")
	_, err := os.Stat(configDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create .jit directory with only read/write/exec access for current user.
			os.Mkdir(configDirPath, 0700)
		} else {
			return nil, err
		}
	}
	gitIgnorePath := path.Join(configDirPath, ".gitignore")
	_, err = os.Stat(gitIgnorePath)
	if err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(gitIgnorePath, []byte("config.yaml\n"), 0600)
		} else {
			return nil, err
		}
	}
	configFilePath := path.Join(configDirPath, "config.yaml")
	_, err = os.Stat(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err = Prompt(configFilePath, nil)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	configData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return nil, err
	}
	if config.Jira.Host == "" || config.Jira.User == "" || config.Jira.Token == "" {
		fmt.Println("Uh-oh. It looks like your Jit config is incomplete. Please bare with us.")
		config, err = Prompt(configFilePath, config)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
