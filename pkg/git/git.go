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
}

func init() {
	for _, directory := range api.ConfigProcessingDirectories {
		entries = append(entries, filepath.Join(api.ConfigDirectory, directory))
	}
}

func EnsureGitIgnoreFromConfig(configuration *api.Config) error {
	path, err := config.BasePathFromConfig(configuration)
	if err != nil {
		return err
	}

	path = filepath.Join(path, ".gitignore")
	fmt.Println("ensure Git exclusion rules in .gitignore:", path)

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
	fmt.Println("  - content:", content)

	for _, entry := range entries {
		index := strings.Index(content, entry)

		if index < 0 {
			fmt.Println("  - add exclusion rule:", entry)
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
