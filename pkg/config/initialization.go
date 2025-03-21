package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/command"
	"github.com/ylallemant/githook-companion/pkg/environment"
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

	return nil
}

func ensureConfigurationDirectory(path string) error {
	configurationDirectory := filepath.Join(path, api.ConfigDirectory)

	exists, _, err := DirectoryExists(configurationDirectory)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", configurationDirectory)
	}

	if !exists {
		err = os.MkdirAll(configurationDirectory, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed create directory %s", configurationDirectory)
		}
	}

	return nil
}

func ensureConfigurationFile(path string, reference *api.ParentConfig, minimalistic bool) error {
	configurationFile := filepath.Join(path, api.ConfigDirectory, api.ConfigFile)

	exists, _, err := fileExists(configurationFile)
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

	exists, _, err := DirectoryExists(path)
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

	return nil
}
