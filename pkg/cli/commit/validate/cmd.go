package validate

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githooks-butler/pkg/cli/commit/validate/options"
	"github.com/ylallemant/githooks-butler/pkg/config"
	"github.com/ylallemant/githooks-butler/pkg/environment"
	"github.com/ylallemant/githooks-butler/pkg/git/commit"
	"github.com/ylallemant/githooks-butler/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "validate",
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

		message := options.Current.Message

		validated, commitType, dictionary := commit.Validate(message, configuration)

		if dictionary != nil {
			// commit type found through dictionary match on first word
			// ensure that the dictionary value is used in the message
			message = commit.EnsureDictionaryValue(message, dictionary)
		}

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
			commitType = configuration.Commit.Types[idx].Type
		} else if !validated {
			// binary has not been called from a terminal
			// no user interaction possible
			// output invalidity information and throw error
			fmt.Println("-- error cause --------------------------------------")
			fmt.Println("commit message malformed : add a commit type prefix")
			fmt.Println("format: \"<commit-type>: <commit-message>\"")
			fmt.Println("available commit types:")
			for _, commitType := range configuration.Commit.Types {
				fmt.Println("  - ", commitType.Type, ": ", commitType.Description)
			}
			fmt.Println("-- error --------------------------------------------")

			return errors.New("message is missing commit type - see \"error cause\" block for more information")
		}

		// ensure commit type prefix format (lower-case)
		message = commit.EnsureFormat(message, commitType)

		if options.Current.OutputFilePath == "" {
			// output to terminal
			fmt.Println(message)
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
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
