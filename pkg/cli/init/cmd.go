package init

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/cli/init/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
	"github.com/ylallemant/githook-companion/pkg/git"
	gitConfig "github.com/ylallemant/githook-companion/pkg/git/config"
	"github.com/ylallemant/githook-companion/pkg/git/hook"
)

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize configuration locally or globally",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var basePath string
		var reference *api.ParentConfig

		globalBasePath, err := config.GetGlobalBasePath()
		if err != nil {
			return err
		}

		localBasePath, err := config.GetLocalBasePath()
		if err != nil {
			return err
		}

		if options.Current.Global {
			basePath = globalBasePath
		} else {
			basePath = localBasePath
		}

		if options.Current.ParentPath != "" && options.Current.ParentRepository != "" {
			reference = &api.ParentConfig{
				GitRepository: options.Current.ParentRepository,
				Path:          options.Current.ParentPath,
				Private:       options.Current.ParentPrivate,
			}
		}

		if err != nil {
			return errors.Wrap(err, "failed to get base path")
		}

		err = config.EnsureConfiguration(basePath, reference, options.Current.Minimalistic)
		if err != nil {
			return errors.Wrap(err, "failed to ensure configuration")
		}

		configContext, err := config.Context(false)
		if err != nil {
			return err
		}

		if configContext.Config().ParentConfig != nil {
			// check and handle configuration reference
			err = config.EnsureReference(configContext.Config().ParentConfig)
			if err != nil {
				return err
			}
		}

		parentBasePath, err := config.BasePathFromConfig(configContext.Config())
		if err != nil {
			return err
		}

		// ensure githooks are present
		err = hook.Ensure(configContext.Config())
		if err != nil {
			return err
		}

		// check for the existance of the hooks directory
		hooksDirectory := config.GithooksPathFromConfig(configContext.Config())
		exists, _, err := filesystem.DirectoryExists(hooksDirectory)
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

		if len(configContext.Config().Dependencies) > 0 {
			directory := dependency.DependencyDirectoryFromConfig(configContext.Config())
			err = dependency.InstallAll(directory, configContext.Config())
			if err != nil {
				return err
			}
		}

		if parentBasePath != basePath {
			// ensure rules exist in the parent .gitignore file
			err = git.EnsureGitIgnoreFromBasePath(parentBasePath)
			if err != nil {
				return err
			}
		}

		// ensure rules exist in the child/project .gitignore file
		err = git.EnsureGitIgnoreFromBasePath(localBasePath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Global, "global", options.Current.Global, "make a global initialization")
	rootCmd.PersistentFlags().StringVar(&options.Current.ParentRepository, "parent-repository", options.Current.ParentRepository, "repository of the parent configuration. will be automatically checked out if necessary")
	rootCmd.PersistentFlags().StringVar(&options.Current.ParentPath, "parent-path", options.Current.ParentPath, fmt.Sprintf("relative path to the parent configuration root (parent of %s)", api.ConfigDirectory))
	rootCmd.PersistentFlags().BoolVar(&options.Current.ParentPrivate, "parent-private", options.Current.ParentPrivate, "specifies if the parent configuration repository is private")
	rootCmd.PersistentFlags().BoolVarP(&options.Current.Minimalistic, "minimalistic", "m", options.Current.Minimalistic, "only install the bare minimum. no hooks, no dictionaries, no nothing")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
