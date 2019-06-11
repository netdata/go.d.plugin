package dnsmasq_dhcp

import (
	"time"

	"github.com/netdata/go-orchestrator/module"

	"github.com/netdata/go.d.plugin/modules/dnsmasq_dhcp/ip"
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
		mx:     make(map[string]int64),
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
	ranges  []*ip.Range
	mx      map[string]int64
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

	for _, raw := range ranges {
		r := ip.NewRange(raw)
		if r == nil {
			continue
		}
		d.ranges = append(d.ranges, r)
	}

	return len(d.ranges) > 0
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
