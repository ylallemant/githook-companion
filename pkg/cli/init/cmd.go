package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/cli/init/options"
	"github.com/ylallemant/githook-companion/pkg/command"
	"github.com/ylallemant/githook-companion/pkg/config"
)

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize configuration locally or globally",
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

		err = config.EnsureConfiguration(path)

		configurationFile := filepath.Join(path, api.ConfigDirectory, api.ConfigFile)

		cfg, err := config.Load(configurationFile, true)
		if err != nil {
			return err
		}

		if cfg.ConfigReference != nil {
			// check and handle configuration reference
			err = config.EnsureReference(cfg.ConfigReference)
			if err != nil {
				return err
			}

			// set path to the right directory
			path = filepath.Join(path, cfg.ConfigReference.Path)
		}

		// check for the existance of the hooks directory
		hooksDirectory := filepath.Join(path, api.GithooksDirectory)
		exists, _, err := config.DirectoryExists(hooksDirectory)
		if err != nil {
			return err
		}

		if exists {
			// set githook property in local git configuration
			fmt.Printf("git config core.hooksPath %s\n", hooksDirectory)
			git := command.New("git")
			git.AddArg("config")
			git.AddArg("core.hooksPath")
			git.AddArg(hooksDirectory)

			_, err = git.Execute()
			if err != nil {
				return errors.Wrapf(err, "failed to set core.hooksPath git config to %s", hooksDirectory)
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
