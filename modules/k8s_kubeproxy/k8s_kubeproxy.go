package k8s_kubeproxy

import (
	"strings"
	"time"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
)

const (
	defaultURL         = "http://127.0.0.1:10249/metrics"
	defaultHTTPTimeout = time.Second * 2
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}

	module.Register("k8s_kubeproxy", creator)
}

// New creates KubeProxy with default values.
func New() *KubeProxy {
	return &KubeProxy{
		Config: Config{
			HTTP: web.HTTP{
				Request: web.Request{URL: defaultURL},
				Client:  web.Client{Timeout: web.Duration{Duration: defaultHTTPTimeout}},
			},
		},
		mx:     newMetrics(),
		charts: charts.Copy(),
	}
}

func newMetrics() *metrics {
	return &metrics{SyncProxyRulesLatency: make(map[string]mtx.Gauge)}
}

type metrics struct {
	SyncProxyRulesLatency map[string]mtx.Gauge `stm:"sync_proxy_rules_latency_microseconds_bucket"`
}

// Config is the KubeProxy module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// DockerEngine DockerEngine module.
type KubeProxy struct {
	module.Base
	Config `yaml:",inline"`

	prom   prometheus.Prometheus
	charts *Charts
	mx     *metrics
}

// Cleanup makes cleanup.
func (KubeProxy) Cleanup() {}

// Init makes initialization.
func (kp *KubeProxy) Init() bool {
	if kp.URL == "" {
		kp.Error("URL parameter is mandatory, please set")
		return false
	}

	client, err := web.NewHTTPClient(kp.Client)
	if err != nil {
		kp.Errorf("error on creating http client : %v", err)
		return false
	}

	kp.prom = prometheus.New(client, kp.Request)

	return true
}

// Check makes check.
func (kp *KubeProxy) Check() bool {
	return len(kp.Collect()) > 0
}

// Charts creates Charts.
func (kp KubeProxy) Charts() *Charts {
	return kp.charts
}

// Collect collects metrics.
func (kp *KubeProxy) Collect() map[string]int64 {
	raw, err := kp.prom.Scrape()

	if err != nil {
		kp.Error(err)
		return nil
	}

	kp.collectSyncProxyRuleLatency(raw)

	return stm.ToMap(kp.mx)
}

func (kp *KubeProxy) collectSyncProxyRuleLatency(raw prometheus.Metrics) {
	metricName := "kubeproxy_sync_proxy_rules_latency_microseconds_bucket"
	for _, metric := range raw.FindByName(metricName) {
		val := metric.Labels.Get("le")
		if val == "" {
			continue
		}
		val = strings.Replace(val, ".", "_", -1)
		kp.mx.SyncProxyRulesLatency[val] = mtx.Gauge(metric.Value)
	}
}
