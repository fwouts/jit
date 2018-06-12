package repo

import (
	"github.com/deckarep/golang-set"
	"github.com/libgit2/git2go"
	"regexp"
)

// BranchesWithJiraKeys contains a list of branches along with their Jira keys.
type BranchesWithJiraKeys struct {
	JiraKeys            mapset.Set
	JiraKeyToBranchName map[string]string
}

// FindBranchesWithJiraKeys finds local branches associated with Jira keys.
func FindBranchesWithJiraKeys(r *git.Repository) (*BranchesWithJiraKeys, error) {
	jiraKeysFromBranches := mapset.NewSet()
	jiraKeyToBranchName := make(map[string]string)
	branchIterator, err := r.NewBranchIterator(git.BranchLocal)
	if err != nil {
		return nil, err
	}
	defer branchIterator.Free()
	jiraIssueFormat, err := regexp.Compile("^([A-Z]+-[0-9]+)")
	if err != nil {
		return nil, err
	}
	err = branchIterator.ForEach(func(b *git.Branch, t git.BranchType) error {
		name, err := b.Name()
		if err != nil {
			return err
		}
		res := jiraIssueFormat.FindStringSubmatch(name)
		if len(res) > 0 {
			jiraKey := res[1]
			jiraKeysFromBranches.Add(jiraKey)
			jiraKeyToBranchName[jiraKey] = name
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BranchesWithJiraKeys{
		JiraKeys:            jiraKeysFromBranches,
		JiraKeyToBranchName: jiraKeyToBranchName,
	}, nil
}
