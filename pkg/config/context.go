package config

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/binary"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

const (
	configLockDescription  = `lock used to mitigate Git provider request limits. Use PAT to enable more requests`
	networkLockDescription = `lock used to mitigate network connectivity problems`
)

func Context(fallbackToDefault bool) (*configContext, error) {
	var err error
	var basePath string

	localBasePath, err := GetLocalBasePath()
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("local base path \"%s\"", localBasePath)

	localBasePathExists, _, err := filesystem.DirectoryExists(localBasePath)
	if err != nil {
		return nil, err
	}

	if localBasePathExists {
		basePath = localBasePath
	} else {
		globalBasePath, err := GetGlobalBasePath()
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("global base path \"%s\"", globalBasePath)

		globalBasePathExists, _, err := filesystem.DirectoryExists(globalBasePath)
		if err != nil {
			return nil, err
		}

		if globalBasePathExists {
			basePath = globalBasePath
		}
	}

	if basePath == "" {
		return nil, errors.Errorf("failed to find a configuration locally or globally")
	}

	log.Debug().Msgf("base path set to \"%s\"", basePath)

	return ContextFromPath(basePath, fallbackToDefault)
}

func ContextFromPath(customPath string, fallbackToDefault bool) (*configContext, error) {
	var err error
	log.Debug().Msgf("use base path \"%s\"", customPath)

	ctx := new(configContext)
	ctx.localPath = customPath

	ctx.localConfig, err = LoadFromBase(customPath, false)
	if err != nil {
		return nil, err
	}

	if ctx.localConfig == nil {
		if !fallbackToDefault {
			return nil, errors.Errorf("failed to find a configuration at %s", customPath)
		}

		log.Debug().Msg("use default configuration")
		ctx.config = Default()
		return ctx, nil
	}

	if ctx.localConfig.ParentConfig != nil {
		log.Debug().Msgf("parent configuration referenced at \"%s\"", ctx.localConfig.ParentConfig.GitRepository)

		ctx.parentPath, err = environment.EnsureAbsolutePath(ctx.localConfig.ParentConfig.Path)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("parent configuration directory at \"%s\"", ctx.parentPath)

		ctx.parentConfig, err = LoadFromBase(ctx.parentPath, false)
		if err != nil {
			return nil, err
		}
	}

	if ctx.localConfig != nil && ctx.parentConfig != nil {
		ctx.config, err = Merge(ctx.parentConfig, ctx.localConfig)
		if err != nil {
			return nil, errors.Wrap(err, "failed to merge parent and local configurations")
		}
		log.Debug().Msg("use merged configuration")
	} else if ctx.localConfig != nil {
		log.Debug().Msg("use local configuration")
		ctx.config = ctx.localConfig
	}

	syncConfig := ctx.Config().ToolSync

	if syncConfig == nil {
		syncConfig = defaultSync()
	}

	if syncConfig.Parent || syncConfig.Binary {
		reloadNeeded := false

		if syncConfig.Parent {
			// ensure parent config has the latest version
			reloadNeeded, err = EnsureVersionSync(ctx)
			if err != nil {
				log.Warn().Msgf("could not ensure parent configuration version sync: %s", err.Error())
			}
		} else {
			log.Debug().Msgf("parent sync disabled")
		}

		if syncConfig.Binary {
			// ensure binary has the latest version
			binaryInSync, err := binary.VersionsInSync()
			if err != nil {
				log.Warn().Msgf("could not ensure binary version sync: %s", err.Error())
			}

			if !binaryInSync {
				return nil, errors.New("new version of the binary is available, please upgrade with: githook-companion upgrade")
			}
		} else {
			log.Debug().Msgf("binary sync disabled")
		}

		if reloadNeeded {
			return ContextFromPath(customPath, fallbackToDefault)
		}
	} else {
		log.Warn().Msgf("auto sync is disabled")
	}

	return ctx, nil
}

var _ api.ConfigContext = &configContext{}

type configContext struct {
	localPath    string
	parentPath   string
	localConfig  *api.Config
	parentConfig *api.Config
	config       *api.Config
}

func (i *configContext) LocalPath() string {
	return i.localPath
}

func (i *configContext) ParentPath() string {
	return i.parentPath
}

func (i *configContext) LocalConfig() *api.Config {
	return i.localConfig
}

func (i *configContext) ParentConfig() *api.Config {
	return i.parentConfig
}

func (i *configContext) Config() *api.Config {
	return i.config
}

func (i *configContext) HasParent() bool {
	return i.parentPath != ""
}
