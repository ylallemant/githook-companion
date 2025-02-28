package repository

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/git/server"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "repository",
	Short: "returns the repository https url extracted \"remote.origin.url\"",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, err := server.Repository()
		if err != nil {
			return err
		}

		fmt.Println(hostname)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
