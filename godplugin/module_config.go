package godplugin

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type rawConfig map[string]interface{}

func (r *rawConfig) merge(src rawConfig) {
	for key, val := range src {
		if _, ok := (*r)[key]; !ok {
			(*r)[key] = val
		}
	}
}

type moduleConfig struct {
	Global rawConfig   `yaml:"global"`
	Jobs   []rawConfig `yaml:"jobs"`
}

func (m *moduleConfig) load(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil && err != io.EOF {
		return err
	}

	return yaml.NewDecoder(file).Decode(m)
}

func (m *moduleConfig) updateJobs(modUpdateEvery, globalUpdateEvery int) {
	defaults := rawConfig{
		"update_every":        1,
		"autodetection_retry": 0,
	}
	if modUpdateEvery > 0 {
		defaults["update_every"] = modUpdateEvery
	}

	for _, job := range m.Jobs {
		job.merge(m.Global)
		job.merge(defaults)

		if v, ok := job["update_every"].(int); ok && v < globalUpdateEvery {
			job["update_every"] = globalUpdateEvery
		}
	}
}
