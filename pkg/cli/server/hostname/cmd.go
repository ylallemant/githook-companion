package hostname

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/git/server"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "hostname",
	Short: "returns the hostname configured in \"remote.origin.url\"",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostname, err := server.Hostname()
		if err != nil {
			return err
		}

		fmt.Println(hostname)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
