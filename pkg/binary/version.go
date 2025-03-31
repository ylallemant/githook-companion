package binary

import (
	"fmt"
	"strings"
)

var semverVersion string
var gitCommitHash string

var (
	defaultSemverVersion = "n/a"
	defaultGitCommitHash = "dirty"
)

func Information() string {
	return fmt.Sprintf(
		`version: %s, commit: %s
source: %s`,
		getOr(semverVersion, defaultSemverVersion),
		getOr(gitCommitHash, defaultGitCommitHash),
		uri,
	)
}

func Commit() string {
	return getOr(gitCommitHash, defaultGitCommitHash)
}

func Semver() string {
	return getOr(semverVersion, defaultSemverVersion)
}

func SemverWithSeparator(sep string) string {
	return strings.ReplaceAll(Semver(), ".", sep)
}

func getOr(this, or string) string {
	if len(this) == 0 {
		return or
	}
	return this
}
