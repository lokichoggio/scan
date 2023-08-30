package config

import (
	"os"

	yaml "gopkg.in/yaml.v2"
)

func Load(file string) (*Config, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err = yaml.Unmarshal(content, c); err != nil {
		return nil, err
	}

	return c, nil
}
