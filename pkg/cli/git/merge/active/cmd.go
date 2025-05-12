package active

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/git/merge"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "active",
	Short: "checks if a merge is in progress",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// configContext, err := config.InitContext()
		// if err != nil {
		// 	return err
		// }

		active, err := merge.Active()
		if err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), active)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	//rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
