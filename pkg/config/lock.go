package config

import (
	"path/filepath"
	"time"

	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

func GithookLockPathFromNameAndConfig(name string, ctx api.ConfigContext) string {
	return filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), name)
}

func SetPermanentLock(filename string, ctx api.ConfigContext) error {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.SetPermanentLock(path)
}

func SetPermanentLockWithDescription(filename, description string, ctx api.ConfigContext) error {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.SetPermanentLockWithDescription(path, description)
}

func SetTimedLock(filename string, duration time.Duration, ctx api.ConfigContext) error {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.SetTimedLock(path, duration)
}

func SetTimedLockWithDescription(filename, description string, duration time.Duration, ctx api.ConfigContext) error {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.SetTimedLockWithDescription(path, description, duration)
}

func PermanentLockExists(filename string, ctx api.ConfigContext) (bool, error) {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.PermanentLockExists(path)
}

func TimeLockActive(filename string, ctx api.ConfigContext) (bool, error) {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.TimeLockActive(path)
}

func LockRemove(filename string, ctx api.ConfigContext) error {
	path := filepath.Join(ContextDirectoryFromBase(ctx.LocalPath()), filename)
	return filesystem.RemoveLock(path)
}
