package list

import (
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "list",
	Short: "list commit message types",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var configuration *api.Config

		fmt.Println("globals.Current.ConfigPath", globals.Current.ConfigPath)
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

		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(cmd.OutOrStdout())

		t.SetTitle("Commit Types")
		t.AppendHeader(table.Row{"Type", "Description", "Auto-Format"})

		for _, commitType := range configuration.Commit.Types {
			t.AppendRow([]interface{}{
				commitType.Type,
				commitType.Description,
				!slices.Contains(configuration.Commit.NoFormatting, commitType.Type),
			})
		}

		t.Render()

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&globals.Current.FallbackConfig, "fallback-config", globals.Current.FallbackConfig, "if no configuration was found, fallback to the default one")
	rootCmd.PersistentFlags().StringVarP(&globals.Current.ConfigPath, "config", "c", globals.Current.ConfigPath, "path to configuration file")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
