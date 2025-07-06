package settings

// ConfigProject is the schema for a projects' settings file
type ConfigProject struct {
	Languages map[string]ConfigLanguage `yaml:"languages"`
	Ignore    []string                  `yaml:"ignore"`
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
