package vsphere

import (
	"net/url"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/modules/vsphere/discover"
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
	"github.com/netdata/go.d.plugin/modules/vsphere/scrape"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/vmware/govmomi/performance"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			UpdateEvery: 5,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("vsphere", creator)
}

const (
	defaultURL              = "http://127.0.0.1"
	defaultHTTPTimeout      = time.Second * 5
	defaultDiscoverInterval = time.Minute * 5

	vCenterURL = "https://192.168.0.154/sdk"
	username   = "administrator@vsphere.local"
	password   = "123qwe!@#QWE"
	timeout    = time.Second * 10
)

type discoverer interface {
	Discover() (*rs.Resources, error)
}

type metricScraper interface {
	ScrapeHostsMetrics(rs.Hosts) []performance.EntityMetric
	ScrapeVMsMetrics(rs.VMs) []performance.EntityMetric
}

func New() *VSphere {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: defaultURL,
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
		DiscoverInterval: web.Duration{Duration: defaultDiscoverInterval},
	}

	return &VSphere{
		collectionLock:  new(sync.RWMutex),
		Config:          config,
		charts:          &Charts{},
		discoveredHosts: make(map[string]int),
		discoveredVMs:   make(map[string]int),
		charted:         make(map[string]bool),
	}
}

type Config struct {
	web.HTTP         `yaml:",inline"`
	DiscoverInterval web.Duration  `yaml:"discover_interval"`
	HostsInclude     []hostInclude `yaml:"host_include"`
	VMsInclude       []vmInclude   `yaml:"vm_include"`
}

type VSphere struct {
	module.Base
	Config `yaml:",inline"`

	discoverer
	metricScraper

	collectionLock  *sync.RWMutex
	resources       *rs.Resources
	discoveryTask   *task
	discoveredHosts map[string]int
	discoveredVMs   map[string]int
	charted         map[string]bool
	charts          *Charts
}

func (vs VSphere) Cleanup() {
	if vs.discoveryTask != nil {
		return
	}
	vs.discoveryTask.stop()
}

func (vs *VSphere) Init() bool {
	u, err := url.Parse(vCenterURL)
	if err != nil {
		vs.Error(err)
		return false
	}

	c, err := client.New(client.Config{
		URL:             u.String(),
		User:            username,
		Password:        password,
		Timeout:         timeout,
		ClientTLSConfig: web.ClientTLSConfig{InsecureSkipVerify: true},
	})
	if err != nil {
		vs.Error(err)
		return false
	}

	hm, err := parseHostIncludes(vs.HostsInclude)
	if err != nil {
		vs.Error(err)
		return false
	}

	vmm, err := parseVMIncludes(vs.VMsInclude)
	if err != nil {
		vs.Error(err)
		return false
	}

	cl := discover.NewVSphereDiscoverer(c)
	if hm != nil {
		cl.HostMatcher = hm
	}
	if vmm != nil {
		cl.VMMatcher = vmm
	}
	mc := scrape.NewVSphereMetricScraper(c)

	res, err := cl.Discover()
	if err != nil {
		vs.Error(err)
		return false
	}
	vs.resources = res
	vs.metricScraper = mc

	return true
}

func (vs VSphere) Check() bool {
	return true
}

func (vs VSphere) Charts() *module.Charts { return vs.charts }

func (vs *VSphere) Collect() map[string]int64 {
	mx, err := vs.collect()

	if err != nil {
		vs.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}

	return mx
}
