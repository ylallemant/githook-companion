package config

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	FailedToReadStatsFromPathFmt = "failed to read stats from path \"%s\""
)

func homeDir() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get home directory")
	}

	return dirname, nil
}

func localDir() (string, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get local directory")
	}
	return dirname, nil
}

func ensureAbsolutePath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		local, err := localDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(local, path)
	}

	return path, nil
}

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
