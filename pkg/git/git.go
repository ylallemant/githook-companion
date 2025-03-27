package git

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/command"
)

func CurrentBranchFromPath(path string) (string, error) {
	git := command.New("git")

	if path != "" {
		git.AddArg("-C")
		git.AddArg(path)
	}

	git.AddArg("rev-parse")
	git.AddArg("--abbrev-ref")
	git.AddArg("HEAD")

	branch, err := git.Execute()
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch git branch at %s", path)
	}

	return branch, nil
}

func CurrentBranch() (string, error) {
	return CurrentBranchFromPath("")
}

func CommitHashFromPath(path, branch string) (string, error) {
	git := command.New("git")

	if path != "" {
		git.AddArg("-C")
		git.AddArg(path)
	}

	git.AddArg("rev-parse")
	git.AddArg(branch)

	hash, err := git.Execute()
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch git hash for branch \"%s\" at %s", branch, path)
	}

	return hash, nil
}

func CommitHash() (string, error) {
	branch, err := CurrentBranch()
	if err != nil {
		return "", err
	}

	return CommitHashFromPath("", branch)
}

func Pull(path string) error {
	git := command.New("git")

	if path != "" {
		git.AddArg("-C")
		git.AddArg(path)
	}

	git.AddArg("pull")

	result, err := git.Execute()
	if err != nil {
		return errors.Wrapf(err, "failed to pull git repository at %s", path)
	}

	fmt.Println(result, "at", path)

	return nil
}
