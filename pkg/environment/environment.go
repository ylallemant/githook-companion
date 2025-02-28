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
	// absolute path of the called executable
	underscore := FindEnvVar("_")

	assertion := pathValid(dirname, oldpwd) &&
		binaryPathValid(underscore)

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

func pathValid(path, oldpath string) bool {
	pathDepth := len(strings.Split(path, string(os.PathSeparator)))
	oldpathDepth := len(strings.Split(oldpath, string(os.PathSeparator)))

	depthDistance := pathDepth - oldpathDepth

	// path and oldpath should start similary but oldpath may be a parent directory
	// but with a maximum of 2 depth level away
	return strings.HasPrefix(path, oldpath) && depthDistance < 3
}

func binaryPathValid(path string) bool {
	// check linux
	return strings.HasPrefix(path, "/usr") ||
		// valid if go rum is used
		strings.HasSuffix(path, "/go")
}
