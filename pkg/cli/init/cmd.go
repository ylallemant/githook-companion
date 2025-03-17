package init

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/cli/init/options"
	"github.com/ylallemant/githook-companion/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize configuration locally or globally",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *api.Config
		var err error

		if options.Current.Global {
			cfg, err = config.GetGlobally()
		} else {
			cfg, err = config.GetLocally()
		}

		if err != nil && !errors.Is(err, api.ConfigurationNotFound) {
			return err
		}

		if cfg == nil {
			fmt.Println("no config found", err.Error())
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Global, "global", options.Current.Global, "make a global initialization")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
