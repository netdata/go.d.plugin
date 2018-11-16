package godplugin

import (
	"os"

	"gopkg.in/yaml.v2"
)

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

func createJobConfigs(modConf moduleConfig, modUpdateEvery *int, globalUpdateEvery int) []map[string]interface{} {
	defaults := map[string]interface{}{
		"update_every":        1,
		"autodetection_retry": 0,
	}

	if modUpdateEvery != nil {
		defaults["update_every"] = *modUpdateEvery
	}

	for _, job := range modConf.Jobs {
		merge(job, modConf.Global)
		merge(job, defaults)

		if v, ok := job["update_every"].(int); ok && v < globalUpdateEvery {
			job["update_every"] = globalUpdateEvery
		}
	}

	return modConf.Jobs
}

func merge(dst, src map[string]interface{}) {
	for key, val := range src {
		if _, ok := dst[key]; !ok {
			dst[key] = val
		}
	}
}
