package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"gopkg.in/yaml.v3"
)

func Get(path string) (*api.Config, error) {
	var err error

	if path == "" {
		// set path to home directory
		home, err := environment.Home()
		if err != nil {
			return nil, err
		}

		path = filepath.Join(home, api.ConfigDirectory, api.ConfigFile)
	}

	mainConfig, err := load(path, false)
	if err != nil {
		return nil, err
	}

	local, err := environment.CurrentDirectory()
	if err != nil {
		return nil, err
	}

	path = filepath.Join(local, api.ConfigDirectory, api.ConfigFile)
	localConfig, err := load(path, false)
	if err != nil {
		return nil, err
	}

	if mainConfig != nil && localConfig != nil {
		merged, err := Merge(mainConfig, localConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge global and local configurations")
		}

		fmt.Println("merged config")
		return merged, nil
	}

	if localConfig != nil {
		fmt.Println("local config")
		return localConfig, nil
	}

	if mainConfig != nil {
		fmt.Println("main config")
		return mainConfig, nil
	}

	fmt.Println("default config")

	return Default(), nil
}

func load(path string, strict bool) (*api.Config, error) {
	var err error

	path, err = environment.EnsureAbsolutePath(path)
	if err != nil {
		return nil, err
	}

	_, stats, err := fileExists(path)
	if err != nil {
		return nil, err
	}

	if stats == nil && strict {
		return nil, errors.Errorf("no config found at %s", path)
	}

	if stats != nil {
		buf, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "config could not be read from %s", path)
		}

		cfg := &api.Config{}
		err = yaml.Unmarshal(buf, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "config could not be loaded from %s", path)
		}

		return cfg, nil
	}

	return nil, nil
}

func GetCommitTypes(config *api.Config) []string {
	commitTypes := make([]string, 0)

	for _, commitType := range config.Commit.Types {
		commitTypes = append(commitTypes, commitType.Type)
	}

	return commitTypes
}

func ToYAML(config *api.Config) ([]byte, error) {
	return yaml.Marshal(config)
}

func ToJSON(config *api.Config) ([]byte, error) {
	return json.Marshal(config)
}

func ToPrettyJSON(config *api.Config) ([]byte, error) {
	return json.MarshalIndent(config, "", "  ")
}
