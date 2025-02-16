package api

const ConfigDirectory = ".githooks-butler"
const ConfigFile = "config.yaml"

type Config struct {
	*Commit `yaml:"commit" json:"commit"`
}
