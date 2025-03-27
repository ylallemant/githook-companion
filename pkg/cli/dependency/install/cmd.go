package install

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/dependency/install/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "install",
	Short: "install dependency",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		installationDirectory := dependency.DependencyDirectoryFromConfig(configContext.Config())

		if options.Current.Directory != "" {
			installationDirectory = options.Current.Directory
		}

		return dependency.InstallAll(installationDirectory, configContext.Config())
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Current.Directory, "directory", "d", options.Current.Directory, "installation directory")
	rootCmd.PersistentFlags().BoolVarP(&options.Current.All, "all", "A", options.Current.All, "install all dependencies")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
