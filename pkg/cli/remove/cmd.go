package remove

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/cli/init/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	gitConfig "github.com/ylallemant/githook-companion/pkg/git/config"
)

var rootCmd = &cobra.Command{
	Use:   "remove",
	Short: "removes the configuration locally or globally",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var path string

		if options.Current.Global {
			path, err = config.GetGlobalBasePath()
		} else {
			path, err = config.GetLocalBasePath()
		}

		if err != nil {
			return errors.Wrap(err, "failed to get base path")
		}

		configurationDirectory := filepath.Join(path, api.ConfigDirectory)
		exists, _, err := config.DirectoryExists(configurationDirectory)
		if err != nil {
			return err
		}

		if exists {
			err = config.Remove(configurationDirectory)
			if err != nil {
				return err
			}
		}

		exists, err = gitConfig.PropertyExists("core.hooksPath", options.Current.Global)
		if err != nil {
			return err
		}

		if exists {
			err = gitConfig.UnsetProperty("core.hooksPath", options.Current.Global)
			if err != nil {
				return err
			}
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
