package settings

// ConfigProject is the schema for a projects' settings file
type ConfigProject struct {
	Languages map[string]ConfigLanguage `json:"languages" yaml:"languages" jsonschema_description:"a map of languages that crie should be able to run"`
	Ignore    []string                  `json:"ignore" yaml:"ignore" jsonschema_description:"list of regexes matched against the file list to ignore them (exact paths also work)"`
}

func (c *ConfigProject) merge(src ConfigProject) {

	for langName, lang := range src.Languages {
		if existing, exists := c.Languages[langName]; exists {
			c.Languages[langName] = existing
		} else {
			c.Languages[langName] = lang
		}
	}

	c.Ignore = append(c.Ignore, src.Ignore...)
}
