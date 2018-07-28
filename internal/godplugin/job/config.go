package job

// NewConfig returns Config with default values
func NewConfig() *Config {
	return &Config{
		UpdEvery:           1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		RetriesMax:         60,
	}
}

type Config struct {
	moduleName         string
	jobName            string
	OverrideName       string `yaml:"name"`
	UpdEvery           int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectionRetry int    `yaml:"autodetection_retry" validate:"gte=0"`
	ChartCleanup       int    `yaml:"chart_cleanup" validate:"gte=0"`
	RetriesMax         int    `yaml:"failedUpdates" validate:"gte=0"`
}

// TODO: ModuleName() prepends "go_"
func (c Config) ModuleName() string {
	return "go_" + c.moduleName
}

func (c Config) FullName() string {
	if c.jobName == "" && c.OverrideName == "" {
		return c.ModuleName()
	}
	return c.ModuleName() + "_" + c.JobName()
}

func (c Config) JobName() string {
	if c.OverrideName != "" {
		return c.OverrideName
	}
	if c.jobName != "" {
		return c.jobName
	}
	return c.moduleName
}

func (c Config) UpdateEvery() int {
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
