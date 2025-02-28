package server

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/server/hostname"
	"github.com/ylallemant/githook-companion/pkg/cli/server/name"
	"github.com/ylallemant/githook-companion/pkg/cli/server/repository"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "parses git server information",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hostname.Command())
	rootCmd.AddCommand(repository.Command())
	rootCmd.AddCommand(name.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
