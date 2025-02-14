package validate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githooks-butler/pkg/config"
	"github.com/ylallemant/githooks-butler/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "validate",
	Short: "validates interactively Git commit messages",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.Get(globals.Current.ConfigPath)
		if err != nil {
			return err
		}

		pretty, _ := config.ToPrettyJSON(configuration)

		fmt.Println(string(pretty))

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "current git branch")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
