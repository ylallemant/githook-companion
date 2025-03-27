package hostname

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/git"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "hostname",
	Short: "returns the hostname configured in \"remote.origin.url\"",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, err := git.Hostname()
		if err != nil {
			return err
		}

		fmt.Println(hostname)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
