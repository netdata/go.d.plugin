package du

import (
	"os"
	"path/filepath"

	"github.com/netdata/go.d.plugin/agent/module"
)

// Config is configuration for Du module.
type Config struct {
	Paths []string `yaml:"paths"`
}

// Du is a module for file/folder size monitoring like linux command 'du'
type Du struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`

	charts        *module.Charts
	collectedDims map[string]bool
}

func init() {
	module.Register("du", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery:        module.UpdateEvery,
			AutoDetectionRetry: module.AutoDetectionRetry,
			Priority:           module.Priority,
			Disabled:           true,
		},
		Create: func() module.Module { return New() },
	})
}

func fileSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// New creates Du with default values.
func New() *Du {
	return &Du{
		Config: Config{
			Paths: make([]string, 0),
		},
		collectedDims: make(map[string]bool),
	}
}

// Init makes initialization.
func (du *Du) Init() bool {
	err := du.validateConfig()
	if err != nil {
		du.Errorf("config validation: %v", err)
		return false
	}

	charts, err := du.initCharts()
	if err != nil {
		du.Errorf("charts init: %v", err)
		return false
	}
	du.charts = charts
	return true
}

// Check makes check.
func (du *Du) Check() bool {
	return len(du.Collect()) > 0
}

// Charts creates Charts.
func (du *Du) Charts() *module.Charts {
	return du.charts
}

// Collect collects metrics.
func (du *Du) Collect() map[string]int64 {
	mx, err := du.collect()
	if err != nil {
		du.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

// Cleanup makes cleanup.
func (Du) Cleanup() {}
