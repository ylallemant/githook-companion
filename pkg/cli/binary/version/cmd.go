package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/binary"
	"github.com/ylallemant/githook-companion/pkg/cli/binary/version/options"
)

var rootCmd = &cobra.Command{
	Use:   "version",
	Short: "outputs the version of the binary",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		if options.Current.Semver {
			fmt.Println(binary.Semver())
			return nil
		}

		if options.Current.Commit {
			fmt.Println(binary.Commit())
			return nil
		}

		if options.Current.Separator != "" {
			fmt.Println(binary.SemverWithSeparator(options.Current.Separator))
			return nil
		}

		fmt.Println(binary.GetInfo())

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.Commit, "commit", options.Current.Commit, "print only the commit hash")
	rootCmd.PersistentFlags().BoolVar(&options.Current.Semver, "semver", options.Current.Semver, "print only the semver string")
	rootCmd.PersistentFlags().StringVarP(&options.Current.Separator, "separator", "s", options.Current.Separator, "replace the point in the semver notation")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
