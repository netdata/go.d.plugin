package godplugin

import (
	"os"

	"gopkg.in/yaml.v2"
)

var globalDefaults = map[string]interface{}{
	"update_every":        1,
	"autodetection_retry": 0,
}

type moduleConfig struct {
	Global map[string]interface{}
	Jobs   []map[string]interface{}
}

func (m *moduleConfig) load(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return err
	}

	return yaml.NewDecoder(file).Decode(m)
}

func (m *moduleConfig) updateJobConfigs() {
	for key, val := range globalDefaults {
		if _, ok := m.Global[key]; !ok {
			m.Global[key] = val
		}
	}

	for _, job := range m.Jobs {
		for key, val := range m.Global {
			if _, ok := job[key]; !ok {
				job[key] = val
			}
		}
	}
}
