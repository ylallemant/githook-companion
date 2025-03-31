package git

import "strings"

const (
	ProviderGitHub      = "github.com"
	ProviderAzureDevOps = "dev.azure.com"
	ProviderUnknown     = "unknown Git provider"
)

func Provider(uri string) string {
	if strings.Contains(uri, ProviderGitHub) {
		return ProviderGitHub
	}

	if strings.Contains(uri, ProviderAzureDevOps) {
		return ProviderAzureDevOps
	}

	return ProviderUnknown
}
