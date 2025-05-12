package hook

import (
	"slices"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

func Locked(hookname string, configContext api.ConfigContext) (bool, error) {
	status := false

	hookList, exist, err := config.ListGithooks(configContext.Config())
	if err != nil {
		return status, err
	}
	log.Debug().Msgf("hook list retrieved %v", hookList)
	log.Debug().Msgf("hook \"%s\" exists: %v", hookname, exist)

	if exist {
		if slices.Contains(hookList, hookname) {
			path := config.GithookLockPathFromNameAndConfig(hookname, configContext)
			exists, _, err := filesystem.FileExists(path)
			if err != nil {
				return status, err
			}

			lockType := filesystem.LockType(path)
			log.Debug().Msgf("lock type: \"%s\"", lockType)

			if exists {
				if lockType == filesystem.LockTypeTemporary {
					status, err = config.TimeLockActive(hookname, configContext)
				} else {
					status = exists
				}
				log.Debug().Msgf("lock status: %v", status)

				if err != nil {
					return status, errors.Wrapf(err, "failed to read lock status for %s", path)
				}
			}
		} else {
			return status, errors.Errorf("githook \"%s\" is unknown. possible values are %v", hookname, hookList)
		}
	} else {
		return status, errors.Errorf("no githooks are defined")
	}

	return status, nil
}
