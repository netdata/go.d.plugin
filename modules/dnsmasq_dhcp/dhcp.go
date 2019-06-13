package dnsmasq_dhcp

import (
	"fmt"
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
	defaultLeasesPath = "/var/lib/misc/dnsmasq/dnsmasq.leases"
	defaultConfPath   = "/etc/dnsmasq.conf"
	defaultConfDir    = "/etc/dnsmasq.d"
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

func (c Config) String() string {
	return fmt.Sprintf("leases_path: [%s], conf_path: [%s], conf_dir: [%s]",
		c.LeasesPath,
		c.ConfPath,
		c.ConfDir,
	)
}

// DnsmasqDHCP DnsmasqDHCP module.
type DnsmasqDHCP struct {
	module.Base
	Config `yaml:",inline"`

	// dnsmasq.leases db modification time
	modTime time.Time
	ranges  []ip.IRange
	mx      map[string]int64
}

// Cleanup makes cleanup.
func (DnsmasqDHCP) Cleanup() {}

// Init makes initialization.
func (d *DnsmasqDHCP) Init() bool {
	d.Infof("start config : %s", d.Config)

	ranges, err := d.findDHCPRanges()
	if err != nil {
		d.Error(err)
		return false
	}

	for _, raw := range ranges {
		r := ip.ParseRange(raw)
		if r == nil {
			d.Warningf("error on parsing '%s' dhcp range, skipping it", raw)
			continue
		}
		d.ranges = append(d.ranges, r)
	}
	return len(d.ranges) > 0
}

// Check makes check.
func (d *DnsmasqDHCP) Check() bool { return len(d.Collect()) > 0 }

// Charts creates Charts.
func (d DnsmasqDHCP) Charts() *Charts { return d.charts() }

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
