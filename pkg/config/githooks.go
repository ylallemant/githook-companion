package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

func GithooksExist(path string) (bool, error) {
	path, err := environment.EnsureAbsolutePath(path)
	if err != nil {
		return false, err
	}

	githooksDirectory := filepath.Join(path, api.GithooksDirectory)

	exists, _, err := filesystem.DirectoryExists(githooksDirectory)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check existance of %s", path)
	}

	return exists, nil
}

func GithooksPathFromConfig(configuration *api.Config) string {
	relativePath := filepath.Join(".", api.ConfigDirectory, api.GithooksDirectory)

	if configuration.GithooksDirectory != "" {
		relativePath = configuration.GithooksDirectory
	}

	if configuration.ParentConfig != nil {
		if configuration.ParentConfig.Path != "" {
			relativePath = filepath.Join(configuration.ParentConfig.Path, api.ConfigDirectory, api.GithooksDirectory)
		}
	}

	path, err := environment.EnsureAbsolutePath(relativePath)
	if err != nil {
		panic(err)
	}

	return path
}

func GithookPathFromNameAndConfig(name string, configuration *api.Config) string {
	directory := GithooksPathFromConfig(configuration)
	return filepath.Join(directory, name)
}

func ListGithooks(configuration *api.Config) ([]string, bool, error) {
	list := make([]string, 0)
	directory := GithooksPathFromConfig(configuration)

	exists, _, err := filesystem.DirectoryExists(directory)
	if err != nil {
		return list, false, errors.Wrap(err, "failed to check existance of hooks directory")
	}

	if exists {
		err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && info.Size() > 0 {
				list = append(list, info.Name())
			}

			return nil
		})
		if err != nil {
			return list, exists, errors.Wrap(err, "failed list content of hooks directory")
		}
	}

	return list, exists, nil
}
