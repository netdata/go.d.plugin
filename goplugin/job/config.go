package job

// NewConf return *Conf with default values
func NewConf() *Conf {
	return &Conf{
		UpdEvery:           1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		RetriesMax:         60,
	}
}

// Conf is a job base configuration struct
// Used in unmarshal of [global] and [.base] ast.Tables.
type Conf struct {
	moduleName         string
	jobName            string
	OverrideName       string `toml:"name,                regexd:<RE>[^[:word:]]</RE>"`
	UpdEvery           int    `toml:"update_every,        range:[1:]"`
	AutoDetectionRetry int    `toml:"autodetection_retry, range:[0:]"`
	ChartCleanup       int    `toml:"chart_cleanup,       range:[0:]"`
	RetriesMax         int    `toml:"retries,             range:[0:]"`
}

// TODO ModuleName() prepends "go_"
func (c *Conf) ModuleName() string {
	return "go_" + c.moduleName
}

func (c *Conf) FullName() string {
	if c.jobName == "" {
		return c.ModuleName()
	}
	return c.ModuleName() + "_" + c.JobName()
}

func (c *Conf) JobName() string {
	if c.jobName == "" {
		return c.ModuleName()
	}
	if c.OverrideName == "" {
		return c.jobName
	}
	return c.OverrideName
}

func (c *Conf) SetModuleName(name string) {
	c.moduleName = name
}

func (c *Conf) SetJobName(name string) {
	c.jobName = name
}

func (c *Conf) UpdateEvery() int {
	return c.UpdEvery
}

func (c *Conf) SetUpdateEvery(u int) {
	c.UpdEvery = u
}
