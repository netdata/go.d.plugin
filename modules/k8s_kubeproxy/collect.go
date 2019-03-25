package k8s_kubeproxy

import (
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

	kp.collectSyncProxyRules(raw)
	kp.collectRESTClientHTTPRequests(raw)

	return stm.ToMap(kp.mx), nil
}

func (kp *KubeProxy) collectSyncProxyRules(raw prometheus.Metrics) {
	val := raw.FindByName("kubeproxy_sync_proxy_rules_latency_microseconds_count").Max()
	kp.mx.SyncProxyRules.Count.Set(val)
	kp.collectSyncProxyRulesLatency(raw)
}

func (kp *KubeProxy) collectSyncProxyRulesLatency(raw prometheus.Metrics) {
	metricName := "kubeproxy_sync_proxy_rules_latency_microseconds_bucket"

	for _, metric := range raw.FindByName(metricName) {
		value := metric.Labels.Get("le")
		switch value {
		case "1000":
			kp.mx.SyncProxyRules.Latency.LE1000.Set(metric.Value)
		case "2000":
			kp.mx.SyncProxyRules.Latency.LE2000.Set(metric.Value)
		case "4000":
			kp.mx.SyncProxyRules.Latency.LE4000.Set(metric.Value)
		case "8000":
			kp.mx.SyncProxyRules.Latency.LE8000.Set(metric.Value)
		case "16000":
			kp.mx.SyncProxyRules.Latency.LE16000.Set(metric.Value)
		case "32000":
			kp.mx.SyncProxyRules.Latency.LE32000.Set(metric.Value)
		case "64000":
			kp.mx.SyncProxyRules.Latency.LE64000.Set(metric.Value)
		case "128000":
			kp.mx.SyncProxyRules.Latency.LE128000.Set(metric.Value)
		case "256000":
			kp.mx.SyncProxyRules.Latency.LE256000.Set(metric.Value)
		case "512000":
			kp.mx.SyncProxyRules.Latency.LE512000.Set(metric.Value)
		case "1.024e+06":
			kp.mx.SyncProxyRules.Latency.LE1024000.Set(metric.Value)
		case "2.048e+06":
			kp.mx.SyncProxyRules.Latency.LE2048000.Set(metric.Value)
		case "4.096e+06":
			kp.mx.SyncProxyRules.Latency.LE4096000.Set(metric.Value)
		case "8.192e+06":
			kp.mx.SyncProxyRules.Latency.LE8192000.Set(metric.Value)
		case "1.6384e+07":
			kp.mx.SyncProxyRules.Latency.LE16384000.Set(metric.Value)
		case "+Inf":
			kp.mx.SyncProxyRules.Latency.Inf.Set(metric.Value)
		}
	}
}

func (kp *KubeProxy) collectRESTClientHTTPRequests(raw prometheus.Metrics) {
	metricName := "rest_client_requests_total"

	for _, metric := range raw.FindByName(metricName) {
		value := metric.Labels.Get("code")
		if value == "" {
			continue
		}
		m := kp.mx.RESTClient.HTTPRequests.ByStatusCode

		if _, ok := m[value]; !ok {
			chart := kp.charts.Get("rest_client_requests_by_code")
			_ = chart.AddDim(&Dim{
				ID:   "rest_client_http_requests_" + value,
				Name: value,
				Algo: module.Incremental,
			})
			chart.MarkNotCreated()
		}
		m[value] = mtx.Gauge(metric.Value)
	}

	for _, metric := range raw.FindByName(metricName) {
		value := metric.Labels.Get("method")
		if value == "" {
			continue
		}
		m := kp.mx.RESTClient.HTTPRequests.ByMethod

		if _, ok := m[value]; !ok {
			chart := kp.charts.Get("rest_client_requests_by_method")
			_ = chart.AddDim(&Dim{
				ID:   "rest_client_http_requests_" + value,
				Name: value,
				Algo: module.Incremental,
			})
			chart.MarkNotCreated()
		}
		m[value] = mtx.Gauge(metric.Value)
	}
}
