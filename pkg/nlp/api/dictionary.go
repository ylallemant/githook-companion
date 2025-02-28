package api

type Dictionary struct {
	Name     string   `yaml:"name" json:"name"`
	Token    string   `yaml:"token" json:"token"`
	Value    string   `yaml:"value" json:"value"`
	Synonyms []string `yaml:"synonyms" json:"synonyms"`
}
