package commit

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	commitType "github.com/ylallemant/githook-companion/pkg/cli/commit/type"
	"github.com/ylallemant/githook-companion/pkg/cli/commit/validate"
)

var rootCmd = &cobra.Command{
	Use:   "commit",
	Short: "tools for commit message processing",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validate.Command())
	rootCmd.AddCommand(commitType.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
