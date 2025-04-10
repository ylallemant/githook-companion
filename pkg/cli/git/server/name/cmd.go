package name

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/server/name/options"
	"github.com/ylallemant/githook-companion/pkg/git"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "name",
	Short: "returns a provider name based on the hostname configured in \"remote.origin.url\", defaults to hostname",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		globals.ProcessGlobals()

		name, err := git.Name(options.Current.Default)
		if err != nil {
			return err
		}

		fmt.Println(name)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&options.Current.Default, "default", options.Current.Default, "returned value if origin hostname can't be resolved to a name")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
