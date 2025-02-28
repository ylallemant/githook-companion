package server

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
	"github.com/ylallemant/githook-companion/pkg/command"
)

func Hostname() (string, error) {
	cmd := command.New("git")
	cmd.AddArg("config")
	cmd.AddArg("--get")
	cmd.AddArg("remote.origin.url")

	origin, err := cmd.Execute()
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve origin from config")
	}

	uri, err := url.Parse(origin)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse origin uri %s", origin)
	}

	return uri.Host, nil
}

func Name(OptionalDefaultValue string) (string, error) {
	hostname, err := Hostname()
	if err != nil {
		return "", errors.Wrapf(err, "could not retrieve hostname")
	}

	if name, found := api.Providers[hostname]; found {
		return name, nil
	}

	if OptionalDefaultValue != "" {
		return OptionalDefaultValue, nil
	}

	return hostname, nil
}

func Repository() (string, error) {
	cmd := command.New("git")
	cmd.AddArg("config")
	cmd.AddArg("--get")
	cmd.AddArg("remote.origin.url")

	origin, err := cmd.Execute()
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve origin from config")
	}

	uri, err := url.Parse(origin)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse origin uri %s", origin)
	}

	path := strings.ReplaceAll(uri.Path, ".git", "")

	return fmt.Sprintf("https://%s%s", uri.Host, path), nil
}
