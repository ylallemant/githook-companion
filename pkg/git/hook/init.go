package hook

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
	"github.com/ylallemant/githook-companion/pkg/git/hook/hooks"
)

func Ensure(configuration *api.Config) error {
	hookDirectory := config.GithooksPathFromConfig(configuration)

	err := ensureHooksDirectory(hookDirectory)
	if err != nil {
		return err
	}

	for _, hook := range hooks.Hooks {
		content, err := hooks.CommonHooks.ReadFile(hook)
		if err != nil {
			return err
		}

		hookPath := filepath.Join(hookDirectory, hook)

		exists, _, err := filesystem.FileExists(hookPath)
		if err != nil {
			return err
		}

		if !exists {
			fmt.Printf("- write hook \"%s\" at %s\n", hook, hookDirectory)

			err = os.WriteFile(hookPath, content, 0755)
			if err != nil {
				return errors.Wrapf(err, "failed to write hook \"%s\" at %s", hook, hookDirectory)
			}
		}

	}

	return nil
}

func ensureHooksDirectory(path string) error {
	exists, _, err := filesystem.DirectoryExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", path)
	}

	if !exists {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return errors.Wrapf(err, "failed create directory %s", path)
		}
	}

	return nil
}
