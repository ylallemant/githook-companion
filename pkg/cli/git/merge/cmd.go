package merge

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/merge/active"
)

var rootCmd = &cobra.Command{
	Use:   "merge",
	Short: "tooling for git merges",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(active.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
