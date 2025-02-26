package options

import (
	"path/filepath"
	"runtime"

	"github.com/ylallemant/githook-companion/pkg/environment"
)

var (
	Current = NewOptions()
)

func NewOptions() *Options {
	options := new(Options)

	if runtime.GOOS == "darwin" {
		options.Directory = "/usr/local/bin"
	} else {
		home, err := environment.Home()
		if err != nil {
			panic(err)
		}

		err = environment.EnsureDirectory(filepath.Join(home, ".local"))
		if err != nil {
			panic(err)
		}

		options.Directory = filepath.Join(home, ".local", "bin")
	}

	return options
}

type Options struct {
	Directory string
}
