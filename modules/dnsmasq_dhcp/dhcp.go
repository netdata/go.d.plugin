package dnsmasq_dhcp

import (
	"fmt"
	"net"
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

// New creates DnsmasqDHCP with default values.
func New() *DnsmasqDHCP {
	config := Config{
		// debian defaults
		LeasesPath: "/var/lib/misc/dnsmasq.leases",
		ConfPath:   "/etc/dnsmasq.conf",
		ConfDir:    "/etc/dnsmasq.d",
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

	leasesModTime time.Time
	ranges        []ip.IRange
	staticIPs     []net.IP
	mx            map[string]int64
}

// Cleanup makes cleanup.
func (DnsmasqDHCP) Cleanup() {}

// Init makes initialization.
func (d *DnsmasqDHCP) Init() bool {
	d.Infof("start config : %s", d.Config)
	err := d.autodetection()
	if err != nil {
		d.Error(err)
		return false
	}
	return true
}

// Check makes check.
func (d *DnsmasqDHCP) Check() bool {
	return len(d.Collect()) > 0
}

// Charts creates Charts.
func (d DnsmasqDHCP) Charts() *Charts {
	return d.charts()
}

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
