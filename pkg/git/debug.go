package git

import (
	"fmt"

	"github.com/ylallemant/githook-companion/pkg/git/config"
)

func Debug() string {

	hostname, err := Hostname()
	if err != nil {
		panic(fmt.Sprintf("failed to get git hostname: %s", err.Error()))
	}

	repository, err := Repository()
	if err != nil {
		panic(fmt.Sprintf("failed to get git repository: %s", err.Error()))
	}

	hookPath, err := config.GetProperty("core.hooksPath", false)
	if err != nil {
		panic(fmt.Sprintf("failed to get git config core.hooksPath: %s", err.Error()))
	}

	user, err := config.GetProperty("user.name", false)
	if err != nil {
		panic(fmt.Sprintf("failed to get git config user.name: %s", err.Error()))
	}

	email, err := config.GetProperty("user.email", false)
	if err != nil {
		panic(fmt.Sprintf("failed to get git config user.email: %s", err.Error()))
	}

	return fmt.Sprintf(
		`
++++ GIT ++++++++++++++++++++++++++
  repository:   "%s"
  hostname:     "%s"

  core.hooksPath: "%s"
  user.name:      "%s"
  user.email:     "%s"
  `,
		repository,
		hostname,
		hookPath,
		user,
		email,
	)
}
