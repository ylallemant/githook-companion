package show

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var configuration *api.Config
		var err error

		if globals.Current.ConfigPath != "" {
			configuration, err = config.Load(globals.Current.ConfigPath, true)
			if err != nil && !globals.Current.FallbackConfig {
				return err
			}
		} else {
			configContext, err := config.InitContext()
			if err != nil && !globals.Current.FallbackConfig {
				return err
			}

			if configContext != nil {
				configuration = configContext.Config()
			}
		}

		if configuration == nil && globals.Current.FallbackConfig {
			configuration = config.Default()
		}

		yaml, err := config.ToYAML(configuration)
		if err != nil {
			return err
		}

		fmt.Println(string(yaml))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
