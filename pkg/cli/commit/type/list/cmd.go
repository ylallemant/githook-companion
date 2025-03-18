package list

import (
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
		var err error
		configuration := config.Default()

		if globals.Current.ConfigPath != "" {
			configuration, err = config.Load(globals.Current.ConfigPath, true)
			if err != nil {
				return err
			}
		} else {
			configuration, err = config.Get()
			if err != nil {
				return err
			}
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(cmd.OutOrStdout())

		t.SetTitle("Commit Types")
		t.AppendHeader(table.Row{"Type", "Description"})

		for _, commitType := range configuration.Commit.Types {
			t.AppendRow([]interface{}{
				commitType.Type,
				commitType.Description,
			})
		}

		t.Render()

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
