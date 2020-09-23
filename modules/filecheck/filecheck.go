package filecheck

import (
	"github.com/netdata/go-orchestrator/module"
)

func init() {
	module.Register("filecheck", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 10,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Filecheck {
	return &Filecheck{
		Config: Config{
			Files: filesConfig{},
			Dirs:  dirsConfig{},
		},
		collectedFiles: make(map[string]bool),
		collectedDirs:  make(map[string]bool),
	}
}

type (
	Config struct {
		Files filesConfig `yaml:"files"`
		Dirs  dirsConfig  `yaml:"dirs"`
	}
	filesConfig struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	}
	dirsConfig struct {
		Include []string `yaml:"include"`
		Exclude []string `yaml:"exclude"`
	}
)

type Filecheck struct {
	module.Base
	Config `yaml:",inline"`

	collectedFiles map[string]bool
	collectedDirs  map[string]bool
	charts         *module.Charts
}

func (Filecheck) Cleanup() {
}

func (fc *Filecheck) Init() bool {
	err := fc.validateConfig()
	if err != nil {
		fc.Errorf("error on validating config: %v", err)
		return false
	}

	charts, err := fc.initCharts()
	if err != nil {
		fc.Errorf("error on charts initialization: %v", err)
		return false
	}
	fc.charts = charts

	fc.Debugf("monitored files: %v", fc.Files.Include)
	fc.Debugf("monitored dirs: %v", fc.Dirs.Include)
	return true
}

func (fc *Filecheck) Check() bool {
	return len(fc.Collect()) > 0
}

func (fc *Filecheck) Charts() *module.Charts {
	return fc.charts
}

func (fc *Filecheck) Collect() map[string]int64 {
	mx, err := fc.collect()
	if err != nil {
		fc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
