package api

type Dictionary struct {
	Name                 string   `yaml:"name" json:"name"`
	Description          string   `yaml:"description" json:"description"`
	LanguageCode         string   `yaml:"language_code" json:"language_code"`
	TokenName            string   `yaml:"token_name" json:"token_name"`
	TokenValue           string   `yaml:"token_value" json:"token_value"`
	TokenValueIsMatch    bool     `yaml:"token_value_is_match" json:"token_value_is_match"`
	Entries              []string `yaml:"synonyms" json:"synonyms"`
	ConfidenceThresthold float64  `yaml:"confidence_thresthold" json:"confidence_thresthold"`
}
