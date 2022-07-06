// SPDX-License-Identifier: GPL-3.0-or-later

package k8s_kubelet

import (
	"io/ioutil"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	creator := module.Creator{
		Defaults: module.Defaults{
			// NETDATA_CHART_PRIO_CGROUPS_CONTAINERS        40000
			Priority: 39900,
		},
		Create: func() module.Module { return New() },
	}

	module.Register("k8s_kubelet", creator)
}

// New creates Kubelet with default values.
func New() *Kubelet {
	config := Config{
		HTTP: web.HTTP{
			Request: web.Request{
				URL:     "http://127.0.0.1:10255/metrics",
				Headers: make(map[string]string),
			},
			Client: web.Client{
				Timeout: web.Duration{Duration: time.Second},
			},
		},
		TokenPath: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	return &Kubelet{
		Config:             config,
		charts:             charts.Copy(),
		collectedVMPlugins: make(map[string]bool),
	}
}

type (
	Config struct {
		web.HTTP  `yaml:",inline"`
		TokenPath string `yaml:"token_path"`
	}

	Kubelet struct {
		module.Base
		Config `yaml:",inline"`

		prom   prometheus.Prometheus
		charts *Charts
		// volume_manager_total_volumes
		collectedVMPlugins map[string]bool
	}
)

// Cleanup makes cleanup.
func (Kubelet) Cleanup() {}

// Init makes initialization.
func (k *Kubelet) Init() bool {
	b, err := ioutil.ReadFile(k.TokenPath)
	if err != nil {
		k.Warningf("error on reading service account token from '%s': %v", k.TokenPath, err)
	} else {
		k.Request.Headers["Authorization"] = "Bearer " + string(b)
	}

	client, err := web.NewHTTPClient(k.Client)
	if err != nil {
		k.Errorf("error on creating http client: %v", err)
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
