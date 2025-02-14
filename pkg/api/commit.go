package api

type Commit struct {
	Dictionaries []*CommitTypeDictionary `yaml:"dictionaries" json:"dictionaries"`
	Types        []*CommitType           `yaml:"types" json:"types"`
	DefaultType  string                  `yaml:"defaut_type" json:"defaut_type"`
}

type CommitType struct {
	Type        string `yaml:"type" json:"type"`
	Description string `yaml:"description" json:"description"`
}

type CommitTypeDictionary struct {
	Name     string   `yaml:"name" json:"name"`
	Value    string   `yaml:"value" json:"value"`
	Type     string   `yaml:"type" json:"type"`
	Synonyms []string `yaml:"synonyms" json:"synonyms"`
}
