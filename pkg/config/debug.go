package config

import (
	"fmt"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/dependency"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

func Debug() string {
	localPath, _ := GetLocalFilePath()
	globalPath, _ := GetGlobalFilePath()

	localExists, _, _ := filesystem.DirectoryExists(localPath)
	globalExists, _, _ := filesystem.DirectoryExists(globalPath)

	configDebug := ""

	if localExists {
		configDebug = debugConfigs(localPath, "local")
	}

	if globalExists {
		configDebug = debugConfigs(globalPath, "global")
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

func debugConfigs(path, origin string) string {
	originConfig, err := Load(path, true)
	if err != nil {
		panic(fmt.Sprintf("failed to load origin config: %s", err.Error()))
	}

	referenceRepository := "none"
	referencePath := "none"
	referencePathExists := false
	parent := ""
	var parentConfig *api.Config

	if originConfig.ParentConfig != nil {
		referenceRepository = originConfig.ParentConfig.GitRepository
		referencePath = ParentPathFromConfig(originConfig)
		referencePathExists, _, _ = filesystem.DirectoryExists(referencePath)

		parent = fmt.Sprintf(`
  parent path (exits=%v):  "%s"
     -> repository: %s
		`,
			referencePathExists,
			referencePath,
			referenceRepository,
		)

		parentConfig, err = Load(FilePathFromBase(referencePath), true)
		if err != nil {
			panic(fmt.Sprintf("failed to load parent config: %s", err.Error()))
		}
	}

	fetchedConfig, err := Get()
	if err != nil {
		panic(fmt.Sprintf("failed to get config: %s", err.Error()))
	}

	originConfigDebug := debugConfig(originConfig, origin)
	parentConfigDebug := debugConfig(parentConfig, "parent")
	fetchedConfigDebug := debugConfig(fetchedConfig, "resulting")

	return fmt.Sprintf(`%s%s%s%s
	`,
		parent,
		originConfigDebug,
		parentConfigDebug,
		fetchedConfigDebug,
	)
}

func debugConfig(cfg *api.Config, name string) string {
	if cfg == nil {
		return ""
	}

	dependencyDirectory := dependency.DependencyDirectoryFromConfig(cfg)
	dependencyDirectoryExists, _, _ := filesystem.DirectoryExists(dependencyDirectory)

	hookDirectory := GithooksPathFromConfig(cfg)
	hookDirectoryExists, _, _ := filesystem.DirectoryExists(hookDirectory)

	dictionaries := 0
	lexemes := 0

	if cfg.Commit != nil && cfg.Commit.TokenizerOptions != nil {
		dictionaries = len(cfg.Commit.TokenizerOptions.Dictionaries)
		lexemes = len(cfg.Commit.TokenizerOptions.Lexemes)
	}

	return fmt.Sprintf(
		`
   -- %s config ---------
     parent:       %v
     dictionaries: %d
     lexemes:      %d
     dependencies: %d

     dependencies
       directory (exists=%v): %s

     hooks
       directory (exists=%v): %s
  `,
		name,
		(cfg.ParentConfig != nil),
		dictionaries,
		lexemes,
		len(cfg.Dependencies),
		dependencyDirectoryExists,
		dependencyDirectory,
		hookDirectoryExists,
		hookDirectory,
	)
}
