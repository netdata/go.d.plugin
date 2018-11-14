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
		setDefault(m.Global, key, val)
	}

	for _, job := range m.Jobs {
		for key, val := range m.Global {
			setDefault(job, key, val)
		}
	}
}

func setDefault(m map[string]interface{}, key string, val interface{}) {
	if _, ok := m[key]; !ok {
		m[key] = val
	}
}
