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
		// "/etc/unbound/unbound.conf"
		Address:  "192.168.88.223:8953",
		ConfPath: "/Users/ilyam/Projects/goland/go.d.plugin/modules/unbound/testdata/unbound.conf",
		Timeout:  web.Duration{Duration: time.Second * 2},
	}

	return &Unbound{
		Config:   config,
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
		DisableTLS          bool         `yaml:"disable_tls"`
		Cumulative          bool         `yaml:"cumulative_stats"`
		web.ClientTLSConfig `yaml:",inline"`
	}
	Unbound struct {
		module.Base
		Config `yaml:",inline"`

		client   unboundClient
		cache    collectCache
		curCache collectCache

		// used in cumulative mode
		prevCacheMiss  float64
		prevRecReplies float64

		hasExtCharts bool

		charts *module.Charts
	}
)

func (Unbound) Cleanup() {}

func (u *Unbound) Init() bool {
	if !u.initConfig() {
		return false
	}

	if err := u.initClient(); err != nil {
		u.Errorf("creating client: %v", err)
		return false
	}
	u.charts = charts(u.Cumulative)
	return true
}

func (u *Unbound) Check() bool {
	return len(u.Collect()) > 0
}

func (u Unbound) Charts() *Charts {
	return u.charts
}

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
