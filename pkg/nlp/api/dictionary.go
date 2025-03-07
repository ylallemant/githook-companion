package api

type Dictionary struct {
	LanguageCode         string   `yaml:"language_code" json:"language_code"`
	Name                 string   `yaml:"name" json:"name"`
	TokenName            string   `yaml:"token" json:"token"`
	Entries              []string `yaml:"synonyms" json:"synonyms"`
	ConfidenceThresthold float64
}
