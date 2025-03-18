package config

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
)

func GithooksExist(path string) (bool, error) {
	path, err := environment.EnsureAbsolutePath(path)
	if err != nil {
		return false, err
	}

	githooksDirectory := filepath.Join(path, api.GithooksDirectory)

	exists, _, err := DirectoryExists(githooksDirectory)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check existance of %s", path)
	}

	return exists, nil
}
