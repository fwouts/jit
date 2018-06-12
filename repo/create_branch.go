package repo

import (
	"github.com/libgit2/git2go"
	"gopkg.in/AlecAivazis/survey.v1"
)

// CreateBranch prompts the user for a branch name based on a Jira key and creates it.
func CreateBranch(repo *git.Repository, jiraKey string) (*string, error) {
	var branchName string
	survey.AskOne(&survey.Input{
		Message: "Enter a branch name:",
		Default: jiraKey,
	}, &branchName, nil)
	master, err := repo.LookupBranch("origin/master", git.BranchRemote)
	if err != nil {
		return nil, err
	}
	defer master.Free()
	lastCommit, err := repo.LookupCommit(master.Target())
	if err != nil {
		return nil, err
	}
	defer lastCommit.Free()
	branch, err := repo.CreateBranch(branchName, lastCommit, false)
	if err != nil {
		return nil, err
	}
	defer branch.Free()
	return &branchName, err
}
