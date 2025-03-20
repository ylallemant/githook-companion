package config

import (
	"io/fs"
	"os"

	"github.com/pkg/errors"
)

const (
	FailedToReadStatsFromPathFmt = "failed to read stats from path \"%s\""
)

var ErrorNoDirectory = errors.New("path target exists but is no directory")

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

func DirectoryExists(path string) (bool, fs.FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil, nil
		}

		return false, nil, errors.Wrapf(err, FailedToReadStatsFromPathFmt, path)
	}

	if fi.IsDir() {
		return true, fi, nil
	}

	return true, fi, ErrorNoDirectory
}
