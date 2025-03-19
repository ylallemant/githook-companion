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

func GithooksPathFromConfig(configuration *api.Config) string {
	if configuration.GithooksDirectory != "" {
		path, err := environment.EnsureAbsolutePath(configuration.GithooksDirectory)
		if err != nil {
			panic(err)
		}

		return path
	}

	if configuration.ConfigReference != nil {
		if configuration.ConfigReference.Path != "" {
			relativePath := filepath.Join(configuration.ConfigReference.Path, api.GithooksDirectory)
			path, err := environment.EnsureAbsolutePath(relativePath)
			if err != nil {
				panic(err)
			}

			return path
		}
	}

	path, err := environment.EnsureAbsolutePath(filepath.Join(".", api.GithooksDirectory))
	if err != nil {
		panic(err)
	}

	return path
}
