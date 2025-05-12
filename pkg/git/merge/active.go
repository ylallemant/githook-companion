package merge

import (
	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/command"
)

const expected = "Already up to date."

func Active() (bool, error) {
	git := command.New("git")

	git.AddArg("merge")
	git.AddArg("HEAD")

	output, err := git.Execute()
	if err != nil {
		return false, errors.Wrapf(err, "failed to check for a merge process")
	}

	if output != expected {
		return true, nil
	}

	return false, nil
}
