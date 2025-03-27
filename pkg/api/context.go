package api

type ConfigContext interface {
	LocalPath() string
	ParentPath() string
	LocalConfig() *Config
	ParentConfig() *Config
	Config() *Config
	HasParent() bool
}
