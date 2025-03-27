package list

import (
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "list",
	Short: "list commit message types",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configContext, err := config.InitContext()
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(cmd.OutOrStdout())

		t.SetTitle("Commit Types")
		t.AppendHeader(table.Row{"Type", "Description", "Auto-Format"})

		for _, commitType := range configContext.Config().Commit.Types {
			t.AppendRow([]interface{}{
				commitType.Type,
				commitType.Description,
				!slices.Contains(configContext.Config().Commit.NoFormatting, commitType.Type),
			})
		}

		t.Render()

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
