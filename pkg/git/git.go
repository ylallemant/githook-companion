package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/config"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

var entries = []string{
	`
# githook-companion rules`,
	filepath.Join(api.ConfigDirectory, "bin", ""),
	filepath.Join(api.ConfigDirectory, "locks", ""),
}

func EnsureGitIgnoreFromConfig(configuration *api.Config) error {
	path, err := config.BasePathFromConfig(configuration)
	if err != nil {
		return err
	}

	path = filepath.Join(path, ".gitignore")

	exists, _, err := filesystem.FileExists(path)
	if err != nil {
		return err
	}

	content := ""

	if exists {
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content = string(raw)
	}

	for _, entry := range entries {
		index := strings.Index(content, entry)

		if index < 0 {
			content = fmt.Sprintf(`%s
%s`,
				content,
				entry,
			)
		}
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", path)
	}

	return nil
}
