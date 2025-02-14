package validate

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githooks-butler/pkg/cli/commit/validate/options"
	"github.com/ylallemant/githooks-butler/pkg/config"
	"github.com/ylallemant/githooks-butler/pkg/git/commit"
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

		message, validated := commit.Validate(options.Current.Message, configuration)

		if !validated {
			prompt := promptui.Select{
				Label: "Select Commit Type",
				Items: config.CommitTypes(configuration),
			}

			_, result, err := prompt.Run()

			if err != nil {
				return err
			}

			message = fmt.Sprintf("%s: %s", result, message)
		}

		fmt.Println(message)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Current.Message, "message", "m", options.Current.Message, "commit message")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
