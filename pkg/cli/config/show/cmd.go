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
		cfg, err := config.Get(globals.Current.ConfigPath)
		if err != nil {
			return err
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
