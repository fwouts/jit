package config

import (
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// Load reads the Jit configuration from a home directory.
func Load(homePath string, gitRepoPath string) (*Config, error) {
	homeConfigDirPath := path.Join(homePath, ".jit")
	configFilePath := path.Join(homeConfigDirPath, "config.yaml")

	// We used to store the Jit config directly within the Git repo. If we find it
	// there, we'll automatically move it to the home directory.
	oldGitConfigFilePath := path.Join(gitRepoPath, ".jit", "config.yaml")

	_, err := os.Stat(homeConfigDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create .jit directory with only read/write/exec access for current user.
			os.Mkdir(homeConfigDirPath, 0700)
		} else {
			return nil, err
		}
	}
	_, err = os.Stat(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Check if there's an old config file inside the Git repo.
			_, oldGitConfigErr := os.Stat(oldGitConfigFilePath)
			if oldGitConfigErr == nil {
				// There is one. Let's move it to the home directory.
				os.Rename(oldGitConfigFilePath, configFilePath)
				fmt.Println("We've moved .jit/config.yaml to your home directory.")
				fmt.Println(color.MagentaString("You can now safely delete the .jit directory in this Git repository."))
			} else {
				// There is no existing config file. Create a new one.
				_, err = Prompt(configFilePath, nil)
				if err != nil {
					return nil, err
				}
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
