package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ylallemant/githook-companion/pkg/cli/binary/upgrade"
	"github.com/ylallemant/githook-companion/pkg/cli/binary/version"
	"github.com/ylallemant/githook-companion/pkg/cli/config"
	"github.com/ylallemant/githook-companion/pkg/cli/debug"
	"github.com/ylallemant/githook-companion/pkg/cli/dependency"
	"github.com/ylallemant/githook-companion/pkg/cli/environment"
	"github.com/ylallemant/githook-companion/pkg/cli/git"
	initCmd "github.com/ylallemant/githook-companion/pkg/cli/init"
	"github.com/ylallemant/githook-companion/pkg/cli/remove"
)

var rootCmd = &cobra.Command{
	Use:          "githook-companion",
	Short:        "githook-companion provides a toolset facilitating complex git-hook workflows",
	SilenceUsage: true,
	Long:         ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(environment.Command())
	rootCmd.AddCommand(dependency.Command())
	rootCmd.AddCommand(git.Command())
	rootCmd.AddCommand(debug.Command())
	rootCmd.AddCommand(initCmd.Command())
	rootCmd.AddCommand(remove.Command())
	rootCmd.AddCommand(config.Command())
	rootCmd.AddCommand(upgrade.Command())
	rootCmd.AddCommand(version.Command())
}

func Command() *cobra.Command {
	return rootCmd
}
