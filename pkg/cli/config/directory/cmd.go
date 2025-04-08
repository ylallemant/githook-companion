package directory

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/config/directory/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "directory",
	Short: "outputs config directory path",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		if options.Current.Child {
			configDirectory := config.DirectoryPathFromBase(configContext.LocalPath())

			fmt.Fprintln(cmd.OutOrStdout(), configDirectory)
			return nil
		}

		installationDirectory := config.ParentPathFromConfig(configContext.Config())

		if installationDirectory == "" {
			// no parent config exists
			// use local path
			installationDirectory = configContext.LocalPath()
		}
		configDirectory := config.DirectoryPathFromBase(installationDirectory)

		fmt.Fprintln(cmd.OutOrStdout(), configDirectory)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Child, "child", options.Current.Child, "forces the path to be relative to the current project")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
