package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/binary"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
	"gopkg.in/yaml.v3"
)

var Current *api.Config

func GetLocally() (*api.Config, error) {
	path, err := GetLocalFilePath()
	if err != nil {
		return nil, err
	}

	localConfig, err := Load(path, false)
	if err != nil {
		return nil, err
	}

	if localConfig == nil {
		return nil, errors.Wrapf(api.ConfigurationNotFound, "no local configuration at %s", path)
	}

	return localConfig, nil
}

func GetGlobally() (*api.Config, error) {
	path, err := GetGlobalFilePath()
	if err != nil {
		return nil, err
	}

	mainConfig, err := Load(path, false)
	if err != nil {
		return nil, err
	}

	if mainConfig == nil {
		return nil, errors.Wrapf(api.ConfigurationNotFound, "no global configuration at %s", path)
	}

	return mainConfig, nil
}

func Get() (*api.Config, error) {
	cfg, err := GetLocally()
	if err != nil && !errors.Is(err, api.ConfigurationNotFound) {
		return nil, err
	}

	if cfg == nil {
		cfg, err = GetGlobally()
		if err != nil && !errors.Is(err, api.ConfigurationNotFound) {
			return nil, err
		}
	}

	if cfg == nil {
		return nil, errors.Wrap(api.ConfigurationNotFound, "failed to find a local or global configuration. use the \"init\" command to create one")
	}

	var referenceCfg *api.Config
	if cfg.ParentConfig != nil {
		path := filepath.Join(cfg.ParentConfig.Path, api.ConfigDirectory, api.ConfigFile)
		path, err := environment.EnsureAbsolutePath(path)
		if err != nil {
			return nil, err
		}

		referenceCfg, err = Load(path, true)
		if err != nil {
			return nil, err
		}
	}

	if cfg != nil && referenceCfg != nil {
		merged, err := Merge(referenceCfg, cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge global and local configurations")
		}

		return merged, nil
	}

	if cfg != nil {
		Current = cfg
		return cfg, nil
	}

	return Default(), nil
}

func LoadFromBase(path string, strict bool) (*api.Config, error) {
	path = FilePathFromBase(path)
	return Load(path, strict)
}

func Load(path string, strict bool) (*api.Config, error) {
	var err error

	path, err = environment.EnsureAbsolutePath(path)
	if err != nil {
		return nil, err
	}

	_, stats, err := filesystem.FileExists(path)
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

		if cfg.Kind != api.ConfigKind {
			return nil, errors.Errorf(
				"configuration kind should be \"%s\" but found \"%s\" at %s",
				api.ConfigKind,
				cfg.Kind,
				path,
			)
		}

		log.Debug().Msgf("configuration parsed from: %s", path)

		compatible, configVersion, binaryVersion := CompatibleToBinary(cfg)
		log.Debug().Msgf("configuration (%s) and binary (%s) versions compatible: %v", configVersion, binaryVersion, compatible)
		if !compatible {
			return nil, errors.Errorf(
				"upgrade your binary: its version is out of sync with the minimum required by the configuration: %s < %s: configuration at %s",
				binaryVersion,
				configVersion,
				path,
			)
		}

		return Merge(Minimalistic(), cfg)
	}

	return nil, nil
}

func Remove(path string) error {
	exists, _, err := filesystem.DirectoryExists(path)
	if err != nil {
		return err
	}

	if exists {
		err = os.RemoveAll(path)
		if err != nil {
			return errors.Wrapf(err, "failed to remove configuration at %s", path)
		}
	}

	return nil
}

func GetCommitTypes(types []*api.CommitType) []string {
	commitTypes := make([]string, 0)

	for _, commitType := range types {
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

func CompatibleToBinary(config *api.Config) (compatible bool, configVersionString, binaryVersionString string) {
	binaryVersionString = binary.Semver()

	if binaryVersionString == "n/a" {
		// only possible locally on "go run"
		binaryVersionString = api.ConfigVersion
	}

	configVersionString = config.Version

	if configVersionString == "" {
		// old configuration have no version
		configVersionString = api.ConfigVersion
	}

	binaryVersion := semver.MustParse(binaryVersionString)
	configVersion := semver.MustParse(configVersionString)

	return configVersion.LessThanEqual(binaryVersion), configVersionString, binaryVersionString
}
