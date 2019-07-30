package vsphere

import (
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/modules/vsphere/discover"
	"github.com/netdata/go.d.plugin/modules/vsphere/match"
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
	defaultURL               = "https://192.168.0.154"
	defaultHTTPTimeout       = time.Second * 10
	defaultDiscoveryInterval = time.Minute * 5
	username                 = "administrator@vsphere.local"
	password                 = "123qwe!@#QWE"
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
				UserURL:  defaultURL,
				Username: username,
				Password: password,
			},
			Client: web.Client{
				Timeout:         web.Duration{Duration: defaultHTTPTimeout},
				ClientTLSConfig: web.ClientTLSConfig{InsecureSkipVerify: true},
			},
		},
		DiscoveryInterval: web.Duration{Duration: defaultDiscoveryInterval},
		HostsInclude:      []string{"/*"},
		VMsInclude:        []string{"/*"},
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
	web.HTTP          `yaml:",inline"`
	DiscoveryInterval web.Duration       `yaml:"discovery_interval"`
	HostsInclude      match.HostIncludes `yaml:"host_include"`
	VMsInclude        match.VMIncludes   `yaml:"vm_include"`
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

func (vs *VSphere) Cleanup() {
	if vs.discoveryTask == nil {
		return
	}
	vs.discoveryTask.stop()
}

func (vs VSphere) createVSphereClient() (*client.Client, error) {
	config := client.Config{
		URL:             vs.UserURL,
		User:            vs.Username,
		Password:        vs.Password,
		Timeout:         vs.Timeout.Duration,
		ClientTLSConfig: vs.ClientTLSConfig,
	}
	return client.New(config)
}

func (vs *VSphere) createVSphereDiscoverer(c *client.Client) error {
	d := discover.NewVSphereDiscoverer(c)
	d.Logger = vs.Logger

	hm, err := vs.HostsInclude.Parse()
	if err != nil {
		return err
	}
	if hm != nil {
		d.HostMatcher = hm
	}
	vmm, err := vs.VMsInclude.Parse()
	if err != nil {
		return err
	}
	if vmm != nil {
		d.VMMatcher = vmm
	}

	vs.discoverer = d
	return nil
}

func (vs *VSphere) createVSphereMetricScraper(c *client.Client) {
	ms := scrape.NewVSphereMetricScraper(c)
	ms.Logger = vs.Logger
	vs.metricScraper = ms
}

func (vs *VSphere) Init() bool {
	c, err := vs.createVSphereClient()
	if err != nil {
		vs.Errorf("error on creating vsphere client : %v", err)
		return false
	}

	err = vs.createVSphereDiscoverer(c)
	if err != nil {
		vs.Errorf("error on creating vsphere discoverer : %v", err)
		return false
	}

	vs.createVSphereMetricScraper(c)

	vs.discoverOnce()
	vs.goDiscovery()
	return true
}

func (vs VSphere) Check() bool {
	return true
}

func (vs VSphere) Charts() *Charts {
	return vs.charts
}

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
