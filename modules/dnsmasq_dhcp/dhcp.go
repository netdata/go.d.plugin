package dnsmasq_dhcp

import (
	"github.com/netdata/go-orchestrator/module"
	"time"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("dnsmasq_dhcp", creator)
}

const (
	defaultLeasesPath = "/var/lib/misc/dnsmasq.leases"
)

// New creates DnsmasqDHCP with default values.
func New() *DnsmasqDHCP {
	config := Config{
		LeasesPath: defaultLeasesPath,
	}

	return &DnsmasqDHCP{
		Config: config,
	}
}

type Pool struct {
	Name  string `yaml:"name"`
	Range string `yaml:"range"`
}

// Config is the DnsmasqDHCP module configuration.
type Config struct {
	LeasesPath string `yaml:"leases_path"`
	Pools      []Pool `yaml:"pools"`
}

// DnsmasqDHCP DnsmasqDHCP module.
type DnsmasqDHCP struct {
	module.Base
	Config `yaml:",inline"`

	modTime time.Time
	pools   []pool
}

// Cleanup makes cleanup.
func (DnsmasqDHCP) Cleanup() {}

// Init makes initialization.
func (DnsmasqDHCP) Init() bool { return true }

// Check makes check.
func (DnsmasqDHCP) Check() bool { return true }

// Charts creates Charts.
func (DnsmasqDHCP) Charts() *Charts { return charts.Copy() }

// Collect collects metrics.
func (d *DnsmasqDHCP) Collect() map[string]int64 {
	mx, err := d.collect()
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}
