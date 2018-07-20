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
	MaxProcs   int             `yaml:"max_procs"`
	Modules    map[string]bool `yaml:"modules"`
}
