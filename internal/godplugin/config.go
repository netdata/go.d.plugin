package godplugin

import (
	"os"

	"io"

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
	err = yaml.NewDecoder(file).Decode(c)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (c *Config) IsModuleEnabled(module string, explicit bool) bool {
	if run, ok := c.Modules[module]; ok {
		return run
	}
	if explicit {
		return false
	}
	return c.DefaultRun
}
