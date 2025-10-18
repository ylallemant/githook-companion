package sync

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/sync/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "sync",
	Short: "synchronize local configuration with remotes (configurations, dependencies)",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configContext, err := config.InitContextWithoutSync()
		if err != nil {
			return err
		}

		if options.Current.SkipParents {
			log.Info().Msg("skipping parent configuration synchronization")
		} else {
			_, err = config.EnsureVersionSync(configContext)
			if err != nil {
				return err
			}
		}

		if options.Current.SkipDependencies {
			log.Info().Msg("skipping dependency synchronization")
		} else {
			path := dependency.DependencyDirectoryFromConfig(configContext.Config())
			err = dependency.InstallAll(path, configContext.Config())
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.SkipParents, "skip-parents", options.Current.SkipParents, "skip parent configuration synchronization")
	rootCmd.PersistentFlags().BoolVar(&options.Current.SkipDependencies, "skip-dependencies", options.Current.SkipDependencies, "skip dependency synchronization")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
