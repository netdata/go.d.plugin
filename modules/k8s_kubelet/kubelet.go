package k8s_kubelet

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:10255/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("k8s_kubelet", creator)
}

// New creates Kubelet with default values.
func New() *Kubelet {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{URL: defaultURL},
			Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
		},
	}

	return &Kubelet{
		Config:                     config,
		charts:                     charts.Copy(),
		activeVolumeManagerPlugins: make(map[string]bool),
	}
}

// Config is the DockerEngine module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// Kubelet Kubelet module.
type Kubelet struct {
	module.Base
	Config `yaml:",inline"`

	prom                       prometheus.Prometheus
	charts                     *Charts
	activeVolumeManagerPlugins map[string]bool
}

// Cleanup makes cleanup.
func (Kubelet) Cleanup() {}

// Init makes initialization.
func (k *Kubelet) Init() bool {
	if k.URL == "" {
		k.Error("URL parameter is mandatory, please set")
		return false
	}

	client, err := web.NewHTTPClient(k.Client)
	if err != nil {
		k.Errorf("error on creating http client : %v", err)
		return false
	}

	k.prom = prometheus.New(client, k.Request)

	return true
}

// Check makes check.
func (k *Kubelet) Check() bool {
	return len(k.Collect()) > 0
}

// Charts creates Charts.
func (k Kubelet) Charts() *Charts {
	return k.charts
}

// Collect collects mx.
func (k *Kubelet) Collect() map[string]int64 {
	mx, err := k.collect()

	if err != nil {
		k.Error(err)
		return nil
	}

	return mx
}
