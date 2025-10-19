package language

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Regex wraps regexp.Regexp with custom YAML unmarshaling
type Regex struct {
	*regexp.Regexp
}

// UnmarshalYAML implements custom YAML unmarshalling
func (cr *Regex) UnmarshalYAML(value *yaml.Node) error {
	var patterns []string
	if err := value.Decode(&patterns); err != nil {
		return err
	}

	if len(patterns) > 0 {
		compiled := regexp.MustCompile(strings.Join(patterns, "|"))
		cr.Regexp = compiled
	}

	return nil
}
