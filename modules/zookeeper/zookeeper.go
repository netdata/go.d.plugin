package zookeeper

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("zookeeper", creator)
}

// Config is the Zookeeper module configuration.
type Config struct {
	Address          string
	Timeout          web.Duration `yaml:"timeout"`
	UseTLS           bool         `yaml:"use_tls"`
	tlscfg.TLSConfig `yaml:",inline"`
}

// New creates Zookeeper with default values.
func New() *Zookeeper {
	config := Config{
		Address: "127.0.0.1:2181",
		Timeout: web.Duration{Duration: time.Second},
		UseTLS:  false,
	}
	return &Zookeeper{Config: config}
}

type zookeeperFetcher interface {
	fetch(command string) ([]string, error)
}

// Zookeeper Zookeeper module.
type Zookeeper struct {
	module.Base
	zookeeperFetcher
	Config `yaml:",inline"`
}

// Cleanup makes cleanup.
func (Zookeeper) Cleanup() {}

func (z *Zookeeper) createZookeeperFetcher() (err error) {
	var tlsConf *tls.Config
	if z.UseTLS {
		tlsConf, err = tlscfg.NewTLSConfig(z.TLSConfig)
		if err != nil {
			return fmt.Errorf("error on creating tls config : %v", err)
		}
	}

	conf := clientConfig{
		network: "tcp",
		address: z.Address,
		timeout: z.Timeout.Duration,
		useTLS:  z.UseTLS,
		tlsConf: tlsConf,
	}
	z.zookeeperFetcher = newClient(conf)
	return nil
}

// Init makes initialization.
func (z *Zookeeper) Init() bool {
	err := z.createZookeeperFetcher()
	if err != nil {
		z.Error(err)
		return false
	}

	return true
}

// Check makes check.
func (z *Zookeeper) Check() bool {
	return len(z.Collect()) > 0
}

// Charts creates Charts.
func (Zookeeper) Charts() *Charts {
	return charts.Copy()
}

// Collect collects metrics.
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
