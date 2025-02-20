package dependencies

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/install/dependencies/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "dependencies",
	Short: "install all dependencies",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Get(globals.Current.ConfigPath)
		if err != nil {
			return err
		}

		installationDirectory := options.Current.Directory

		if installationDirectory == "" {
			home, err := environment.Home()
			if err != nil {
				return err
			}

			err = environment.EnsureDirectory(filepath.Join(home, ".local"))
			if err != nil {
				return err
			}

			installationDirectory = filepath.Join(home, ".local", "bin")
		}

		err = environment.EnsureDirectory(installationDirectory)
		if err != nil {
			return err
		}

		fmt.Println("install all dependencies in directory", installationDirectory)

		for _, tool := range cfg.Dependencies {
			fmt.Println("check", tool.Name, dependency.Version(tool))

			installed, err := dependency.Available(tool, installationDirectory)
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
					availableVersion, err := dependency.AvailableVersion(filepath.Join(installationDirectory, tool.Name))
					if err != nil {
						return errors.Wrap(err, "failed to get available version")
					}

					if strings.Contains(availableVersion, dependency.Version(tool)) {
						fmt.Println("  - is already installed in the right version")
					} else {
						fmt.Println("  - is available in an wanted version", availableVersion)
						deleteOldBinary = true
					}
				}

				if deleteOldBinary {
					fmt.Println("  - delete available binary")
					dependency.Delete(tool, installationDirectory)
					if err != nil {
						return err
					}

					needsInstallation = true
				}
			}
			fmt.Println("  - installation needed", needsInstallation)

			if needsInstallation {
				fmt.Println("  - installing", tool.Name)
				err = dependency.Install(tool, installationDirectory)
				if err != nil {
					return errors.Wrap(err, "failed to install all dependencies")
				}
			}
		}

		fmt.Println("all dependencies installed")
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Current.Directory, "directory", "d", options.Current.Directory, "installation directory")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
