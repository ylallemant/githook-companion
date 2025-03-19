package dependency

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/config/dependency/path"
)

var rootCmd = &cobra.Command{
	Use:   "dependency",
	Short: "dependency configuration information",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("please use a subcommand...")
		cmd.Usage()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(path.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
