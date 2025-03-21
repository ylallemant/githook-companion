package config

import (
	"fmt"

	"github.com/ylallemant/githook-companion/pkg/dependency"
)

func Debug() string {
	localPath, _ := GetLocalPath()
	globalPath, _ := GetGlobalPath()

	localExists, _, _ := DirectoryExists(localPath)
	globalExists, _, _ := DirectoryExists(globalPath)

	configDebug := ""

	if localExists {
		configDebug = debugConfig(localPath)
	}

	if globalExists {
		configDebug = debugConfig(globalPath)
	}

	return fmt.Sprintf(
		`
++++ CONFIG(S) ++++++++++++++++++++
  local path  (exits=%v):  "%s"
  global path (exits=%v):  "%s"

%s
  `,
		localExists,
		localPath,
		globalExists,
		globalPath,
		configDebug,
	)
}

func debugConfig(path string) string {
	cfg, err := Load(path, true)
	if err != nil {
		panic("failes to load local config")
	}

	referenceRepository := "none"
	referencePath := "none"
	referencePathExists := false

	if cfg.ParentConfig != nil {
		referenceRepository = cfg.ParentConfig.GitRepository
		referencePath = ParentPathFromConfig(cfg)
		referencePathExists, _, _ = DirectoryExists(referencePath)
	}

	dependencyDirectory := dependency.InstallDirectoryFromConfig(cfg)
	dependencyDirectoryExists, _, _ := DirectoryExists(dependencyDirectory)

	hookDirectory := GithooksPathFromConfig(cfg)
	hookDirectoryExists, _, _ := DirectoryExists(hookDirectory)

	return fmt.Sprintf(
		`
-- parent ---------
  repository: %s
  directory (exists=%v): %s

-- dependencies ---
  count: %d
  directory (exists=%v): %s

-- hooks ----------
  directory (exists=%v): %s
  `,
		referenceRepository,
		referencePathExists,
		referencePath,
		len(cfg.Dependencies),
		dependencyDirectoryExists,
		dependencyDirectory,
		hookDirectoryExists,
		hookDirectory,
	)
}
