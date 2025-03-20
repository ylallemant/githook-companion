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
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	gitConfig "github.com/ylallemant/githook-companion/pkg/git/config"
)

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize configuration locally or globally",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var path string
		var reference *api.ConfigReference

		if options.Current.Global {
			path, err = config.GetGlobalBasePath()
		} else {
			path, err = config.GetLocalBasePath()
		}

		if options.Current.ReferencePath != "" && options.Current.ReferenceRepository != "" {
			reference = &api.ConfigReference{
				GitRepository: options.Current.ReferenceRepository,
				Path:          options.Current.ReferencePath,
			}
		}

		if err != nil {
			return errors.Wrap(err, "failed to get base path")
		}

		err = config.EnsureConfiguration(path, reference, options.Current.Minimalistic)

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
			err = gitConfig.SetProperty("core.hooksPath", hooksDirectory, options.Current.Global)
			if err != nil {
				return err
			}
		}

		if len(cfg.Dependencies) > 0 {
			directory := dependency.InstallDirectoryFromConfig(cfg)
			err = dependency.InstallAll(directory, cfg)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Global, "global", options.Current.Global, "make a global initialization")
	rootCmd.PersistentFlags().StringVar(&options.Current.ReferenceRepository, "reference-repository", options.Current.ReferenceRepository, "repository of the centralised configuration")
	rootCmd.PersistentFlags().StringVar(&options.Current.ReferencePath, "reference-path", options.Current.ReferencePath, fmt.Sprintf("relative path to the centralised configuration root (parent of %s)", api.ConfigDirectory))
	rootCmd.PersistentFlags().BoolVarP(&options.Current.Minimalistic, "minimalistic", "m", options.Current.Minimalistic, "only install the bare minimum. no hooks, no dictionaries, no nothing")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
