package config

import (
	"encoding/json"
	"fmt"

	"github.com/ylallemant/githooks-butler/pkg/api"
	"gopkg.in/yaml.v3"
)

func Get(path string) (*api.Config, error) {
	if path != "" {
		return load(path)
	}

	return Default(), nil
}

func load(path string) (*api.Config, error) {
	fmt.Println("load config from", path)
	return Default(), nil
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
