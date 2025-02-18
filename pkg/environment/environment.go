package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func CalledFromTerminal() (bool, error) {
	dirname, err := os.Getwd()
	if err != nil {
		return false, errors.Wrapf(err, "failed to get local directory")
	}

	oldpwd := FindEnvVar("OLDPWD")
	underscore := FindEnvVar("_")

	assertion := oldpwd == dirname &&
		strings.HasPrefix(underscore, "/usr")

	return assertion, nil
}

func FindEnvVar(name string) string {
	envVars := os.Environ()
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, name) {
			return strings.ReplaceAll(envVar, fmt.Sprintf("%s=", name), "")
		}
	}

	return ""
}
