package k8s_kubeproxy

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go-orchestrator/module"
)

func newMetrics() *metrics {
	return &metrics{
		RestClientMetrics: restClientMetrics{
			HTTPRequestsByStatusCode: make(map[string]mtx.Gauge),
			HTTPRequestsByMethod:     make(map[string]mtx.Gauge),
		},
	}
}

type metrics struct {
	SyncProxyRules struct {
		Count   mtx.Gauge `stm:"count"`
		Latency struct {
			LE1000     mtx.Gauge `stm:"1000"`
			LE2000     mtx.Gauge `stm:"2000"`
			LE4000     mtx.Gauge `stm:"4000"`
			LE8000     mtx.Gauge `stm:"8000"`
			LE16000    mtx.Gauge `stm:"16000"`
			LE32000    mtx.Gauge `stm:"32000"`
			LE64000    mtx.Gauge `stm:"64000"`
			LE128000   mtx.Gauge `stm:"128000"`
			LE256000   mtx.Gauge `stm:"256000"`
			LE512000   mtx.Gauge `stm:"512000"`
			LE1024000  mtx.Gauge `stm:"1024000"`
			LE2048000  mtx.Gauge `stm:"2048000"`
			LE4096000  mtx.Gauge `stm:"4096000"`
			LE8192000  mtx.Gauge `stm:"8192000"`
			LE16384000 mtx.Gauge `stm:"16384000"`
			Inf        mtx.Gauge `stm:"+Inf"`
		} `stm:"bucket"`
	} `stm:"sync_proxy_rules"`
	RestClientMetrics restClientMetrics `stm:"rest_client"`
}

type restClientMetrics struct {
	HTTPRequestsByStatusCode map[string]mtx.Gauge `stm:"requests"`
	HTTPRequestsByMethod     map[string]mtx.Gauge `stm:"requests"`
}

func (kp *KubeProxy) collect() (map[string]int64, error) {
	raw, err := kp.prom.Scrape()

	if err != nil {
		return nil, err
	}

	kp.mx.SyncProxyRules.Count.Set(
		raw.FindByName("kubeproxy_sync_proxy_rules_latency_microseconds_count").Max())

	kp.collectSyncProxyRuleLatency(raw)
	kp.collectRESTClientHTTPRequests(raw)

	return stm.ToMap(kp.mx), nil
}

func (kp *KubeProxy) collectSyncProxyRuleLatency(raw prometheus.Metrics) {
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

		m := kp.mx.RestClientMetrics.HTTPRequestsByStatusCode

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

	for _, metric := range raw.FindByName(metricName) {
		value := metric.Labels.Get("method")
		if value == "" {
			continue
		}

		m := kp.mx.RestClientMetrics.HTTPRequestsByMethod

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
