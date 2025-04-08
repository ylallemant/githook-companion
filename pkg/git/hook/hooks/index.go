package hooks

import "embed"

var (
	//go:embed prepare-commit-msg
	//go:embed pre-commit
	CommonHooks embed.FS

	Hooks = []string{
		"prepare-commit-msg",
		"pre-commit",
	}
)
