package k8s_kubeproxy

import (
	"math"

	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go-orchestrator/module"
)

func (kp *KubeProxy) collect() (map[string]int64, error) {
	raw, err := kp.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	kp.collectSyncProxyRules(raw, mx)
	kp.collectRESTClientHTTPRequests(raw, mx)
	kp.collectHTTPRequestDuration(raw, mx)

	return stm.ToMap(mx), nil
}

func (kp *KubeProxy) collectSyncProxyRules(raw prometheus.Metrics, mx *metrics) {
	m := raw.FindByName("kubeproxy_sync_proxy_rules_latency_microseconds_count")
	mx.SyncProxyRules.Count.Set(m.Max())
	kp.collectSyncProxyRulesLatency(raw, mx)
}

func (kp *KubeProxy) collectSyncProxyRulesLatency(raw prometheus.Metrics, mx *metrics) {
	metricName := "kubeproxy_sync_proxy_rules_latency_microseconds_bucket"

	for _, metric := range raw.FindByName(metricName) {
		bucket := metric.Labels.Get("le")
		switch bucket {
		case "1000":
			mx.SyncProxyRules.Latency.LE1000.Set(metric.Value)
		case "2000":
			mx.SyncProxyRules.Latency.LE2000.Set(metric.Value)
		case "4000":
			mx.SyncProxyRules.Latency.LE4000.Set(metric.Value)
		case "8000":
			mx.SyncProxyRules.Latency.LE8000.Set(metric.Value)
		case "16000":
			mx.SyncProxyRules.Latency.LE16000.Set(metric.Value)
		case "32000":
			mx.SyncProxyRules.Latency.LE32000.Set(metric.Value)
		case "64000":
			mx.SyncProxyRules.Latency.LE64000.Set(metric.Value)
		case "128000":
			mx.SyncProxyRules.Latency.LE128000.Set(metric.Value)
		case "256000":
			mx.SyncProxyRules.Latency.LE256000.Set(metric.Value)
		case "512000":
			mx.SyncProxyRules.Latency.LE512000.Set(metric.Value)
		case "1.024e+06":
			mx.SyncProxyRules.Latency.LE1024000.Set(metric.Value)
		case "2.048e+06":
			mx.SyncProxyRules.Latency.LE2048000.Set(metric.Value)
		case "4.096e+06":
			mx.SyncProxyRules.Latency.LE4096000.Set(metric.Value)
		case "8.192e+06":
			mx.SyncProxyRules.Latency.LE8192000.Set(metric.Value)
		case "1.6384e+07":
			mx.SyncProxyRules.Latency.LE16384000.Set(metric.Value)
		case "+Inf":
			mx.SyncProxyRules.Latency.Inf.Set(metric.Value)
		}
	}
}

func (kp *KubeProxy) collectRESTClientHTTPRequests(raw prometheus.Metrics, mx *metrics) {
	metricName := "rest_client_requests_total"
	chart := kp.charts.Get("rest_client_requests_by_code")

	for _, metric := range raw.FindByName(metricName) {
		code := metric.Labels.Get("code")
		if code == "" {
			continue
		}
		dimID := "rest_client_requests_" + code
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: code, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.RESTClient.Requests.ByStatusCode[code] = mtx.Gauge(metric.Value)
	}

	chart = kp.charts.Get("rest_client_requests_by_method")

	for _, metric := range raw.FindByName(metricName) {
		method := metric.Labels.Get("method")
		if method == "" {
			continue
		}
		dimID := "rest_client_requests_" + method
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: method, Algo: module.Incremental})
			chart.MarkNotCreated()
		}
		mx.RESTClient.Requests.ByMethod[method] = mtx.Gauge(metric.Value)
	}
}

func (kp *KubeProxy) collectHTTPRequestDuration(raw prometheus.Metrics, mx *metrics) {
	// Summary
	for _, metric := range raw.FindByName("http_request_duration_microseconds") {
		if math.IsNaN(metric.Value) {
			continue
		}
		quantile := metric.Labels.Get("quantile")
		switch quantile {
		case "0.5":
			mx.HTTP.Request.Duration.Quantile05.Set(metric.Value)
		case "0.9":
			mx.HTTP.Request.Duration.Quantile09.Set(metric.Value)
		case "0.99":
			mx.HTTP.Request.Duration.Quantile099.Set(metric.Value)
		}
	}
}
