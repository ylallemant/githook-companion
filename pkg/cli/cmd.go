package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ylallemant/githooks-butler/pkg/cli/commit"
	"github.com/ylallemant/githooks-butler/pkg/cli/config"
	"github.com/ylallemant/githooks-butler/pkg/cli/install"
	"github.com/ylallemant/githooks-butler/pkg/cli/server"
	"github.com/ylallemant/githooks-butler/pkg/cli/update"
	"github.com/ylallemant/githooks-butler/pkg/cli/version"
)

var rootCmd = &cobra.Command{
	Use:   "githooks-butler",
	Short: "githooks-butler provides a toolset facilitating complex git-hook workflows",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(config.Command())
	rootCmd.AddCommand(install.Command())
	rootCmd.AddCommand(commit.Command())
	rootCmd.AddCommand(server.Command())
	rootCmd.AddCommand(update.Command())
	rootCmd.AddCommand(version.Command())
}

func Command() *cobra.Command {
	return rootCmd
}
