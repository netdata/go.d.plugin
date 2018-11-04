package godplugin

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func NewConfig() *Config {
	return &Config{
		DefaultRun: true,
		Enabled:    true,
	}
}

type GlobalConfig struct {
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

func (c *Config) IsModuleEnabled(module string, explicit bool) bool {
	if run, ok := c.Modules[module]; ok {
		return run
	}
	if explicit {
		return false
	}
	return c.DefaultRun
}

type rawModConfig []yaml.MapSlice

func (r *rawModConfig) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	var raw yaml.MapSlice

	err = yaml.NewDecoder(file).Decode(&raw)
	if err != nil && err != io.EOF {
		return err
	}

	var (
		global yaml.MapSlice
		jobs   yaml.MapSlice
	)

	for _, v := range raw {
		if _, ok := v.Value.(yaml.MapSlice); ok {
			jobs = append(jobs, v)
		} else {
			global = append(global, v)
		}
	}

	if len(jobs) == 0 {
		global = append(global, yaml.MapItem{Key: "job_name", Value: ""})
		*r = append(*r, global)
		return nil
	} else if len(jobs) == 1 {
		job := merge(jobs[0].Value.(yaml.MapSlice), global)
		job = append(job, yaml.MapItem{Key: "job_name", Value: ""})
		*r = append(*r, job)
		return nil
	}
	for _, conf := range jobs {
		job := merge(conf.Value.(yaml.MapSlice), global)
		job = append(job, yaml.MapItem{Key: "job_name", Value: conf.Key.(string)})
		*r = append(*r, job)
	}
	return nil
}

func merge(dst, src yaml.MapSlice) yaml.MapSlice {
	for _, v := range src {
		if hasKey(dst, v) {
			continue
		}
		dst = append(dst, v)
	}
	return dst
}

func hasKey(item yaml.MapSlice, key yaml.MapItem) bool {
	for _, v := range item {
		if key.Key == v.Key {
			return true
		}
	}
	return false
}
