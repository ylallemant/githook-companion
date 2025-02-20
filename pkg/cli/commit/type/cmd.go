package committype

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/commit/type/list"
)

var rootCmd = &cobra.Command{
	Use:   "type",
	Short: "commit message type commands",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(list.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
