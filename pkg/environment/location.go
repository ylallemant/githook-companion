package environment

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func Home() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get home directory")
	}

	return dirname, nil
}

func CurrentDirectory() (string, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get current directory")
	}
	return dirname, nil
}

func EnsureAbsolutePath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		var pathPrefix string
		var err error

		if strings.HasPrefix(path, "~") {
			pathPrefix, err = Home()
		} else {
			pathPrefix, err = CurrentDirectory()
		}

		if err != nil {
			return "", errors.Wrapf(err, "failed to ensure absolute path for %s", path)
		}

		path = filepath.Join(pathPrefix, path)
	}

	return path, nil
}

func EnsureDirectory(path string) error {
	err := os.Mkdir(path, os.FileMode(0755))
	if err == nil {
		return nil
	}

	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(path)
		if err != nil {
			return errors.Wrapf(err, "failed to check directory %s", path)
		}

		if !info.IsDir() {
			return errors.Errorf("path exists but is not a directory %s", path)
		}

		return nil
	}

	return errors.Wrapf(err, "failed to create directory %s", path)
}
