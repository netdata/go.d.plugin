package vsphere

import (
	"net/url"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/modules/vsphere/client"
	"github.com/netdata/go.d.plugin/modules/vsphere/collect"
	"github.com/netdata/go.d.plugin/modules/vsphere/discover"
	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"
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
	defaultURL         = "http://127.0.0.1"
	defaultHTTPTimeout = time.Second * 5

	vCenterURL = "https://192.168.0.154/sdk"
	username   = "administrator@vsphere.local"
	password   = "123qwe!@#QWE"
	timeout    = time.Second * 10
)

type discoverer interface {
	Discover() (*rs.Resources, error)
}

type metricCollector interface {
	CollectHostsMetrics(rs.Hosts) []performance.EntityMetric
	CollectVMsMetrics(rs.VMs) []performance.EntityMetric
}

// New creates VSphere with default values.
func New() *VSphere {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				UserURL: defaultURL,
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &VSphere{
		resLock:        new(sync.RWMutex),
		Config:         config,
		charts:         &Charts{},
		collectedHosts: make(map[string]int),
		collectedVMs:   make(map[string]int),
		charted:        make(map[string]bool),
	}
}

// Config is the VSphere module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// VSphere VSphere module.
type VSphere struct {
	module.Base
	Config `yaml:",inline"`

	resLock   *sync.RWMutex
	resources *rs.Resources
	discoverer
	metricCollector
	discoveryTask *task

	charts *module.Charts

	collectedHosts map[string]int
	collectedVMs   map[string]int
	charted        map[string]bool
}

// Cleanup makes cleanup.
func (vs VSphere) Cleanup() {
	if vs.discoveryTask != nil {
		return
	}
	vs.discoveryTask.stop()
}

// Init makes initialization.
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

	cl := discover.NewVSphereDiscoverer(c)
	mc := collect.NewVSphereMetricCollector(c)

	res, err := cl.Discover()
	if err != nil {
		vs.Error(err)
		return false
	}
	vs.resources = res
	vs.metricCollector = mc

	return true
}

// Check makes check.
func (vs VSphere) Check() bool {
	return true

}

// Charts returns Charts.
func (vs VSphere) Charts() *module.Charts { return vs.charts }

// Collect collects metricList.
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
