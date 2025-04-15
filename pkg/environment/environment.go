package environment

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"golang.org/x/term"

	"github.com/rs/zerolog/log"
)

func CalledFromTerminal() (bool, error) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		log.Debug().Msg("called from a terminal")
		return true, nil
	} else {
		log.Debug().Msg("NOT called from a terminal")
	}

	return false, nil
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
	log.Debug().Msgf("binary path assessment for \"%s\"", path)
	osAssessment := false

	switch runtime.GOOS {
	case "darwin":
		osAssessment = assessForDarwin(path)
		log.Debug().Msgf("path assessment for darwin: %v", osAssessment)
	case "linux":
		osAssessment = assessForLinux(path)
		log.Debug().Msgf("path assessment for linux: %v", osAssessment)
	default:
		log.Debug().Msgf("default OS path assessment: %v", osAssessment)
		return osAssessment
	}

	// check linux
	return osAssessment ||
		// dependency installation paths are valid
		strings.Contains(path, ".githook-companion/bin") ||
		// valid if go rum is used
		strings.HasSuffix(path, "/go")
}

func assessForUnixFamily(path string) bool {
	return strings.HasPrefix(path, "/usr") ||
		strings.Contains(path, "/.local/bin")
}

func assessForLinux(path string) bool {
	return assessForUnixFamily(path)
}

func assessForDarwin(path string) bool {
	return assessForUnixFamily(path) ||
		strings.Contains(path, "com.apple.Terminal")
}
