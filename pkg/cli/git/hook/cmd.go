package hook

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/disable"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/lock"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/locked"
	"github.com/ylallemant/githook-companion/pkg/cli/git/hook/unlock"
)

var rootCmd = &cobra.Command{
	Use:   "hook",
	Short: "tools for githook management",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(disable.Command())
	rootCmd.AddCommand(unlock.Command())
	rootCmd.AddCommand(locked.Command())
	rootCmd.AddCommand(lock.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
