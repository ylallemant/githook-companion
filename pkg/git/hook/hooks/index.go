package hooks

import "embed"

var (
	//go:embed pre-commit
	//go:embed prepare-commit-msg
	//go:embed post-commit
	//go:embed pre-push
	CommonHooks embed.FS

	Hooks = []string{
		"pre-commit",
		"prepare-commit-msg",
		"post-commit",
		"pre-push",
	}
)
