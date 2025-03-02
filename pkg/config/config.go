package config

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

type Config struct {
	OutputFormat string
	SourcePath   string
	ExcludePaths []string
	Cleaner      struct {
		ErrorIf       bool
		NoneCodeLines bool
		Generated     bool
		CustomIf      []string
	}
	Complexity struct {
		Active bool
		Type   string
	}
}

func ReadConfig(path string) (Config, error) {
	// Read config from file
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	c := Config{}
	if err := json.Unmarshal(content, &c); err != nil {
		return Config{}, err
	}

	return c, nil
}

func (c *Config) Validate() error {
	// Validate config
	if !slices.Contains([]string{"textfmt", "cobertura"}, c.OutputFormat) {
		return fmt.Errorf("output format must be one of textfmt or cobertura")
	}

	if c.SourcePath == "" {
		return fmt.Errorf("source path is empty")
	}

	if c.Complexity.Active {
		if !slices.Contains([]string{"cyclomatic", "cognitive"}, c.Complexity.Type) {
			return fmt.Errorf("complexity type must be one of cyclomatic or cognitive")
		}
	}

	return nil
}
