package unbound

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("unbound", creator)
}

func New() *Unbound {
	config := Config{
		Address: "192.168.88.223:8953",
		Timeout: web.Duration{Duration: time.Second * 2},
	}

	return &Unbound{
		Config:   config,
		charts:   charts.Copy(),
		curCache: newCollectCache(),
		cache:    newCollectCache(),
	}
}

type unboundClient interface {
	send(command string) ([]string, error)
}

type (
	Config struct {
		Address             string       `yaml:"address"`
		ConfPath            string       `yaml:"conf_path"`
		Timeout             web.Duration `yaml:"timeout"`
		web.ClientTLSConfig `yaml:",inline"`
	}
	Unbound struct {
		module.Base
		Config `yaml:",inline"`

		client   unboundClient
		cache    collectCache
		curCache collectCache

		cumulative   bool
		hasExtCharts bool

		charts *module.Charts
	}
)

func (Unbound) Cleanup() {}

func (u *Unbound) Init() bool {
	u.client = newClient(clientConfig{
		address: u.Address,
		timeout: u.Timeout.Duration,
		useTLS:  false,
		tlsConf: nil,
	})
	return true
}

func (u Unbound) Check() bool {
	return true
}

func (u Unbound) Charts() *module.Charts { return u.charts }

func (u *Unbound) Collect() map[string]int64 {
	mx, err := u.collect()

	if err != nil {
		u.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}
