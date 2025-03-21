package config

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/command"
)

const (
	ignoreErrorCheckUnknownSectionOrProperty = "command execution failed: exit status 1"
	ignoreErrorUnsetUnknownSection           = "command execution failed: exit status 2"
	ignoreErrorUnsetUnknownProperty          = "command execution failed: exit status 5"
)

func PropertyExists(property string, globally bool) (bool, error) {
	// unset githook property in local git configuration
	git := command.New("git")
	git.AddArg("config")

	if globally {
		git.AddArg("--global")
	}

	git.AddArg(property)

	message, err := git.Execute()
	if err != nil {
		if err.Error() != ignoreErrorCheckUnknownSectionOrProperty {
			return false, errors.Wrapf(err, "failed to check git property \"%s\" (global=%v): %s", property, globally, message)
		}
	}

	exists := true

	if err != nil {
		exists = false
	}

	fmt.Printf("check for git config property \"%s\" (exists=%v)\n", property, exists)

	return exists, nil
}

func GetProperty(property string, globally bool) (string, error) {
	// unset githook property in local git configuration
	git := command.New("git")
	git.AddArg("config")

	if globally {
		git.AddArg("--global")
	}

	git.AddArg(property)

	message, err := git.Execute()
	if err != nil {
		if err.Error() != ignoreErrorCheckUnknownSectionOrProperty {
			return "", errors.Wrapf(err, "failed to check git property \"%s\" (global=%v): %s", property, globally, message)
		}
	}

	return message, nil
}

func SetProperty(property, value string, globally bool) error {
	// unset githook property in local git configuration
	git := command.New("git")
	git.AddArg("config")

	if globally {
		git.AddArg("--global")
	}

	git.AddArg(property)
	git.AddArg(value)

	message, err := git.Execute()
	if err != nil {
		return errors.Wrapf(err, "failed to set git property \"%s\" = \"%s\" (global=%v)", property, value, globally)
	}

	fmt.Printf("set git config property (global=%v) \"%s\" = \"%s\"\n%s", globally, property, value, message)

	return nil
}

func UnsetProperty(property string, globally bool) error {
	// unset githook property in local git configuration
	git := command.New("git")
	git.AddArg("config")

	if globally {
		git.AddArg("--global")
	}

	git.AddArg("--unset")
	git.AddArg(property)

	_, err := git.Execute()
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() != ignoreErrorUnsetUnknownSection && err.Error() != ignoreErrorUnsetUnknownProperty {
			return errors.Wrapf(err, "failed to unset git property \"%s\" (global=%v)", property, globally)
		}
	}

	fmt.Printf("unset git config property (global=%v) \"%s\"\n", globally, property)
	return nil
}
