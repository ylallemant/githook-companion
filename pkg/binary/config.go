package binary

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var ConfigDirectory string

func init() {
	var err error
	ConfigDirectory, err = configDirectory()
	if err != nil {
		panic(fmt.Sprintf("failed init: %s", err.Error()))
	}
}

func configDirectory() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get home directory")
	}

	executable, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get binary name")
	}
	executable = filepath.Base(executable)

	dirname = filepath.Join(dirname, fmt.Sprintf(".%s", executable))

	exists, err := directoryExists(dirname)
	if err != nil {
		return "", errors.Wrapf(err, "failed to check configuration directory")
	}

	if !exists {
		err = os.MkdirAll(dirname, 0755)
		if err != nil {
			return "", errors.Wrapf(err, "failed create configuration directory %s", dirname)
		}
	}

	return dirname, nil
}

func directoryExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}

		return false, errors.Wrapf(err, "failed to read path %s", path)
	}

	if fi.IsDir() {
		return true, nil
	}

	return true, errors.Wrapf(err, "path is no directory %s", path)
}
