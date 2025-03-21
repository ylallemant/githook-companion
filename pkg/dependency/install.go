package dependency

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/environment"
)

func InstallDirectoryFromConfig(configuration *api.Config) string {
	if configuration.DependencyDirectory != "" {
		path, err := environment.EnsureAbsolutePath(configuration.DependencyDirectory)
		if err != nil {
			panic(err)
		}

		return path
	}

	if configuration.ParentConfig != nil {
		if configuration.ParentConfig.Path != "" {
			relativePath := filepath.Join(configuration.ParentConfig.Path, "bin")
			path, err := environment.EnsureAbsolutePath(relativePath)
			if err != nil {
				panic(err)
			}

			return path
		}
	}

	return DefaultInstallDirectory()
}

func DefaultInstallDirectory() string {
	switch runtime.GOOS {
	case "darwin":
		return "/usr/local/bin"
	default:
		home, err := environment.Home()
		if err != nil {
			panic(err)
		}

		err = environment.EnsureDirectory(filepath.Join(home, ".local"))
		if err != nil {
			panic(err)
		}

		return filepath.Join(home, ".local", "bin")
	}
}

func InstallAll(installationDirectory string, configuration *api.Config) error {
	err := environment.EnsureDirectory(installationDirectory)
	if err != nil {
		return err
	}

	fmt.Println("install all dependencies in directory", installationDirectory)

	for _, tool := range configuration.Dependencies {
		fmt.Println("check", tool.Name, Version(tool))

		installed, err := Available(tool, installationDirectory)
		if err != nil {
			return errors.Wrap(err, "failed to install all dependencies")
		}

		needsInstallation := !installed

		if installed {
			deleteOldBinary := false

			if tool.ForceReplace {
				fmt.Println("  - force replace enabled")
				deleteOldBinary = true
			} else {
				availableVersion, err := AvailableVersion(filepath.Join(installationDirectory, tool.Name))
				if err != nil {
					return errors.Wrap(err, "failed to get available version")
				}

				if strings.Contains(availableVersion, Version(tool)) {
					fmt.Println("  - is already installed in the right version")
				} else {
					fmt.Println("  - is available in an wanted version", availableVersion)
					deleteOldBinary = true
				}
			}

			if deleteOldBinary {
				fmt.Println("  - delete available binary")
				Delete(tool, installationDirectory)
				if err != nil {
					return err
				}

				needsInstallation = true
			}
		}
		fmt.Println("  - installation needed", needsInstallation)

		if needsInstallation {
			fmt.Println("  - installing", tool.Name)
			err = Install(tool, installationDirectory)
			if err != nil {
				return errors.Wrap(err, "failed to install all dependencies")
			}
		}
	}

	fmt.Println("all dependencies installed")
	return nil

}
