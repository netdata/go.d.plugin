package godplugin

import (
	"os"

	"gopkg.in/yaml.v2"
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

func (c *Config) isModuleEnabled(module string, explicit bool) bool {
	if run, ok := c.Modules[module]; ok {
		return run
	}
	if explicit {
		return false
	}
	return c.DefaultRun
}
