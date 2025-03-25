package filesystem

import (
	"io/fs"
	"os"

	"github.com/pkg/errors"
)

const (
	FailedToReadStatsFromPathFmt = "failed to read stats from path \"%s\""
)

var ErrorNoDirectory = errors.New("path target exists but is no directory")

func FileExists(path string) (bool, fs.FileInfo, error) {
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

func EnsureDirectory(path string) error {
	exists, _, err := DirectoryExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", path)
	}

	if !exists {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed create directory %s", path)
		}
	}

	return nil
}
