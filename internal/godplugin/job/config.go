package job

type Config struct {
	RealModuleName     string
	RealJobName        string
	OverrideName       string `yaml:"name"`
	UpdateEvery        int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectionRetry int    `yaml:"autodetection_retry" validate:"gte=0"`
	ChartCleanup       int    `yaml:"chart_cleanup" validate:"gte=0"`
	MaxRetries         int    `yaml:"retries" validate:"gte=0"`
}

// NewConfig returns Config with default values
func NewConfig() *Config {
	return &Config{
		UpdateEvery:        1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		MaxRetries:         60,
	}
}

// TODO: ModuleName() prepends "go_"
func (c *Config) ModuleName() string {
	return "go_" + c.RealModuleName
}

func (c *Config) FullName() string {
	if c.RealJobName == "" && c.OverrideName == "" {
		return c.ModuleName()
	}
	return c.ModuleName() + "_" + c.JobName()
}

func (c *Config) JobName() string {
	if c.OverrideName != "" {
		return c.OverrideName
	}
	if c.RealJobName != "" {
		return c.RealJobName
	}
	return c.RealModuleName
}
