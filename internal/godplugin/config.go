package godplugin

func newConfig() config {
	return config{
		DefaultRun: true,
		Enabled:    true,
	}
}

type config struct {
	Enabled    bool            `yaml:"enabled"`
	DefaultRun bool            `yaml:"default_run"`
	MaxProcs   int             `yaml:"max_procs,inrange:[1:]"`
	Modules    map[string]bool `yaml:"modules"`
}
