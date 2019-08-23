package zookeeper

import (
	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("zookeeper", creator)
}

// Config is the Zookeeper module configuration.
type Config struct {
	Address string
	Timeout web.Duration `yaml:"timeout"`
}

// New creates Zookeeper with default values.
func New() *Zookeeper {
	return &Zookeeper{}
}

type zookeeperFetcher interface {
	fetch(command string) ([]string, error)
}

// Zookeeper Zookeeper module.
type Zookeeper struct {
	module.Base
	zookeeperFetcher
}

// Cleanup makes cleanup.
func (Zookeeper) Cleanup() {}

// Init makes initialization.
func (Zookeeper) Init() bool {
	return true
}

// Check makes check.
func (Zookeeper) Check() bool {
	return true
}

// Charts creates Charts.
func (Zookeeper) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics
func (z *Zookeeper) Collect() map[string]int64 {
	mx, err := z.collect()
	if err != nil {
		z.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}
