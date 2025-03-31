package debug

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/binary"
	"github.com/ylallemant/githook-companion/pkg/cli/debug/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/git"
)

var rootCmd = &cobra.Command{
	Use:   "debug",
	Short: "displays key information about the environment",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("githook-companion version", binary.Information())
		fmt.Println(environment.Debug())
		fmt.Println(git.Debug())
		fmt.Println(config.Debug())

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Global, "global", options.Current.Global, "make a global initialization")
	rootCmd.PersistentFlags().StringVar(&options.Current.ReferenceRepository, "reference-repository", options.Current.ReferenceRepository, "repository of the centralised configuration")
	rootCmd.PersistentFlags().StringVar(&options.Current.ReferencePath, "reference-path", options.Current.ReferencePath, fmt.Sprintf("relative path to the centralised configuration root (parent of %s)", api.ConfigDirectory))
	rootCmd.PersistentFlags().BoolVarP(&options.Current.Minimalistic, "minimalistic", "m", options.Current.Minimalistic, "only install the bare minimum. no hooks, no dictionaries, no nothing")
	rootCmd.SetOut(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
