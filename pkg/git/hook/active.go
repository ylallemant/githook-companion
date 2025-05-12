package hook

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/git/merge"
)

func Active(hookname string, configContext api.ConfigContext) (bool, error) {
	status := false

	// check if a lock exists
	locked, err := Locked(hookname, configContext)
	if err != nil {
		return status, errors.Wrap(err, "failed to check lock")
	}

	// check is a merge process is active
	mergeActive, err := merge.Active()
	if err != nil {
		return status, errors.Wrap(err, "failed to check for merge status")
	}

	log.Debug().Msgf("hook is locked:       %v", locked)
	log.Debug().Msgf("merge process active: %v", mergeActive)

	// githook only active if none of the above is true
	status = !locked && !mergeActive

	return status, nil
}
