package godplugin

import (
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfig create new config
func NewConfig() *Config {
	return &Config{
		DefaultRun: true,
		Enabled:    true,
	}
}

// Config go.d.conf config
type Config struct {
	Enabled    bool            `yaml:"enabled"`
	DefaultRun bool            `yaml:"default_run"`
	MaxProcs   int             `yaml:"max_procs"`
	Modules    map[string]bool `yaml:"modules"`
}

// Load load go.d.conf config
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
