package testutils

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ResetCobraFlags gets rid of the persistence of flag values inside cobra.
// This persistence may break tests that you define for cobra commands.
// this problem is being discussed in this issue : https://github.com/spf13/cobra/issues/2079
func ResetCobraFlags(cmd *cobra.Command) {
	if cmd.Flags().Parsed() {
		cmd.Flags().Visit(func(pf *pflag.Flag) {
			if err := pf.Value.Set(pf.DefValue); err != nil {
				panic(fmt.Errorf("reset argument[%s] value error %v", pf.Name, err))
			}
		})
	}
}
