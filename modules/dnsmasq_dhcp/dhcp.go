package dnsmasq_dhcp

import (
	"github.com/netdata/go-orchestrator/module"
	"github.com/netdata/go.d.plugin/modules/dnsmasq_dhcp/ip"
	"time"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("dnsmasq_dhcp", creator)
}

const (
	//defaultLeasesPath = "/var/lib/misc/dnsmasq.leases"
	defaultLeasesPath = "/home/ilyam/leases"
	defaultConfPath   = "/home/ilyam/dnsmasq.conf"
	defaultConfDir    = "/home/ilyam/dnsmasq.d"
)

// New creates DnsmasqDHCP with default values.
func New() *DnsmasqDHCP {
	config := Config{
		LeasesPath: defaultLeasesPath,
		ConfPath:   defaultConfPath,
		ConfDir:    defaultConfDir,
	}

	return &DnsmasqDHCP{
		Config: config,
	}
}

// Config is the DnsmasqDHCP module configuration.
type Config struct {
	LeasesPath string `yaml:"leases_path"`
	ConfPath   string `yaml:"conf_path"`
	ConfDir    string `yaml:"conf_dir"`
}

// DnsmasqDHCP DnsmasqDHCP module.
type DnsmasqDHCP struct {
	module.Base
	Config `yaml:",inline"`

	// leases db modification time
	modTime time.Time
	pools   []*ip.Pool
}

// Cleanup makes cleanup.
func (DnsmasqDHCP) Cleanup() {}

// Init makes initialization.
func (d *DnsmasqDHCP) Init() bool {
	ranges, err := d.findDHCPRanges()
	if err != nil {
		d.Error(err)
		return false
	}

	for _, r := range ranges {
		pool := ip.NewPool(r)
		if pool == nil {
			d.Errorf("failed to parse %s", r)
			return false
		}
		d.pools = append(d.pools, pool)
	}

	return len(d.pools) > 0
}

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
