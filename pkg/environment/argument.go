package environment

import "regexp"

var (
	argumentPrefixRegexp = regexp.MustCompile("^(-|--)")
)

func IsAnArgument(text string) bool {
	if argumentPrefixRegexp.MatchString(text) {
		return true
	}

	return false
}
