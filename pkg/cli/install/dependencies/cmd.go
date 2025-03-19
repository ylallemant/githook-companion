package dependencies

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
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
		var err error
		var configuration *api.Config

		if globals.Current.ConfigPath != "" {
			configuration, err = config.Load(globals.Current.ConfigPath, true)
		} else {
			configuration, err = config.Get()
		}

		if err != nil && !globals.Current.FallbackConfig {
			return err
		}

		if configuration == nil {
			configuration = config.Default()
		}

		installationDirectory := options.Current.Directory

		err = environment.EnsureDirectory(installationDirectory)
		if err != nil {
			return err
		}

		fmt.Println("install all dependencies in directory", installationDirectory)

		for _, tool := range configuration.Dependencies {
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
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
