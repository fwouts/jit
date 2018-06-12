package repo

import (
	"github.com/libgit2/git2go"
)

// CheckoutBranch checks out a branch safely in the provided repository.
func CheckoutBranch(repo *git.Repository, branchName string) error {
	branch, err := repo.LookupBranch(branchName, git.BranchLocal)
	if err != nil {
		return err
	}
	defer branch.Free()

	commit, err := repo.LookupCommit(branch.Target())
	if err != nil {
		return err
	}
	defer commit.Free()

	tree, err := repo.LookupTree(commit.TreeId())
	if err != nil {
		return err
	}
	defer tree.Free()

	err = repo.CheckoutTree(tree, &git.CheckoutOpts{
		Strategy: git.CheckoutSafe,
	})
	if err != nil {
		return err
	}
	repo.SetHead("refs/heads/" + branchName)
	return nil
}
