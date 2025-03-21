package config

import (
	"path/filepath"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
)

func ParentPathFromConfig(configuration *api.Config) string {
	if configuration.ParentConfig != nil {
		if configuration.ParentConfig.Path != "" {
			relativePath := filepath.Join(configuration.ParentConfig.Path)
			path, err := environment.EnsureAbsolutePath(relativePath)
			if err != nil {
				panic(err)
			}

			return path
		}
	}

	return ""
}
