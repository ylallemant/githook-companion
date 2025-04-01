package lock

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/lock/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "lock \"hookname\"",
	Short: "locks a specific hook permanently or for a specific duration",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		hookname := args[0]

		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		hookList, exist, err := config.ListGithooks(configContext.Config())
		if err != nil {
			return err
		}

		if exist {
			if slices.Contains(hookList, hookname) {
				if options.Current.Duration > 0 {
					err = config.SetTimedLock(hookname, options.Current.Duration, configContext)
				} else {
					err = config.SetPermanentLock(hookname, configContext)
				}
			} else {
				err = errors.Errorf("githook \"%s\" is unknown. possible values are %v", hookname, hookList)
			}
		} else {
			err = errors.Errorf("no githooks are defined")
		}

		fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("locked githook \"%s\"", hookname))
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().DurationVarP(&options.Current.Duration, "duration", "d", options.Current.Duration, "set a validity duration for the lock. without this setting the lock is permanent")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
