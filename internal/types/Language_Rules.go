package types

import "regexp"

type PatternList []string

type LanguageList []string

type Heuristics struct {
	NamedPatterns   map[string]PatternList `yaml:"named_patterns"`
	Disambiguations []Disambiguation       `yaml:"disambiguations"`
}

type Disambiguation struct {
	Extensions []string `yaml:"extensions"`
	Rules      []Rule   `yaml:"rules"`
}

type Rule struct {
	Language        LanguageList `yaml:"language"`
	Pattern         PatternList  `yaml:"pattern"`
	NegativePattern PatternList  `yaml:"negative_pattern"`
	NamedPattern    string       `yaml:"named_pattern"`
	And             []Rule       `yaml:"and"`
}

type CompiledRule struct {
	Language string

	Patterns []string

	Negative []string
}

type HeuristicRule struct {
	Language string
	Patterns []*regexp.Regexp
	Negative []*regexp.Regexp
}
