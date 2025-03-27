package config

import (
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/globals"
)

func InitContext() (api.ConfigContext, error) {
	globals.ProcessGlobals()

	var configContext api.ConfigContext
	var err error

	if globals.Current.ConfigPath != "" {
		configContext, err = ContextFromPath(globals.Current.ConfigPath, globals.Current.FallbackConfig)
		if err != nil {
			return nil, err
		}
	} else {
		configContext, err = Context(globals.Current.FallbackConfig)
		if err != nil {
			return nil, err
		}
	}

	return configContext, nil
}
