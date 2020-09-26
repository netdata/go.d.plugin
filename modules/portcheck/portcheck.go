package portcheck

import (
	"net"
	"sort"
	"time"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("portcheck", creator)
}

const defaultConnectTimeout = time.Second * 2

// New creates PortCheck with default values.
func New() *PortCheck {
	config := Config{
		Timeout: web.Duration{Duration: defaultConnectTimeout},
	}

	return &PortCheck{
		Config: config,
		dial:   net.DialTimeout,
	}
}

/// Config is the Portcheck module configuration file.
type Config struct {
	Host    string       `yaml:"host"`
	Ports   []int        `yaml:"ports"`
	Timeout web.Duration `yaml:"timeout"`
}

type dialFunc func(network, address string, timeout time.Duration) (net.Conn, error)

type state string

const (
	success state = "success"
	timeout state = "timeout"
	failed  state = "failed"
)

type port struct {
	number  int
	state   state
	inState int
	latency int
}

// PortCheck portcheck module.
type PortCheck struct {
	module.Base
	Config      `yaml:",inline"`
	UpdateEvery int `yaml:"update_every"`
	dial        dialFunc
	ports       []*port
}

// Cleanup makes cleanup.
func (PortCheck) Cleanup() {}

// Init makes initialization.
func (pc *PortCheck) Init() bool {
	if pc.Host == "" {
		pc.Error("host parameter is not set")
		return false
	}

	if len(pc.Ports) == 0 {
		pc.Error("ports parameter is not set")
		return false
	}

	sort.Ints(pc.Ports)

	for _, p := range pc.Ports {
		pc.ports = append(pc.ports, &port{number: p})
	}

	pc.Debugf("using host %s", pc.Host)
	pc.Debugf("using ports %v", pc.Ports)
	pc.Debugf("using TCP connection timeout: %s", pc.Timeout)

	return true
}

// Check makes check.
func (PortCheck) Check() bool { return true }

// Charts creates charts.
func (pc PortCheck) Charts() *Charts {
	charts := &Charts{}

	for _, port := range pc.Ports {
		_ = charts.Add(*newPortCharts(port)...)
	}

	return charts
}

// Collect collects metrics.
func (pc *PortCheck) Collect() map[string]int64 {
	mx, err := pc.collect()

	if err != nil {
		pc.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}
