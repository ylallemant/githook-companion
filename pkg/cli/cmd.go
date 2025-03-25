package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ylallemant/githook-companion/pkg/cli/config"
	"github.com/ylallemant/githook-companion/pkg/cli/debug"
	"github.com/ylallemant/githook-companion/pkg/cli/dependency"
	"github.com/ylallemant/githook-companion/pkg/cli/git"
	initCmd "github.com/ylallemant/githook-companion/pkg/cli/init"
	"github.com/ylallemant/githook-companion/pkg/cli/remove"
	"github.com/ylallemant/githook-companion/pkg/cli/update"
	"github.com/ylallemant/githook-companion/pkg/cli/version"
)

var rootCmd = &cobra.Command{
	Use:   "githook-companion",
	Short: "githook-companion provides a toolset facilitating complex git-hook workflows",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dependency.Command())
	rootCmd.AddCommand(git.Command())
	rootCmd.AddCommand(debug.Command())
	rootCmd.AddCommand(initCmd.Command())
	rootCmd.AddCommand(remove.Command())
	rootCmd.AddCommand(config.Command())
	rootCmd.AddCommand(update.Command())
	rootCmd.AddCommand(version.Command())
}

func Command() *cobra.Command {
	return rootCmd
}
