package directory

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/cli/config/directory/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "directory",
	Short: "outputs config directory path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		if options.Current.Child {
			// always return the local directory as "child"
			basePath, err := config.GetLocalBasePath()
			if err != nil && !globals.Current.FallbackConfig {
				return err
			}

			configDirectory := config.DirectoryPathFromBase(basePath)

			fmt.Fprintln(cmd.OutOrStdout(), configDirectory)
			return nil
		}

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

		installationDirectory := config.ParentPathFromConfig(configuration)
		configDirectory := config.DirectoryPathFromBase(installationDirectory)

		fmt.Fprintln(cmd.OutOrStdout(), configDirectory)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Child, "child", options.Current.Child, "forces the path to be relative to the current project")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
