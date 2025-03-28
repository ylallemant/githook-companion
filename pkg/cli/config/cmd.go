package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/config/arch"
	"github.com/ylallemant/githook-companion/pkg/cli/config/dependency"
	"github.com/ylallemant/githook-companion/pkg/cli/config/directory"
	"github.com/ylallemant/githook-companion/pkg/cli/config/hook"
	"github.com/ylallemant/githook-companion/pkg/cli/config/os"
	"github.com/ylallemant/githook-companion/pkg/cli/config/show"
	"github.com/ylallemant/githook-companion/pkg/cli/config/sync"
)

var rootCmd = &cobra.Command{
	Use:   "config",
	Short: "tools for config processing",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(directory.Command())
	rootCmd.AddCommand(arch.Command())
	rootCmd.AddCommand(os.Command())
	rootCmd.AddCommand(show.Command())
	rootCmd.AddCommand(dependency.Command())
	rootCmd.AddCommand(hook.Command())
	rootCmd.AddCommand(sync.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
