package active

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/git/hook"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "active \"hookname\"",
	Short: "checks if a specific hook is active or not",
	Long: `checks if a specific hook is active or not:
	the status is affected by different things like a lock,
	an active merge process, ...`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		status := false
		hookname := args[0]

		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		status, err = hook.Active(hookname, configContext)
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), status)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
