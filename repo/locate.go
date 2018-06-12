package repo

import (
	"errors"
	"os"
	"path"
)

// Locate returns the absolute path of the Git repository containing the current working directory.
func Locate() (*string, error) {
	potentialGitRepoPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for potentialGitRepoPath != "/" {
		_, err := os.Stat(path.Join(potentialGitRepoPath, ".git"))
		if err == nil {
			return &potentialGitRepoPath, nil
		} else if os.IsNotExist(err) {
			potentialGitRepoPath = path.Join(potentialGitRepoPath, "..")
		} else {
			return nil, err
		}
	}
	return nil, errors.New("please call Jit from a Git repository")
}
