package config

import (
	"io/fs"
	"os"

	"github.com/pkg/errors"
)

const (
	FailedToReadStatsFromPathFmt = "failed to read stats from path \"%s\""
)

func fileExists(path string) (bool, fs.FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil, nil
		}

		return false, nil, errors.Wrapf(err, FailedToReadStatsFromPathFmt, path)
	}

	return true, fi, nil
}
