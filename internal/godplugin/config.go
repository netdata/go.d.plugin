package godplugin

import (
	"os"

	"github.com/go-yaml/yaml"
)

func NewConfig() *Config {
	return &Config{
		DefaultRun: true,
		Enabled:    true,
	}
}

type Config struct {
	Enabled    bool            `yaml:"enabled"`
	DefaultRun bool            `yaml:"default_run"`
	MaxProcs   int             `yaml:"max_procs"`
	Modules    map[string]bool `yaml:"modules"`
}

func (c *Config) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	return yaml.NewDecoder(file).Decode(c)
}
