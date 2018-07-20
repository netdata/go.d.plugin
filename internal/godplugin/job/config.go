package job

// NewConf returns Config with default values
func NewConf() *Config {
	return &Config{
		UpdEvery:           1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		RetriesMax:         60,
	}
}

type Config struct {
	moduleName         string // standalone struct ?
	jobName            string // standalone struct ?
	OverrideName       string `yaml:"name"`
	UpdEvery           int    `yaml:"update_every"`
	AutoDetectionRetry int    `yaml:"autodetection_retry"`
	ChartCleanup       int    `yaml:"chart_cleanup"`
	RetriesMax         int    `yaml:"retries"`
}

// TODO: ModuleName() prepends "go_"
func (c *Config) ModuleName() string {
	return "go_" + c.moduleName
}

func (c *Config) FullName() string {
	if c.jobName == "" {
		return c.ModuleName()
	}
	return c.ModuleName() + "_" + c.JobName()
}

func (c *Config) JobName() string {
	if c.jobName == "" {
		return c.ModuleName()
	}
	if c.OverrideName == "" {
		return c.jobName
	}
	return c.OverrideName
}

func (c *Config) UpdateEvery() int {
	return c.UpdEvery
}

func (c *Config) SetModuleName(name string) {
	c.moduleName = name
}

func (c *Config) SetJobName(name string) {
	c.jobName = name
}

func (c *Config) SetUpdateEvery(u int) {
	c.UpdEvery = u
}
