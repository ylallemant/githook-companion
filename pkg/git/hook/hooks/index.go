package hooks

import "embed"

var (
	//go:embed prepare-commit-msg
	CommonHooks embed.FS

	Hooks = []string{
		"prepare-commit-msg",
	}
)
