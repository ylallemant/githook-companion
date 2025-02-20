package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/version/options"
	"github.com/ylallemant/githook-companion/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:   "version",
	Short: "outputs the version of the binary",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		if options.Current.Semver {
			fmt.Println(version.Semver())
			return nil
		}

		if options.Current.Commit {
			fmt.Println(version.Commit())
			return nil
		}

		fmt.Println(version.GetInfo())

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Commit, "commit", options.Current.Commit, "print only the commit hash")
	rootCmd.PersistentFlags().BoolVar(&options.Current.Semver, "semver", options.Current.Semver, "print only the semver string")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
