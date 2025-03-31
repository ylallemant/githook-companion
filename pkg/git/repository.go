package git

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
)

var (
	gitExtention    = regexp.MustCompile(`.git$`)
	azureSshVersion = regexp.MustCompile(`:v\d+`)
)

func OwnerAndRepositoryFromUri(uri string) (string, string, error) {
	parsed, err := parseGitURI(uri)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(parsed.Path, "/")

	switch Provider(uri) {
	case ProviderAzureDevOps:
		return parts[1], parts[4], err
	default:
		return parts[1], parts[2], err
	}
}

func OwnerFromUri(uri string) (string, error) {
	parsed, err := parseGitURI(uri)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsed.Path, "/")

	return parts[1], err
}

func RepositoryFromUri(uri string) (string, error) {
	parsed, err := parseGitURI(uri)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsed.Path, "/")
	maxIndex := len(parts) - 1

	return parts[maxIndex], err
}

func Hostname() (string, error) {
	origin, err := Origin()

	uri, err := parseGitURI(origin)
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
	signature, err := RepositorySignature("")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s", signature), nil
}

func RepositorySignature(path string) (string, error) {
	origin, err := OriginFromPath(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to retrieve origin from config")
	}

	return RepositorySignatureFromUri(origin)
}

func RepositorySignatureFromUri(uri string) (string, error) {
	parsed, err := parseGitURI(uri)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse origin uri %s", uri)
	}

	return fmt.Sprintf("%s%s", parsed.Host, parsed.Path), nil
}

func parseGitURI(uri string) (*url.URL, error) {
	uri = NormaliseUri(uri)

	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse git uri %s", uri)
	}

	return parsed, nil
}

func NormaliseUri(uri string) string {
	switch Provider(uri) {
	case ProviderAzureDevOps:
		uri = normaliseAzureDevOpsUri(uri)
	default:
		uri = nomaliseGitHubLikeUri(uri)
	}

	return uri
}

func nomaliseGitHubLikeUri(uri string) string {
	isGitProtocol := strings.HasPrefix(uri, "git@")
	if isGitProtocol {
		uri = strings.Replace(uri, ":", "/", 1)
		uri = strings.Replace(uri, "git@", "https://", 1)
	}

	isGitUri := strings.HasPrefix(uri, "git://")
	if isGitUri {
		uri = strings.Replace(uri, "git://", "https://", 1)
	}

	uri = gitExtention.ReplaceAllString(uri, "")

	return uri
}

func normaliseAzureDevOpsUri(uri string) string {
	if strings.Contains(uri, "git@ssh.") {
		lastSlach := strings.LastIndex(uri, "/")
		repository := uri[lastSlach+1:]

		uri = azureSshVersion.ReplaceAllString(uri, "")
		uri = strings.Replace(uri, "git@ssh.", "https://", 1)
		uri = strings.Replace(uri, repository, fmt.Sprintf("_git/%s", repository), 1)
	}

	return uri
}
