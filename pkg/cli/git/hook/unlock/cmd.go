package unlock

import (
	"fmt"
	"os"
	"slices"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/unlock/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "unlock [\"hookname\"]",
	Short: "unlocks a specific hook or all hooks",
	Long:  ``,
	Args:  cobra.MatchAll(cobra.MinimumNArgs(0), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		var hookname string

		if len(args) > 0 {
			hookname = args[0]
		}

		if (hookname == "" && !options.Current.All) || (hookname != "" && options.Current.All) {
			return errors.Errorf("you have to specify wether the githook name argument or the --all flag")
		}

		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		hookList, exist, err := config.ListGithooks(configContext.Config())
		if err != nil {
			return err
		}

		if exist {
			unlockList := []string{hookname}

			if options.Current.All {
				unlockList = hookList
			}

			for _, githook := range unlockList {
				if slices.Contains(hookList, githook) {
					err = config.LockRemove(githook, configContext)
				} else {
					err = errors.Errorf("githook \"%s\" is unknown. possible values are %v", githook, hookList)
				}

				if err != nil {
					if os.IsNotExist(err) {
						continue
					} else {
						return err
					}
				}

				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("unlocked githook \"%s\"", githook))
			}

		} else {
			err = errors.Errorf("no githooks are defined")
		}

		return err
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&options.Current.All, "all", "a", options.Current.All, "unlock all githooks")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
