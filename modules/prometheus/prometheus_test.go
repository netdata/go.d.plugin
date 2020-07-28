package prometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.Implements(t, (*module.Module)(nil), New())
}

func TestPrometheus_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"non empty URL": {
			config: Config{HTTP: web.HTTP{Request: web.Request{UserURL: "http://127.0.0.1:9090/metric"}}},
		},
		//"default": {
		//	config:   New().Config,
		//	wantFail: true,
		//},
		"nonexistent TLS CA": {
			config: Config{HTTP: web.HTTP{
				Request: web.Request{UserURL: "http://127.0.0.1:9090/metric"},
				Client:  web.Client{ClientTLSConfig: web.ClientTLSConfig{TLSCA: "testdata/tls"}}}},
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom := New()
			prom.Config = test.config

			if test.wantFail {
				assert.False(t, prom.Init())
			} else {
				assert.True(t, prom.Init())
			}
		})
	}
}

func TestPrometheus_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func(*testing.T) (prom *Prometheus, cleanup func())
		wantFail bool
	}{
		"valid data":         {prepare: preparePrometheusValidData},
		"invalid data":       {prepare: preparePrometheusInvalidData, wantFail: true},
		"404":                {prepare: preparePrometheus404, wantFail: true},
		"connection refused": {prepare: preparePrometheusConnectionRefused, wantFail: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantFail {
				assert.False(t, prom.Check())
			} else {
				assert.True(t, prom.Check())
			}
		})
	}

}

func TestPrometheus_Charts(t *testing.T) {
	assert.NotNil(t, New().Charts())
}

func TestPrometheus_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestPrometheus_Collect_ReturnsNilOnError(t *testing.T) {
	tests := map[string]func(*testing.T) (prom *Prometheus, cleanup func()){
		"invalid data":       preparePrometheusInvalidData,
		"404":                preparePrometheus404,
		"connection refused": preparePrometheusConnectionRefused,
	}

	for name, prepare := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := prepare(t)
			defer cleanup()
			assert.Nil(t, prom.Collect())
		})
	}
}

func TestPrometheus_Collect(t *testing.T) {
	tests := map[string]struct {
		input             []string
		expectedCollected map[string]int64
	}{
		"GAUGE metrics": {
			input: []string{
				`# HELP prometheus_sd_discovered_targets Current number of discovered targets.`,
				`# TYPE prometheus_sd_discovered_targets gauge`,
				`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
				`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
				`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
				`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
			},
			expectedCollected: map[string]int64{
				"prometheus_sd_discovered_targets|config=config-0,name=notify":               0,
				"prometheus_sd_discovered_targets|config=node_exporter_notebook,name=scrape": 1000,
				"prometheus_sd_discovered_targets|config=node_exporter,name=scrape":          1000,
				"prometheus_sd_discovered_targets|config=prometheus,name=scrape":             1000,
				"series":  4,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		"GAUGE no meta metrics": {
			input: []string{
				`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
				`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
				`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
				`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
			},
			expectedCollected: map[string]int64{
				"prometheus_sd_discovered_targets|config=config-0,name=notify":               0,
				"prometheus_sd_discovered_targets|config=node_exporter_notebook,name=scrape": 1000,
				"prometheus_sd_discovered_targets|config=node_exporter,name=scrape":          1000,
				"prometheus_sd_discovered_targets|config=prometheus,name=scrape":             1000,
				"series":  4,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		"COUNTER metrics": {
			input: []string{
				`# HELP prometheus_sd_kubernetes_events_total The number of Kubernetes events handled.`,
				`# TYPE prometheus_sd_kubernetes_events_total counter`,
				`prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 1`,
				`prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 2`,
				`prometheus_sd_kubernetes_events_total{event="add",role="node"} 3`,
				`prometheus_sd_kubernetes_events_total{event="add",role="pod"} 4`,
				`prometheus_sd_kubernetes_events_total{event="add",role="service"} 5`,
			},
			expectedCollected: map[string]int64{
				"prometheus_sd_kubernetes_events_total|event=add,role=endpoints": 1000,
				"prometheus_sd_kubernetes_events_total|event=add,role=ingress":   2000,
				"prometheus_sd_kubernetes_events_total|event=add,role=node":      3000,
				"prometheus_sd_kubernetes_events_total|event=add,role=pod":       4000,
				"prometheus_sd_kubernetes_events_total|event=add,role=service":   5000,
				"series":  5,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		"COUNTER no meta metrics": {
			input: []string{
				`prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 1`,
				`prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 2`,
				`prometheus_sd_kubernetes_events_total{event="add",role="node"} 3`,
				`prometheus_sd_kubernetes_events_total{event="add",role="pod"} 4`,
				`prometheus_sd_kubernetes_events_total{event="add",role="service"} 5`,
			},
			expectedCollected: map[string]int64{
				"prometheus_sd_kubernetes_events_total|event=add,role=endpoints": 1000,
				"prometheus_sd_kubernetes_events_total|event=add,role=ingress":   2000,
				"prometheus_sd_kubernetes_events_total|event=add,role=node":      3000,
				"prometheus_sd_kubernetes_events_total|event=add,role=pod":       4000,
				"prometheus_sd_kubernetes_events_total|event=add,role=service":   5000,
				"series":  5,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		"SUMMARY metrics": {
			input: []string{
				`# HELP prometheus_target_interval_length_seconds Actual intervals between scrapes.`,
				`# TYPE prometheus_target_interval_length_seconds summary`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
				`prometheus_target_interval_length_seconds_sum{interval="15s"} 314505.6192476938`,
				`prometheus_target_interval_length_seconds_count{interval="15s"} 20967`,
			},
			expectedCollected: map[string]int64{
				"prometheus_target_interval_length_seconds_count|interval=15s":         20967000,
				"prometheus_target_interval_length_seconds_sum|interval=15s":           314505619,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.01": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.05": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.5":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.9":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.99": 15000,
				"series":  7,
				"metrics": 3,
				"charts":  int64(4 + len(statsCharts)),
			},
		},
		"SUMMARY no meta metrics": {
			input: []string{
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
				`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
				`prometheus_target_interval_length_seconds_sum{interval="15s"} 314505.6192476938`,
				`prometheus_target_interval_length_seconds_count{interval="15s"} 20967`,
			},
			expectedCollected: map[string]int64{
				"prometheus_target_interval_length_seconds_count|interval=15s":         20967000,
				"prometheus_target_interval_length_seconds_sum|interval=15s":           314505619,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.01": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.05": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.5":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.9":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.99": 15000,
				"series":  7,
				"metrics": 3,
				"charts":  int64(4 + len(statsCharts)),
			},
		},
		"HISTOGRAM metrics": {
			input: []string{
				`# HELP prometheus_tsdb_compaction_chunk_range_seconds Final time range`,
				`# TYPE prometheus_tsdb_compaction_chunk_range_seconds histogram`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="100"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="400"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1600"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6400"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="25600"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="102400"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="409600"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1.6384e+06"} 2000`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6.5536e+06"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="2.62144e+07"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="+Inf"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_sum 1.50091952011e+11`,
				`prometheus_tsdb_compaction_chunk_range_seconds_count 84164`,
			},
			expectedCollected: map[string]int64{
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=+Inf":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=1.6384e+06":  1999000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=100":         0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=102400":      0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=1600":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=2.62144e+07": 0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=25600":       1000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=400":         0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=409600":      0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=6.5536e+06":  82164000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=6400":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_count":                 84164000,
				"prometheus_tsdb_compaction_chunk_range_seconds_sum":                   150091952011000,
				"series":  13,
				"metrics": 3,
				"charts":  int64(4 + len(statsCharts)),
			},
		},
		"HISTOGRAM no meta metrics": {
			input: []string{
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="100"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="400"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1600"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6400"} 0`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="25600"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="102400"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="409600"} 1`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1.6384e+06"} 2000`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6.5536e+06"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="2.62144e+07"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="+Inf"} 84164`,
				`prometheus_tsdb_compaction_chunk_range_seconds_sum 1.50091952011e+11`,
				`prometheus_tsdb_compaction_chunk_range_seconds_count 84164`,
			},
			expectedCollected: map[string]int64{
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=+Inf":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=1.6384e+06":  1999000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=100":         0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=102400":      0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=1600":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=2.62144e+07": 0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=25600":       1000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=400":         0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=409600":      0,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=6.5536e+06":  82164000,
				"prometheus_tsdb_compaction_chunk_range_seconds_bucket|le=6400":        0,
				"prometheus_tsdb_compaction_chunk_range_seconds_count":                 84164000,
				"prometheus_tsdb_compaction_chunk_range_seconds_sum":                   150091952011000,
				"series":  13,
				"metrics": 3,
				"charts":  int64(4 + len(statsCharts)),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			input := strings.Join(test.input, "\n")
			prom, cleanup := preparePrometheus(t, input)
			defer cleanup()

			var collected map[string]int64

			for i := 0; i < 10; i++ {
				collected = prom.Collect()
			}

			assert.Equal(t, test.expectedCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, prom, collected)
		})
	}
}

func TestPrometheus_Collect_Split(t *testing.T) {
	tests := map[string]struct {
		input                   [][]string
		expectedNumCharts       int
		expectedNumActiveCharts int
	}{
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: metrics <= desired": {
			input: [][]string{
				genMetrics(desiredDim),
			},
			expectedNumCharts:       1,
			expectedNumActiveCharts: 1,
		},
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: > desired": {
			input: [][]string{
				genMetrics(desiredDim + 1),
			},
			expectedNumCharts:       2,
			expectedNumActiveCharts: 2,
		},
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: <= desired, 2nd: == max": {
			input: [][]string{
				genMetrics(desiredDim),
				genMetrics(maxDim),
			},
			expectedNumCharts:       1,
			expectedNumActiveCharts: 1,
		},
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: <= desired, 2nd: > max": {
			input: [][]string{
				genMetrics(desiredDim),
				genMetrics(maxDim + 1),
			},
			expectedNumCharts:       4,
			expectedNumActiveCharts: 3,
		},
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: > desired, 2nd: > max": {
			input: [][]string{
				genMetrics(desiredDim + 1),
				genMetrics(maxDim*2 + 1),
			},
			expectedNumCharts:       7,
			expectedNumActiveCharts: 5,
		},
		"GAUGE|COUNTER|UNKNOWN, scrapes: 1st: <= desired, 2nd: == max, 3rd: > max": {
			input: [][]string{
				genMetrics(desiredDim),
				genMetrics(maxDim),
				genMetrics(maxDim + 1),
			},
			expectedNumCharts:       4,
			expectedNumActiveCharts: 3,
		},
		"SUMMARY, several time series": {
			input: [][]string{
				{
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
					`prometheus_target_interval_length_seconds_sum{interval="15s"} 314505.6192476938`,
					`prometheus_target_interval_length_seconds_count{interval="15s"} 20967`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.99"} 15.000169848`,
					`prometheus_target_interval_length_seconds_sum{interval="30s"} 314505.6192476938`,
					`prometheus_target_interval_length_seconds_count{interval="30s"} 20967`,
				},
			},
			expectedNumCharts:       6,
			expectedNumActiveCharts: 6,
		},
		"HISTOGRAM, several time series": {
			input: [][]string{
				{
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="0.1"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="0.2"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="0.4"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="1"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="3"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="8"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="20"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="60"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="120"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/",le="+Inf"} 1`,
					`prometheus_http_request_duration_seconds_sum{handler="/"} 5.9042e-05`,
					`prometheus_http_request_duration_seconds_count{handler="/"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.1"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.2"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.4"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="1"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="3"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="8"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="20"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="60"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="120"} 1`,
					`prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="+Inf"} 1`,
					`prometheus_http_request_duration_seconds_sum{handler="/metrics"} 5.9042e-05`,
					`prometheus_http_request_duration_seconds_count{handler="/metrics"} 1`,
				},
			},
			expectedNumCharts:       6,
			expectedNumActiveCharts: 6,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var input []string
			for _, v := range test.input {
				input = append(input, strings.Join(v, "\n"))
			}
			prom, cleanup := preparePrometheusDynamic(t, input)
			defer cleanup()

			for i := 0; i < len(input); i++ {
				prom.Collect()
			}

			var active int
			for _, chart := range *prom.Charts() {
				if !chart.Obsolete {
					active++
				}
			}

			assert.Equalf(t, test.expectedNumCharts+len(statsCharts), len(*prom.Charts()), "expected charts")
			assert.Equalf(t, test.expectedNumActiveCharts+len(statsCharts), active, "expected active charts")
		})
	}
}

func genMetrics(num int) (metrics []string) {
	for i := 0; i < num; i++ {
		line := fmt.Sprintf(`netdata_generated_metric_count{number="%d"} %d`, i, i)
		metrics = append(metrics, line)
	}
	return metrics
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, prom *Prometheus, collected map[string]int64) {
	for _, chart := range *prom.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok, "collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok, "collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func preparePrometheus(t *testing.T, metrics string) (*Prometheus, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(metrics))
		}))

	prom := New()
	prom.UserURL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusDynamic(t *testing.T, metrics []string) (*Prometheus, func()) {
	t.Helper()
	require.NotZero(t, metrics)
	var i int
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if i <= len(metrics)-1 {
				_, _ = w.Write([]byte(metrics[i]))
			} else {
				_, _ = w.Write([]byte(metrics[len(metrics)-1]))
			}
			i++
		}))

	prom := New()
	prom.UserURL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusValidData(t *testing.T) (*Prometheus, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`some_metric_units 1`))
		}))

	prom := New()
	prom.UserURL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusInvalidData(t *testing.T) (*Prometheus, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello and\n goodbye"))
		}))

	prom := New()
	prom.UserURL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheus404(t *testing.T) (*Prometheus, func()) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
	prom := New()
	prom.UserURL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusConnectionRefused(t *testing.T) (*Prometheus, func()) {
	t.Helper()
	prom := New()
	prom.UserURL = "http://127.0.0.1:38001/metrics"
	require.True(t, prom.Init())

	return prom, func() {}
}
