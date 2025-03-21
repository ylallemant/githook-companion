package environment

import (
	"fmt"
	"math"
	"os"
	"runtime"
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

	depthDistance := int(math.Abs(float64(pathDepth - oldpathDepth)))

	pathStartEqual := strings.HasPrefix(path, oldpath)

	if pathDepth < oldpathDepth {
		pathStartEqual = strings.HasPrefix(oldpath, path)
	}

	// path and oldpath should start similary but oldpath may be a parent directory
	// but with a maximum of 2 depth level away
	return pathStartEqual && depthDistance < 3
}

func binaryPathValid(path string) bool {
	osAssessment := false

	switch runtime.GOOS {
	case "dawin":
		osAssessment = assessForDarwin(path)
	case "linux":
		osAssessment = assessForLinux(path)
	default:
		return false
	}
	fmt.Println(`osAssessment`, osAssessment)
	// check linux
	return osAssessment ||
		// valid if go rum is used
		strings.HasSuffix(path, "/go")
}

func assessForUnixFamily(path string) bool {
	return strings.HasPrefix(path, "/usr")
}

func assessForLinux(path string) bool {
	return assessForUnixFamily(path)
}

func assessForDarwin(path string) bool {
	fmt.Println(`strings.Contains(path, "com.apple.Terminal")`, strings.Contains(path, "com.apple.Terminal"))
	return assessForUnixFamily(path) ||
		strings.Contains(path, "com.apple.Terminal")
}
