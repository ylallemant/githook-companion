package locked

import (
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/lock/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "locked \"hookname\"",
	Short: "checks if a specific hook is permanently or temporary locked",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		status := false
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
				path := config.GithookLockPathFromNameAndConfig(hookname, configContext)
				exists, _, err := filesystem.FileExists(path)
				if err != nil {
					return err
				}

				lockType := filesystem.LockType(path)

				if exists {
					if lockType == filesystem.LockTypeTemporary {
						status, err = config.TimeLockActive(hookname, configContext)
					} else {
						status, err = config.PermanentLockExists(hookname, configContext)
					}

					if err != nil {
						return errors.Wrapf(err, "failed to read lock status for %s", path)
					}
				}
			} else {
				return errors.Errorf("githook \"%s\" is unknown. possible values are %v", hookname, hookList)
			}
		} else {
			return errors.Errorf("no githooks are defined")
		}

		fmt.Fprintln(cmd.OutOrStdout(), status)
		return nil
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
