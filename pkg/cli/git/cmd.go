package git

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/commit"
	"github.com/ylallemant/githook-companion/pkg/cli/git/server"
)

var rootCmd = &cobra.Command{
	Use:   "git",
	Short: "helpers for Git configuration",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(server.Command())
	rootCmd.AddCommand(commit.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
