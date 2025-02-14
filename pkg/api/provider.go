package api

const (
	Github      = "github"
	AzureDevOps = "azure-devops"
)

var (
	Providers = map[string]string{
		"github.com":    Github,
		"dev.azure.com": AzureDevOps,
	}
)
