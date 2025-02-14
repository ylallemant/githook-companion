package config

import "github.com/ylallemant/githooks-butler/pkg/api"

func CommitTypes(config *api.Config) []string {
	types := make([]string, 0)

	for _, current := range config.Commit.Types {
		types = append(types, current.Type)
	}

	return types
}
