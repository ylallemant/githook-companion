package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/command"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

func EnsureConfiguration(path string, reference *api.ParentConfig, minimalistic bool) error {
	err := ensureConfigurationDirectory(path)
	if err != nil {
		return err
	}

	err = ensureConfigurationFile(path, reference, minimalistic)
	if err != nil {
		return err
	}

	for _, directory := range api.ConfigProcessingDirectories {
		err = filesystem.EnsureDirectory(filepath.Join(path, api.ConfigDirectory, directory))
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureConfigurationDirectory(path string) error {
	configurationDirectory := DirectoryPathFromBase(path)
	return filesystem.EnsureDirectory(configurationDirectory)
}

func ensureConfigurationFile(path string, reference *api.ParentConfig, minimalistic bool) error {
	configurationFile := FilePathFromBase(path)

	exists, _, err := filesystem.FileExists(configurationFile)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", configurationFile)
	}

	if !exists {
		cfg := Minimalistic()

		if !minimalistic {
			cfg = Default()
		}

		if reference != nil {
			cfg.ParentConfig = reference
		}

		content, err := ToYAML(cfg)
		if err != nil {
			return errors.Wrap(err, "failed to marshal default config")
		}

		err = os.WriteFile(configurationFile, content, 0644)
		if err != nil {
			return errors.Wrapf(err, "failed to write default config at %s", configurationFile)
		}
	}

	return nil
}

func EnsureReference(reference *api.ParentConfig) error {
	path, err := environment.EnsureAbsolutePath(reference.Path)
	if err != nil {
		return err
	}

	exists, _, err := filesystem.DirectoryExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", path)
	}

	if !exists {
		fmt.Println("clone reference repository", reference.GitRepository)
		git := command.New("git")
		git.AddArg("clone")
		git.AddArg(reference.GitRepository)
		git.AddArg(path)

		_, err = git.Execute()
		if err != nil {
			return errors.Wrapf(err, "failed to clone reference repository %s", reference.GitRepository)
		}
	}

	for _, directory := range api.ConfigProcessingDirectories {
		err = filesystem.EnsureDirectory(filepath.Join(path, api.ConfigDirectory, directory))
		if err != nil {
			return err
		}
	}

	return nil
}
