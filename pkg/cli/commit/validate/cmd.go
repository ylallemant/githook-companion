package validate

import (
	"fmt"
	"os"
	"slices"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
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
		globals.ProcessGlobals()

		// assess if binary was called in a terminal or from some editor/git client
		calledFromTerminal, err := environment.CalledFromTerminal()
		if err != nil {
			return errors.Wrap(err, "failed to assess if called from terminal")
		}

		var configuration *api.Config

		if globals.Current.ConfigPath != "" {
			configuration, err = config.Load(globals.Current.ConfigPath, true)
		} else {
			configuration, err = config.Get()
		}

		if err != nil && !globals.Current.FallbackConfig {
			return err
		}

		if configuration == nil {
			configuration = config.Default()
		}

		if !environment.IsAnArgument(os.Args[3]) && options.Current.Message != "" {
			return errors.Errorf("too many messages provided, choose whether per argument or flag")
		}

		message := os.Args[3]
		log.Debug().Msgf("message from arguments \"%s\"", message)

		if options.Current.Message != "" {
			message = options.Current.Message
			log.Debug().Msgf("message from flags \"%s\"", message)
		}

		if message == "" {
			return errors.Errorf("providing a commit message by argument or flag is mendatory")
		}

		log.Debug().Msgf("validate \"%s\"", message)

		languageCode, validated, commitTypeToken, tokens, err := commit.Validate(message, configuration)
		if err != nil {
			return errors.Wrap(err, "failed validation")
		}

		log.Debug().Msgf("message validated \"%v\"", validated)
		log.Debug().Msgf("command called from terminal \"%v\"", calledFromTerminal)

		if !validated && calledFromTerminal {
			log.Debug().Msg("invalid message, user input required")
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

			commiType := configuration.Commit.Types[idx].Type
			log.Debug().Msgf("user selected commit type number %d \"%s\"", idx, commiType)

			// commit type from user input
			commitTypeToken = commit.CommitTypeTokenFromString(commiType, languageCode)
		} else if !validated {
			log.Debug().Msg("invalid message error because no user input is possible")
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

		log.Debug().Msgf("commit type token: %s", commitTypeToken.Value)
		log.Debug().Msgf("commit types without formatting: %v", configuration.Commit.NoFormatting)

		if !slices.Contains(configuration.Commit.NoFormatting, commitTypeToken.Value) {
			// ensure commit type prefix format (lower-case)
			message, err = commit.EnsureFormat(message, configuration.Commit.MessageTemplate, commitTypeToken, tokens)
			if err != nil {
				return errors.Wrap(err, "failed to format commit message")
			}
		}

		if options.Current.OutputFilePath == "" {
			log.Debug().Msg("output to terminal")
			// output to terminal
			fmt.Fprintln(cmd.OutOrStdout(), message)
		} else {
			log.Debug().Msg("output to file")
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
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
	rootCmd.PersistentFlags().BoolVar(&globals.Current.Debug, "debug", globals.Current.Debug, "outputs processing information")
	rootCmd.SetOutput(os.Stderr)
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
