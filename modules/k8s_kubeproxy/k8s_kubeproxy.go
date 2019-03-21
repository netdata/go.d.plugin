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
	return &metrics{
		SyncProxyRules: syncProxyRulesMetrics{
			Latency: make(map[string]mtx.Gauge),
		},
		RestClientMetrics: restClientMetrics{
			HTTPRequestsByStatusCode: make(map[string]mtx.Gauge),
			HTTPRequestsByMethod:     make(map[string]mtx.Gauge),
		},
	}
}

type metrics struct {
	SyncProxyRules    syncProxyRulesMetrics `stm:"sync_proxy_rules"`
	RestClientMetrics restClientMetrics     `stm:"rest_client"`
}

type syncProxyRulesMetrics struct {
	Count   mtx.Gauge            `stm:"count"`
	Latency map[string]mtx.Gauge `stm:"bucket"`
}

type restClientMetrics struct {
	HTTPRequestsByStatusCode map[string]mtx.Gauge `stm:"requests"`
	HTTPRequestsByMethod     map[string]mtx.Gauge `stm:"requests"`
}

// Config is the KubeProxy module configuration.
type Config struct {
	web.HTTP `yaml:",inline"`
}

// KubeProxy is KubeProxy module.
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

	kp.mx.SyncProxyRules.Count.Set(
		raw.FindByName("kubeproxy_sync_proxy_rules_latency_microseconds_count").Max())
	kp.collectSyncProxyRuleLatency(raw)
	kp.collectRestClientHTTPRequests(raw)

	return stm.ToMap(kp.mx)
}

func (kp *KubeProxy) collectSyncProxyRuleLatency(raw prometheus.Metrics) {
	metricName := "kubeproxy_sync_proxy_rules_latency_microseconds_bucket"

	for _, metric := range raw.FindByName(metricName) {
		val := metric.Labels.Get("le")
		if val == "" {
			continue
		}
		// TODO: FIX
		newVal := strings.Replace(val, ".", "_", -1)
		kp.mx.SyncProxyRules.Latency[newVal] = mtx.Gauge(metric.Value)
	}
}

func (kp *KubeProxy) collectRestClientHTTPRequests(raw prometheus.Metrics) {
	metricName := "rest_client_requests_total"

	for _, metric := range raw.FindByName(metricName) {
		value := metric.Labels.Get("code")
		m := kp.mx.RestClientMetrics.HTTPRequestsByStatusCode

		if value != "" {
			if _, ok := m[value]; !ok {
				chart := kp.charts.Get("rest_client_requests_by_code")
				_ = chart.AddDim(&Dim{
					ID:   "rest_client_requests_" + value,
					Name: value,
					Algo: module.Incremental,
				})
				chart.MarkNotCreated()
			}
			m[value] = mtx.Gauge(metric.Value)
		}

		value = metric.Labels.Get("method")
		m = kp.mx.RestClientMetrics.HTTPRequestsByMethod

		if value != "" {
			if _, ok := m[value]; !ok {
				chart := kp.charts.Get("rest_client_requests_by_method")
				_ = chart.AddDim(&Dim{
					ID:   "rest_client_requests_" + value,
					Name: value,
					Algo: module.Incremental,
				})
				chart.MarkNotCreated()
			}
			m[value] = mtx.Gauge(metric.Value)
		}
	}
}
