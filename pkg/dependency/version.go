package dependency

import (
	"fmt"
	"os/exec"

	"github.com/ylallemant/githook-companion/pkg/api"
)

func Version(dependency *api.Dependency) string {
	return fmt.Sprintf("%s%s", dependency.SemverPrefix, dependency.Version)
}

func AvailableVersion(binary string) (string, error) {
	var output []byte

	// try call with "--version"
	cmd1 := exec.Command(binary, "version")
	output, err := cmd1.CombinedOutput()
	if err != nil {
		output = []byte{}
	}

	if len(output) == 0 {
		// try call with "version"
		cmd2 := exec.Command(binary, "--version")
		output, err = cmd2.CombinedOutput()
		if err != nil {
			return "", err
		}
	}

	return string(output), nil
}
