package version

import "fmt"

var Version string
var GitCommit string

var (
	defaultVersion   = "n/a"
	defaultGitCommit = "dirty"
)

func GetInfo() string {
	return fmt.Sprintf("version: %s, commit: %s", getOr(Version, defaultVersion), getOr(GitCommit, defaultGitCommit))
}

func Commit() string {
	return getOr(GitCommit, defaultGitCommit)
}

func Semver() string {
	return getOr(Version, defaultVersion)
}

func getOr(this, or string) string {
	if len(this) == 0 {
		return or
	}
	return this
}
