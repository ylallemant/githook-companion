package git

import (
	"bufio"
	"fmt"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/environment"
	"github.com/ylallemant/githook-companion/pkg/filesystem"
)

const gitCredentialsRelativePath = "~/.git-credentials"

var (
	credentials        = []*url.URL{}
	gitCredentialsPath = ""
)

func init() {
	var err error
	gitCredentialsPath, err = environment.EnsureAbsolutePath(gitCredentialsRelativePath)
	if err != nil {
		panic(fmt.Sprintf("failed to ensure absolute path for %s: %s", gitCredentialsRelativePath, err.Error()))
	}
}

func LoadCredentials() error {
	log.Debug().Msgf("loading git credentials from \"%s\"", gitCredentialsPath)
	exists, _, err := filesystem.FileExists(gitCredentialsPath)
	if err != nil {
		return errors.Wrapf(err, "failed to check existance of %s", gitCredentialsPath)
	}

	if exists {
		file, err := os.Open(gitCredentialsPath)
		if err != nil {
			return errors.Wrapf(err, "failed to open %s", gitCredentialsPath)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			uri, err := url.Parse(line)
			if err != nil {
				return errors.Wrapf(err, "failed to parse git credential entry \"%s\"", line)
			}

			credentials = append(credentials, uri)
		}
	} else {
		log.Warn().Msgf("file does not exist \"%s\"", gitCredentialsPath)
	}

	return nil
}

func UserInfoFromUri(uri *url.URL) (*url.Userinfo, bool, error) {
	if len(credentials) == 0 {
		LoadCredentials()
	}

	log.Debug().Msgf("check %d credentials for repository %s%s", len(credentials), uri.Hostname(), uri.Path)

	filtered := filterCredentials(uri)
	log.Debug().Msgf("%d credentials left after filtering", len(filtered))

	if len(filtered) == 0 {
		return nil, false, nil
	}

	if len(filtered) == 1 {
		return filtered[0], true, nil
	}

	return nil, false, nil
}

func AuthMethodFromUri(uri string) (transport.AuthMethod, error) {
	repositoryUri, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse repository uri \"%s\"", uri)
	}

	userInfo, found, err := UserInfoFromUri(repositoryUri)
	if err != nil {
		return nil, err
	}

	if found {
		password, _ := userInfo.Password()

		authMethod := &http.BasicAuth{
			Username: userInfo.Username(),
			Password: password,
		}

		log.Debug().Msgf("returns a BasicAuth method for user %s", userInfo.Username())
		return authMethod, nil
	}

	// fallback to a ssh-agent method
	authMethod, err := ssh.NewSSHAgentAuth("git")
	if err != nil {
		return nil, errors.Wrap(err, "failed initialte a SSH-Agent method")
	}

	log.Debug().Msg("returns a SSH-Agent method")
	return authMethod, nil
}

func filterCredentials(uri *url.URL) []*url.Userinfo {
	filtered := filterByHostname(uri, credentials)

	if len(filtered) == 1 {
		return filtered
	}

	// try to filter by user
	filtered = filterByUser(uri, filtered)

	return filtered
}

func filterByHostname(uri *url.URL, list []*url.URL) []*url.Userinfo {
	filtered := make([]*url.Userinfo, 0)
	for _, credential := range list {
		if credential.Hostname() == uri.Hostname() {
			if credential.User != nil {
				filtered = append(filtered, credential.User)
			}
		}
	}
	return filtered
}

func filterByUser(uri *url.URL, list []*url.Userinfo) []*url.Userinfo {
	filtered := make([]*url.Userinfo, 0)
	for _, user := range list {
		if uri.User != nil {
			if user.Username() == uri.User.Username() {
				filtered = append(filtered, user)
			}
		}
	}
	return filtered
}
