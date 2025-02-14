package api

type Tool struct {
	Name                string
	Version             string
	Args                []string
	DownloadUrlTemplate string
}

type ToolList struct {
	Tools []Tool
}
