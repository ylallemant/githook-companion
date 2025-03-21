package config

import (
	"path/filepath"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
)

func DirectoryPathFromBase(path string) string {
	return filepath.Join(path, api.ConfigDirectory)
}

func FilePathFromBase(path string) string {
	return filepath.Join(path, api.ConfigDirectory, api.ConfigFile)
}

func GetLocalBasePath() (string, error) {
	local, err := environment.CurrentDirectory()
	if err != nil {
		return "", err
	}

	return local, nil
}

func GetLocalFilePath() (string, error) {
	local, err := GetLocalBasePath()
	if err != nil {
		return "", err
	}

	path := FilePathFromBase(local)

	return path, nil
}

func GetGlobalBasePath() (string, error) {
	home, err := environment.Home()
	if err != nil {
		return "", err
	}

	return home, nil
}

func GetGlobalFilePath() (string, error) {
	home, err := GetGlobalBasePath()
	if err != nil {
		return "", err
	}

	path := FilePathFromBase(home)

	return path, nil
}
