package main

import (
	"fmt"
	"github.com/libgit2/git2go"
	"github.com/zenclabs/jit/config"
	"github.com/zenclabs/jit/repo"
	"github.com/zenclabs/jit/ui"
	"github.com/zenclabs/jit/versioning"
	"gopkg.in/AlecAivazis/survey.v1"
	"log"
	"os/user"
)

func main() {
	versioning.CheckNewRelease()

	gitRepoPath, err := repo.Locate()
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	c, err := config.Load(usr.HomeDir, *gitRepoPath)
	if err != nil {
		log.Fatal(err)
	}

	r, err := git.OpenRepository(*gitRepoPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Free()
	branches, err := repo.FindBranchesWithJiraKeys(r)
	if err != nil {
		log.Fatal(err)
	}

	api, err := c.JiraClient()
	if err != nil {
		log.Fatal(err)
	}
	issues, _, err := api.Issue.Search("assignee = currentUser() order by created desc", nil)
	if err != nil {
		log.Fatal(err)
	}
	choices := make([]string, len(issues))
	for i, issue := range issues {
		choices[i] = ui.IssueSummary(&issue, branches.JiraKeys.Contains(issue.Key))
	}
	var pickedIssue string
	err = survey.AskOne(&survey.Select{
		Message: "Choose an issue:",
		Options: choices,
	}, &pickedIssue, nil)
	if err != nil {
		log.Fatal(err)
	}
	jiraKey := ui.JiraKeyFromIssueSummary(pickedIssue)
	var branchName *string
	if branches.JiraKeys.Contains(jiraKey) {
		// Switch to the existing branch.
		b := branches.JiraKeyToBranchName[jiraKey]
		branchName = &b
	} else {
		// Create a new branch.
		branchName, err = repo.CreateBranch(r, jiraKey)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created a new branch %s\n", *branchName)
	}
	err = repo.CheckoutBranch(r, *branchName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Checked out branch %s\n", *branchName)
}
