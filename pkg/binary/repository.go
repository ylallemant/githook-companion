package binary

import (
	"github.com/ylallemant/githook-companion/pkg/git"
)

var repository string

var (
	defaultRepository = "git@github.com:test/some-repo.git"
	uri               = ""
)

func init() {
	uri = git.NormaliseUri(Repository())
}

func Repository() string {
	return getOr(repository, defaultRepository)
}

func Uri() string {
	return uri
}
