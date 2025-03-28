package disable

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/init/options"
	gitConfig "github.com/ylallemant/githook-companion/pkg/git/config"
)

var rootCmd = &cobra.Command{
	Use:   "disable",
	Short: "disable configuration locally or globally",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// set githook property in local git configuration
		err := gitConfig.UnsetProperty("core.hooksPath", options.Current.Global)
		if err != nil {
			return err
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
