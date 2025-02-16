package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githooks-butler/pkg/api"
	"gopkg.in/yaml.v3"
)

func Get(path string) (*api.Config, error) {
	var err error

	if path != "" {
		path, err = ensureAbsolutePath(path)
		if err != nil {
			return nil, err
		}

		return load(path, true)
	}

	local, err := localDir()
	if err != nil {
		return nil, err
	}

	path = filepath.Join(local, api.ConfigDirectory, api.ConfigFile)
	config, err := load(path, false)
	if err != nil {
		return nil, err
	}
	if config != nil {
		return config, nil
	}

	home, err := homeDir()
	if err != nil {
		return nil, err
	}

	path = filepath.Join(home, api.ConfigDirectory, api.ConfigFile)
	config, err = load(path, false)
	if err != nil {
		return nil, err
	}
	if config != nil {
		return config, nil
	}

	return Default(), nil
}

func load(path string, strict bool) (*api.Config, error) {
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

func ToYAML(config *api.Config) ([]byte, error) {
	return yaml.Marshal(config)
}

func ToJSON(config *api.Config) ([]byte, error) {
	return json.Marshal(config)
}

func ToPrettyJSON(config *api.Config) ([]byte, error) {
	return json.MarshalIndent(config, "", "  ")
}
