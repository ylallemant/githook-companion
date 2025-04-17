package terminal

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/environment"
)

var rootCmd = &cobra.Command{
	Use:   "terminal",
	Short: "returns whether or not the command runs in a terminal with \"true\" or \"false\"",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		isTerminal, err := environment.CalledFromTerminal()
		if err != nil {
			return errors.Wrap(err, "failed to check if command runs in a terminal")
		}

		if isTerminal {
			fmt.Fprintln(cmd.OutOrStdout(), "true")
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "false")
		return nil
	},
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
