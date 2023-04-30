package cgminer

import (
	"math/rand"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("cgminer", module.Creator{
		Defaults: module.Defaults{
			UpdateEvery:        module.UpdateEvery,
			AutoDetectionRetry: module.AutoDetectionRetry,
			Priority:           module.Priority,
			Disabled:           true,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *Cgminer {
	return &Cgminer{
		Config: Config{
			Charts: ConfigCharts{
				Num:  1,
				Dims: 4,
			},
			HiddenCharts: ConfigCharts{
				Num:  0,
				Dims: 4,
			},
		},

		randInt:       func() int64 { return rand.Int63n(100) },
		collectedDims: make(map[string]bool),
	}
}

type (
	Config struct {
		Charts       ConfigCharts `yaml:"charts"`
		HiddenCharts ConfigCharts `yaml:"hidden_charts"`
	}
	ConfigCharts struct {
		Type     string `yaml:"type"`
		Num      int    `yaml:"num"`
		Contexts int    `yaml:"contexts"`
		Dims     int    `yaml:"dimensions"`
		Labels   int    `yaml:"labels"`
	}
)

type Cgminer struct {
	module.Base // should be embedded by every module
	Config      `yaml:",inline"`

	randInt       func() int64
	charts        *module.Charts
	collectedDims map[string]bool
}

func (c *Cgminer) Init() bool {
	err := c.validateConfig()
	if err != nil {
		c.Errorf("config validation: %v", err)
		return false
	}

	charts, err := c.initCharts()
	if err != nil {
		c.Errorf("charts init: %v", err)
		return false
	}
	c.charts = charts
	return true
}

func (c *Cgminer) Check() bool {
	return len(c.Collect()) > 0
}

func (c *Cgminer) Charts() *module.Charts {
	return c.charts
}

func (c *Cgminer) Collect() map[string]int64 {
	mx, err := c.collect()
	if err != nil {
		c.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (c *Cgminer) Cleanup() {}
