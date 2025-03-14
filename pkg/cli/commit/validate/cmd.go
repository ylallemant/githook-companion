package validate

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/commit/validate/options"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/git/commit"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "validate [message]",
	Short: "validates interactively Git commit messages",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		// assess if binary was called in a terminal or from some editor/git client
		calledFromTerminal, err := environment.CalledFromTerminal()
		if err != nil {
			return errors.Wrap(err, "failed to assess if called from terminal")
		}

		configuration, err := config.Get(globals.Current.ConfigPath)
		if err != nil {
			return err
		}

		if !environment.IsAnArgument(os.Args[3]) && options.Current.Message != "" {
			return errors.Errorf("too many messages provided, choose whether per argument or flag")
		}

		message := os.Args[3]

		if options.Current.Message != "" {
			message = options.Current.Message
		}

		if message == "" {
			return errors.Errorf("providing a commit message by argument or flag is mendatory")
		}

		languageCode, validated, commitTypeToken, tokens := commit.Validate(message, configuration)

		if !validated && calledFromTerminal {
			// message does not have a commit type prefix
			// and no commit type could be correlated through dictionaries
			// request user input through interactive commit type selection
			// this can only work in a terminal
			templates := &promptui.SelectTemplates{
				Active:   fmt.Sprintf("%s {{ .Type | underline }} : {{ .Description | underline }}", promptui.IconSelect),
				Inactive: "  {{ .Type }} : {{ .Description }}",
				Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Type | faint }}`, promptui.IconGood),
			}

			prompt := promptui.Select{
				Label:     "Select Commit Type",
				Items:     configuration.Commit.Types,
				Templates: templates,
			}

			idx, _, err := prompt.Run()
			if err != nil {
				return errors.Wrap(err, "interactive commit type user selection failed")
			}

			// commit type from user input
			commitTypeToken = commit.CommitTypeTokenFromString(configuration.Commit.Types[idx].Type, languageCode)
		} else if !validated {
			// binary has not been called from a terminal
			// no user interaction possible
			// output invalidity information and throw error
			typeList := ""
			for _, commitType := range configuration.Commit.Types {
				typeList = typeList + fmt.Sprintf("    - %s: %s\n", commitType.Type, commitType.Description)
			}

			nonInteractiveErrorMessage := fmt.Sprintf(`commit message malformed
  you didn't commit on the command line, commit type can not be added interactively
  please make sure to provide a commit type prefix in your message
  format: "<commit-type-prefix>: <commit-message>"
  available commit types:
%s
			`, typeList)

			//fmt.Fprintln(cmd.OutOrStderr(), nonInteractiveErrorMessage)
			//os.Exit(1)
			return errors.New(nonInteractiveErrorMessage)
		}

		// ensure commit type prefix format (lower-case)
		message, err = commit.EnsureFormat(message, configuration.Commit.MessageTemplate, commitTypeToken, tokens)
		if err != nil {
			return errors.Wrap(err, "failed to format commit message")
		}

		if options.Current.OutputFilePath == "" {
			// output to terminal
			fmt.Fprintln(cmd.OutOrStdout(), message)
		} else {
			// output to file
			file, err := os.OpenFile(options.Current.OutputFilePath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to write to output file %s", options.Current.OutputFilePath)
			}

			defer file.Close()

			file.WriteString(message)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&options.Current.Message, "message", "m", options.Current.Message, "commit message")
	rootCmd.PersistentFlags().StringVarP(&options.Current.OutputFilePath, "output", "o", options.Current.OutputFilePath, "output file path")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
