package godplugin

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func newConfig() *config {
	return &config{
		Enabled:    true,
		DefaultRun: true,
	}
}

type config struct {
	Enabled    bool            `yaml:"enabled"`
	DefaultRun bool            `yaml:"default_run"`
	MaxProcs   int             `yaml:"max_procs"`
	Modules    map[string]bool `yaml:"modules"`
}

func (c config) isModuleEnabled(module string, explicit bool) bool {
	if run, ok := c.Modules[module]; ok {
		return run
	}
	if explicit {
		return false
	}
	return c.DefaultRun
}

func newModuleGlobal() *moduleGlobal {
	return &moduleGlobal{
		UpdateEvery:        1,
		AutoDetectionRetry: 0,
	}
}

type moduleGlobal struct {
	UpdateEvery        int
	AutoDetectionRetry int
}

func newModuleConfig() *moduleConfig {
	return &moduleConfig{
		Global: newModuleGlobal(),
	}
}

type moduleConfig struct {
	name   string
	Global *moduleGlobal            `yaml:"global"`
	Jobs   []map[string]interface{} `yaml:"jobs"`
}

func (m *moduleConfig) updateJobs(moduleUpdateEvery, pluginUpdateEvery int) {
	if m.Global == nil {
		m.Global = newModuleGlobal()
	}

	if moduleUpdateEvery > 0 {
		m.Global.UpdateEvery = moduleUpdateEvery
	}

	for _, job := range m.Jobs {
		if _, ok := job["update_every"]; !ok {
			job["update_every"] = m.Global.UpdateEvery
		}

		if _, ok := job["autodetection_retry"]; !ok {
			job["autodetection_retry"] = m.Global.AutoDetectionRetry
		}

		if v, ok := job["update_every"].(int); ok && v < pluginUpdateEvery {
			job["update_every"] = pluginUpdateEvery
		}
	}
}

func load(conf interface{}, filename string) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil && err != io.EOF {
		return err
	}

	if err = yaml.NewDecoder(file).Decode(conf); err != nil {
		return err
	}

	return nil
}
