package upgrade

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ylallemant/githook-companion/pkg/binary"
	"github.com/ylallemant/githook-companion/pkg/cli/binary/upgrade/options"
)

var (
	owner      = "ylallemant"
	repo       = "githook-companion"
	binaryName = "githook-companion"
)

var rootCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade the binary",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentLocation, err := binary.Location()
		if err != nil {
			return errors.Wrapf(err, "failed to get current installation path")
		}
		fmt.Println("current binary location", currentLocation)
		tempDir, err := os.MkdirTemp(os.TempDir(), binaryName)
		if err != nil {
			return errors.Wrapf(err, "failed to create temporary directory")
		}
		defer os.RemoveAll(tempDir)
		fmt.Println("temp folder created at", tempDir)
		fmt.Println("targeting binary for", runtime.GOOS, runtime.GOARCH)

		releases, err := binary.ListReleases()
		if err != nil {
			return errors.Wrapf(err, "failed to list releases for repo %s/%s", owner, repo)
		}

		if len(releases) == 0 {
			return errors.Errorf("no release found for repo %s/%s", owner, repo)
		}

		// latest release
		wanted := binary.Latest(releases, options.Current.AllowPrerelease)

		currentVersion := binary.Semver()

		if wanted.GetTagName() == currentVersion && !options.Current.Force {
			fmt.Printf("binary with version \"%s\" is up to date : skipping upgrade\n", currentVersion)
			return nil
		}

		if options.Current.DryRun {
			fmt.Printf("upgrade would replace binary from \"%s\" to \"%s\" at its current location %s\n", currentVersion, wanted.GetTagName(), currentLocation)
			return nil
		} else {
			fmt.Printf("upgrade will replace binary from \"%s\" to \"%s\" at its current location %s\n", currentVersion, wanted.GetTagName(), currentLocation)
		}

		err = binary.Upgrade(currentLocation, tempDir, wanted)

		fmt.Println("upgrade done.")
		return err
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&options.Current.DryRun, "dry-run", options.Current.DryRun, "does not replace the binary")
	rootCmd.PersistentFlags().BoolVar(&options.Current.Force, "force", options.Current.Force, "force the replacement of the binary")
	rootCmd.PersistentFlags().BoolVar(&options.Current.AllowPrerelease, "allow-prerelease", options.Current.AllowPrerelease, "allow the installation of pre-release binary versions")
}

func Command() *cobra.Command {
	pflag.CommandLine.AddFlagSet(rootCmd.Flags())
	return rootCmd
}
