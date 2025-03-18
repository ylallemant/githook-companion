package show

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "show",
	Short: "show config",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg := config.Default()

		if globals.Current.ConfigPath != "" {
			cfg, err = config.Load(globals.Current.ConfigPath, true)
			if err != nil {
				return err
			}
		} else {
			cfg, err = config.Get()
			if err != nil {
				return err
			}
		}

		yaml, err := config.ToYAML(cfg)
		if err != nil {
			return err
		}

		fmt.Println(string(yaml))
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
