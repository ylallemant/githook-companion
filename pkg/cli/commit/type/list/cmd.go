package list

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/cli/commit/validate"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

var rootCmd = &cobra.Command{
	Use:   "list",
	Short: "list commit message types",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		configuration, err := config.Get(globals.Current.ConfigPath)
		if err != nil {
			return err
		}

		t := table.NewWriter()
		t.SetStyle(table.StyleLight)
		t.SetOutputMirror(os.Stdout)

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
	rootCmd.AddCommand(validate.Command())
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
