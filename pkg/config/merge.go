package config

import (
	"slices"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

func Merge(cfgA, cfgB *api.Config) (*api.Config, error) {
	merged := &api.Config{}
	err := mergo.Merge(merged, cfgA, mergo.WithOverride, mergo.WithAppendSlice)
	if err != nil {
		return nil, errors.New("failed to merge first provided config")
	}

	err = mergo.Merge(merged, cfgB, mergo.WithOverride, mergo.WithAppendSlice)
	if err != nil {
		return nil, errors.New("failed to merge first provided config")
	}

	removeDependencyDuplicates(merged)
	removeDictionaryDuplicates(merged)
	removeTypeDuplicates(merged)

	return merged, nil
}

func removeDictionaryDuplicates(cfg *api.Config) {
	if cfg.Commit.TokenizerOptions != nil {
		uniques := make([]*nlpapi.Dictionary, 0)
		slices.Reverse(cfg.Commit.TokenizerOptions.Dictionaries)
		for _, element := range cfg.Commit.TokenizerOptions.Dictionaries {
			found := false
			for _, unique := range uniques {
				if unique.Name == element.Name {
					found = true
					break
				}
			}

			if found {
				continue
			}

			uniques = append(uniques, element)
		}

		slices.Reverse(uniques)
		cfg.Commit.TokenizerOptions.Dictionaries = uniques
	}
}

func removeTypeDuplicates(cfg *api.Config) {
	uniques := make([]*api.CommitType, 0)
	slices.Reverse(cfg.Commit.Types)
	for _, element := range cfg.Commit.Types {
		found := false
		for _, unique := range uniques {
			if unique.Type == element.Type {
				found = true
				break
			}
		}

		if found {
			continue
		}

		uniques = append(uniques, element)
	}

	slices.Reverse(uniques)
	cfg.Commit.Types = uniques
}

func removeDependencyDuplicates(cfg *api.Config) {
	uniques := make([]*api.Dependency, 0)
	slices.Reverse(cfg.Dependencies)
	for _, element := range cfg.Dependencies {
		found := false
		for _, unique := range uniques {
			if unique.Name == element.Name {
				found = true
				break
			}
		}

		if found {
			continue
		}

		uniques = append(uniques, element)
	}

	slices.Reverse(uniques)
	cfg.Dependencies = uniques
}
