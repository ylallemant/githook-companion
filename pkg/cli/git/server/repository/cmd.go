package repository

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/git"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "repository",
	Short: "returns the repository https url extracted \"remote.origin.url\"",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, err := git.Repository()
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
