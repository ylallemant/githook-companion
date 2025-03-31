package filesystem

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func SetPermanentLock(path string) error {
	description := "permanent lock"
	return SetPermanentLockWithDescription(path, description)
}

func SetPermanentLockWithDescription(path, description string) error {
	log.Debug().Msgf("set permanent lock at %s", path)
	return os.WriteFile(path, []byte(description), 0644)
}

func SetTimedLock(path string, duration time.Duration) error {
	description := "timed lock"
	return SetTimedLockWithDescription(path, description, duration)
}

func SetTimedLockWithDescription(path, description string, duration time.Duration) error {
	description = fmt.Sprintf(`%s
	valid until: %s`, description, time.Now().Add(duration))

	err := os.WriteFile(path, []byte(description), 0644)
	if err != nil {
		return errors.Wrap(err, "failed to add credentials to repository uri")
	}

	_, info, err := FileExists(path)
	if err != nil {
		return errors.Wrap(err, "failed to add credentials to repository uri")
	}
	endtime := info.ModTime().Add(duration)

	err = os.Chtimes(path, endtime, endtime)
	if err != nil {
		return errors.Wrap(err, "failed to add credentials to repository uri")
	}

	log.Debug().Msgf("set timed lock at %s", path)
	return nil
}

func PermanentLockExists(path string) (bool, error) {
	exists, _, err := FileExists(path)
	if err != nil {
		return false, errors.Wrap(err, "failed to add credentials to repository uri")
	}

	log.Debug().Msgf("permanent lock exists (%v) at %s", exists, path)
	return exists, nil
}

func TimeLockActive(path string) (bool, error) {
	exists, info, err := FileExists(path)
	if err != nil {
		return false, errors.Wrap(err, "failed to add credentials to repository uri")
	}

	if exists {
		validity := info.ModTime().After(time.Now())
		log.Debug().Msgf("timed lock is valid (%v) at %s", validity, path)

		if !validity {
			err = RemoveLock(path)
			if err != nil {
				return false, err
			}
		}

		return validity, nil
	}

	log.Debug().Msgf("timed lock exists (%v) at %s", exists, path)
	return false, err
}

func RemoveLock(path string) error {
	log.Debug().Msgf("remove lock at %s", path)
	return os.Remove(path)
}
