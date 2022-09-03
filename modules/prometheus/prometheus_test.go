// SPDX-License-Identifier: GPL-3.0-or-later

package prometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/netdata/go.d.plugin/pkg/prometheus/selector"
	"github.com/netdata/go.d.plugin/pkg/tlscfg"
	"github.com/netdata/go.d.plugin/pkg/web"

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
			config: Config{HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9090/metric"}}},
		},
		"invalid selector syntax": {
			config: Config{
				HTTP:     web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9090/metric"}},
				Selector: selector.Expr{Allow: []string{`name{label=#"value"}`}},
			},
			wantFail: true,
		},
		"invalid group selector syntax": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9090/metric"}},
				Grouping: []GroupOption{
					{Selector: `name{label=#"value"}`, ByLabel: "label"},
				},
			},
			wantFail: true,
		},
		"empty group selector": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9090/metric"}},
				Grouping: []GroupOption{
					{Selector: "", ByLabel: "label"},
				},
			},
			wantFail: true,
		},
		"empty group 'by_label'": {
			config: Config{
				HTTP: web.HTTP{Request: web.Request{URL: "http://127.0.0.1:9090/metric"}},
				Grouping: []GroupOption{
					{Selector: "name", ByLabel: ""},
				},
			},
			wantFail: true,
		},
		"default": {
			config:   New().Config,
			wantFail: true,
		},
		"nonexistent TLS CA": {
			config: Config{HTTP: web.HTTP{
				Request: web.Request{URL: "http://127.0.0.1:9090/metric"},
				Client:  web.Client{TLSConfig: tlscfg.TLSConfig{TLSCA: "testdata/tls"}}}},
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

func TestPrometheus_Collect_WithExpectedPrefix(t *testing.T) {
	tests := map[string]struct {
		prepare       func(t *testing.T) (prom *Prometheus, cleanup func())
		wantCollected bool
	}{
		"fails on metrics without expected prefix": {
			wantCollected: false,
			prepare: func(t *testing.T) (prom *Prometheus, cleanup func()) {
				prom, cleanup = preparePrometheusValidData(t)
				prom.ExpectedPrefix = "prefix_"
				return prom, cleanup
			},
		},
		"success on metrics with expected prefix": {
			wantCollected: true,
			prepare: func(t *testing.T) (prom *Prometheus, cleanup func()) {
				prom, cleanup = preparePrometheusValidData(t)
				prom.ExpectedPrefix = "some_metric_"
				return prom, cleanup
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := test.prepare(t)
			defer cleanup()

			if test.wantCollected {
				assert.NotNil(t, prom.Collect())
			} else {
				assert.Nil(t, prom.Collect())
			}
		})
	}
}

func TestPrometheus_Collect(t *testing.T) {
	type testGroup map[string]struct {
		input         [][]string
		maxTS         int
		wantCollected map[string]int64
	}

	testGauge := testGroup{
		"with metadata": {
			input: [][]string{
				{
					`# HELP prometheus_sd_discovered_targets Current number of discovered targets.`,
					`# TYPE prometheus_sd_discovered_targets gauge`,
					`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
				},
			},
			wantCollected: map[string]int64{
				"prometheus_sd_discovered_targets|config=config-0,name=notify":               0,
				"prometheus_sd_discovered_targets|config=node_exporter_notebook,name=scrape": 1000,
				"prometheus_sd_discovered_targets|config=node_exporter,name=scrape":          1000,
				"prometheus_sd_discovered_targets|config=prometheus,name=scrape":             1000,
				"series":  4,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		"without metadata": {
			input: [][]string{
				{
					`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
				},
			},
			wantCollected: map[string]int64{
				"prometheus_sd_discovered_targets|config=config-0,name=notify":               0,
				"prometheus_sd_discovered_targets|config=node_exporter_notebook,name=scrape": 1000,
				"prometheus_sd_discovered_targets|config=node_exporter,name=scrape":          1000,
				"prometheus_sd_discovered_targets|config=prometheus,name=scrape":             1000,
				"series":  4,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
	}

	testCounter := testGroup{
		"with metadata": {
			input: [][]string{
				{
					`# HELP prometheus_sd_kubernetes_events_total The number of Kubernetes events handled.`,
					`# TYPE prometheus_sd_kubernetes_events_total counter`,
					`prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 1`,
					`prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 2`,
					`prometheus_sd_kubernetes_events_total{event="add",role="node"} 3`,
					`prometheus_sd_kubernetes_events_total{event="add",role="pod"} 4`,
					`prometheus_sd_kubernetes_events_total{event="add",role="service"} 5`,
				},
			},
			wantCollected: map[string]int64{
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
		"without metadata": {
			input: [][]string{
				{
					`prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 1`,
					`prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 2`,
					`prometheus_sd_kubernetes_events_total{event="add",role="node"} 3`,
					`prometheus_sd_kubernetes_events_total{event="add",role="pod"} 4`,
					`prometheus_sd_kubernetes_events_total{event="add",role="service"} 5`,
				},
			},
			wantCollected: map[string]int64{
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
	}

	testSummary := testGroup{
		"with metadata": {
			input: [][]string{
				{
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
			},
			wantCollected: map[string]int64{
				"prometheus_target_interval_length_seconds_count|interval=15s":         20967000,
				"prometheus_target_interval_length_seconds_sum|interval=15s":           314505619,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.01": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.05": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.5":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.9":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.99": 15000,
				"series":  7,
				"metrics": 3,
				"charts":  int64(3 + len(statsCharts)),
			},
		},
		"without metadata": {
			input: [][]string{
				{
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
					`prometheus_target_interval_length_seconds_sum{interval="15s"} 314505.6192476938`,
					`prometheus_target_interval_length_seconds_count{interval="15s"} 20967`,
				},
			},
			wantCollected: map[string]int64{
				"prometheus_target_interval_length_seconds_count|interval=15s":         20967000,
				"prometheus_target_interval_length_seconds_sum|interval=15s":           314505619,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.01": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.05": 14999,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.5":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.9":  15000,
				"prometheus_target_interval_length_seconds|interval=15s,quantile=0.99": 15000,
				"series":  7,
				"metrics": 3,
				"charts":  int64(3 + len(statsCharts)),
			},
		},
	}

	testHistogram := testGroup{
		"with metadata": {
			input: [][]string{
				{
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
			},
			wantCollected: map[string]int64{
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
				"charts":  int64(3 + len(statsCharts)),
			},
		},
		"without metadata": {
			input: [][]string{
				{
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
			},
			wantCollected: map[string]int64{
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
				"charts":  int64(3 + len(statsCharts)),
			},
		},
	}

	testSanitize := testGroup{
		`label value contains '\x' encoding`: {
			input: [][]string{
				{
					`node_restart_total{state="systemd-fsck@dev-disk-by\\x2duuid\\x2dDD70.service"} 0`,
				},
			},
			wantCollected: map[string]int64{
				"node_restart_total|state=systemd-fsck@dev-disk-by-uuid-DD70.service": 0,

				"series":  1,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		`label value contains '\'`: {
			input: [][]string{
				{
					`http_time_to_write_seconds{code="200",route="^/([^/]+/){1,}[^/]+/uploads\z"} 0`,
				},
			},
			wantCollected: map[string]int64{
				"http_time_to_write_seconds|code=200,route=^/([^/]+/){1,}[^/]+/uploads_z": 0,

				"series":  1,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
		`label value contains ' '`: {
			input: [][]string{
				{
					`jvm_memory_max_bytes{area="heap",id="Eden Space"} 0`,
				},
			},
			wantCollected: map[string]int64{
				"jvm_memory_max_bytes|area=heap,id=Eden_Space": 0,

				"series":  1,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
		},
	}

	maxTSPerMetric := New().MaxTSPerMetric
	maxTS := 4
	testTSLimits := testGroup{
		"len(ts) > per metric limit": {
			input: [][]string{
				genSimpleMetrics(maxTSPerMetric + 1),
			},
			wantCollected: map[string]int64{
				"series":  int64(maxTSPerMetric + 1),
				"metrics": 0,
				"charts":  int64(len(statsCharts)),
			},
		},
		"len(ts) > global limit": {
			wantCollected: map[string]int64{
				"prometheus_sd_discovered_targets|config=config-0,name=notify":               0,
				"prometheus_sd_discovered_targets|config=node_exporter_notebook,name=scrape": 1000,
				"prometheus_sd_discovered_targets|config=node_exporter,name=scrape":          1000,
				"prometheus_sd_discovered_targets|config=prometheus,name=scrape":             1000,
				"series":  4,
				"metrics": 1,
				"charts":  int64(1 + len(statsCharts)),
			},
			maxTS: maxTS,
			input: [][]string{
				{
					`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets2{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="prometheus",name="scrape"} 1`,
				},
				{
					`prometheus_sd_discovered_targets{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets{config="prometheus",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="config-0",name="notify"} 0`,
					`prometheus_sd_discovered_targets2{config="node_exporter",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="node_exporter_notebook",name="scrape"} 1`,
					`prometheus_sd_discovered_targets2{config="prometheus",name="scrape"} 1`,
				},
			},
		},
	}

	tests := map[string]testGroup{
		"Gauge":     testGauge,
		"Counter":   testCounter,
		"Summary":   testSummary,
		"Histogram": testHistogram,
		"Sanitize":  testSanitize,
		"TSLimits":  testTSLimits,
	}

	for groupName, group := range tests {
		for name, test := range group {
			name = fmt.Sprintf("%s: %s", groupName, name)

			t.Run(name, func(t *testing.T) {
				prom, cleanup := preparePrometheus(t, test.input)
				defer cleanup()
				if test.maxTS != 0 {
					prom.MaxTS = test.maxTS
				}

				var collected map[string]int64

				for i := 0; i < len(test.input); i++ {
					collected = prom.Collect()
				}

				assert.Equal(t, test.wantCollected, collected)
				if test.wantCollected != nil {
					ensureCollectedHasAllChartsDimsVarsIDs(t, prom, collected)
				}
			})
		}
	}
}

func TestPrometheus_ForceAbsoluteAlgorithm(t *testing.T) {
	input := [][]string{
		{
			`# HELP prometheus_sd_kubernetes_events_total The number of Kubernetes events handled.`,
			`# TYPE prometheus_sd_kubernetes_events_total counter`,
			`prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 1`,
			`prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 2`,
			`prometheus_sd_kubernetes_events_total{event="add",role="node"} 3`,
			`prometheus_sd_kubernetes_events_total{event="add",role="pod"} 4`,
			`prometheus_sd_kubernetes_events_total{event="add",role="service"} 5`,
		},
	}

	prom, cleanup := preparePrometheus(t, input)
	defer cleanup()
	prom.forceAbsoluteAlgorithm = matcher.TRUE()

	assert.NotEmpty(t, prom.Collect())

	for _, c := range *prom.Charts() {
		if c.ID != "prometheus_sd_kubernetes_events_total" {
			continue
		}
		for _, d := range c.Dims {
			assert.Equal(t, module.Absolute, d.Algo)
		}
	}
}

func TestPrometheus_Collect_WithSelector(t *testing.T) {
	tests := map[string]struct {
		input             [][]string
		sr                selector.Expr
		expectedCollected map[string]int64
	}{
		"simple filtering": {
			input: [][]string{
				{
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
			},
			sr: selector.Expr{
				Allow: []string{
					"prometheus_*_sum prometheus_*_count",
				},
			},
			expectedCollected: map[string]int64{
				"prometheus_tsdb_compaction_chunk_range_seconds_count": 84164000,
				"prometheus_tsdb_compaction_chunk_range_seconds_sum":   150091952011000,
				"series":  2,
				"metrics": 2,
				"charts":  int64(2 + len(statsCharts)),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := preparePrometheusWithSelector(t, test.input, test.sr)
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

func TestPrometheus_Collect_DefaultGrouping(t *testing.T) {
	type testGroup map[string]struct {
		input            [][]string
		wantCharts       int
		wantActiveCharts int
	}

	testGaugeAndCounter := testGroup{
		"scrapes| 1st: <= desired": {
			input: [][]string{
				genSimpleMetrics(desiredDim),
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
		"scrapes| 1st: > desired (split into 2)": {
			input: [][]string{
				genSimpleMetrics(desiredDim + 1),
			},
			wantCharts:       2,
			wantActiveCharts: 2,
		},
		"scrapes| 1st: > desired (split into 4)": {
			input: [][]string{
				genSimpleMetrics(desiredDim*3 + 1),
			},
			wantCharts:       4,
			wantActiveCharts: 4,
		},
		"scrapes| 1st: <= desired, 2nd: == max": {
			input: [][]string{
				genSimpleMetrics(desiredDim),
				genSimpleMetrics(maxDim),
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
		"scrapes| 1st: <= desired, 2nd: > max": {
			input: [][]string{
				genSimpleMetrics(desiredDim),
				genSimpleMetrics(maxDim + 1),
			},
			wantCharts:       3,
			wantActiveCharts: 2,
		},
		"scrapes| 1st: > desired, 2nd: > max": {
			input: [][]string{
				genSimpleMetrics(desiredDim + 1),
				genSimpleMetrics(maxDim*2 + 1),
			},
			wantCharts:       5,
			wantActiveCharts: 3,
		},
		"scrapes| 1st: <= desired, 2nd: == max, 3rd: > max": {
			input: [][]string{
				genSimpleMetrics(desiredDim),
				genSimpleMetrics(maxDim),
				genSimpleMetrics(maxDim + 1),
			},
			wantCharts:       3,
			wantActiveCharts: 2,
		},
		"scrapes| 1st: = desired, 2nd: = desired but different set": {
			input: [][]string{
				genSimpleMetrics(desiredDim),
				genSimpleMetricsFrom(desiredDim+1, desiredDim),
			},
			wantCharts:       2,
			wantActiveCharts: 1,
		},
	}

	testSummary := testGroup{
		"several time series in one scrape": {
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
			wantCharts:       4,
			wantActiveCharts: 4,
		},
		"several time series in several scrapes": {
			input: [][]string{
				{
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
					`prometheus_target_interval_length_seconds_sum{interval="15s"} 314505.6192476938`,
					`prometheus_target_interval_length_seconds_count{interval="15s"} 20967`,
				},
				{
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="30s",quantile="0.99"} 15.000169848`,
					`prometheus_target_interval_length_seconds_sum{interval="30s"} 314505.6192476938`,
					`prometheus_target_interval_length_seconds_count{interval="30s"} 20967`,
				},
			},
			wantCharts:       4,
			wantActiveCharts: 4,
		},
	}

	testHistogram := testGroup{
		"several time series in one scrape": {
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
			wantCharts:       4,
			wantActiveCharts: 4,
		},
		"several time series in several scrapes": {
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
				},
				{
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
			wantCharts:       4,
			wantActiveCharts: 4,
		},
	}

	tests := map[string]testGroup{
		"Gauge,Counter": testGaugeAndCounter,
		"Summary":       testSummary,
		"Histogram":     testHistogram,
	}

	for groupName, group := range tests {
		for name, test := range group {
			name = fmt.Sprintf("%s: %s", groupName, name)

			t.Run(name, func(t *testing.T) {
				prom, cleanup := preparePrometheus(t, test.input)
				defer cleanup()

				for i := 0; i < len(test.input); i++ {
					prom.Collect()
				}

				var active int
				for _, chart := range *prom.Charts() {
					if !chart.Obsolete {
						active++
					}
				}
				test.wantCharts += len(statsCharts)
				test.wantActiveCharts += len(statsCharts)

				assert.Equalf(t, test.wantCharts, len(*prom.Charts()), "expected charts")
				assert.Equalf(t, test.wantActiveCharts, active, "expected active charts")
			})
		}
	}
}

func TestPrometheus_Collect_UserDefinedGrouping(t *testing.T) {
	type testGroup map[string]struct {
		input            [][]string
		grouping         []GroupOption
		wantCharts       int
		wantActiveCharts int
	}

	testGaugeAndCounter := testGroup{
		"not matches, one grouping": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated_not_exists", ByLabel: "value"},
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
		"not matches, several groupings": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated_not_exists1", ByLabel: "value"},
				{Selector: "generated_not_exists2", ByLabel: "value"},
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
		"matches but no label": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value_not_exists"},
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
		"matches, one grouping": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value"},
			},
			wantCharts:       3,
			wantActiveCharts: 3,
		},
		"partial matches, one grouping": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: `generated{value="odd"}`, ByLabel: "value"},
			},
			wantCharts:       3,
			wantActiveCharts: 3,
		},
		"matches, several groupings": {
			input: [][]string{
				genMetrics("generated1", 0, desiredDim),
				genMetrics("generated2", 0, desiredDim),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated1", ByLabel: "value"},
				{Selector: "generated2", ByLabel: "value"},
			},
			wantCharts:       5,
			wantActiveCharts: 5,
		},
		"partial matches, several groupings": {
			input: [][]string{
				genMetrics("generated1", 0, desiredDim),
				genMetrics("generated2", 0, desiredDim),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: `generated1{value="odd"}`, ByLabel: "value"},
				{Selector: `generated2{value="odd"}`, ByLabel: "value"},
			},
			wantCharts:       5,
			wantActiveCharts: 5,
		},
		"matches, one grouping, one chart > desired": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim*2+1),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value"},
			},
			wantCharts:       4,
			wantActiveCharts: 4,
		},
		"matches, one grouping, two charts > desired": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim*2+2),
				genMetrics("generated_not_matches", 0, desiredDim),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value"},
			},
			wantCharts:       5,
			wantActiveCharts: 5,
		},
		"matches, one grouping, 1st scrape < desired, 2nd = max": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim*2),
				genMetrics("generated", 0, maxDim*2),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value"},
			},
			wantCharts:       2,
			wantActiveCharts: 2,
		},
		"matches, one grouping, 1st scrape < desired, 2nd > max": {
			input: [][]string{
				genMetrics("generated", 0, desiredDim*2),
				genMetrics("generated", 0, maxDim*2+1),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "value"},
			},
			wantCharts:       6,
			wantActiveCharts: 4,
		},
		"group by 2 labels": {
			input: [][]string{
				genMetrics("generated", 0, 4),
			},
			grouping: []GroupOption{
				{Selector: "generated", ByLabel: "number value"},
			},
			wantCharts:       4,
			wantActiveCharts: 4,
		},
	}

	testSummary := testGroup{
		"grouping matches, but doesnt apply": {
			input: [][]string{
				{
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.01"} 14.999892842`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.05"} 14.999933467`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.5"} 15.000030499`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.9"} 15.000099345`,
					`prometheus_target_interval_length_seconds{interval="15s",quantile="0.99"} 15.000169848`,
				},
			},
			grouping: []GroupOption{
				{
					Selector: "prometheus_target_interval_length_seconds",
					ByLabel:  "quantile",
				},
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
	}

	testHistogram := testGroup{
		"grouping matches, but doesnt apply": {
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
				},
			},
			grouping: []GroupOption{
				{
					Selector: "prometheus_http_request_duration_seconds_bucket",
					ByLabel:  "le",
				},
			},
			wantCharts:       1,
			wantActiveCharts: 1,
		},
	}

	tests := map[string]testGroup{
		"Gauge,Counter": testGaugeAndCounter,
		"Summary":       testSummary,
		"Histogram":     testHistogram,
	}

	for groupName, group := range tests {
		for name, test := range group {
			name = fmt.Sprintf("%s: %s", groupName, name)

			t.Run(name, func(t *testing.T) {
				prom, cleanup := preparePrometheusWithGrouping(t, test.input, test.grouping)
				defer cleanup()

				for i := 0; i < len(test.input); i++ {
					prom.Collect()
				}

				var active int
				for _, chart := range *prom.Charts() {
					if !chart.Obsolete {
						active++
					}
				}
				test.wantCharts += len(statsCharts)
				test.wantActiveCharts += len(statsCharts)

				assert.Equalf(t, test.wantCharts, len(*prom.Charts()), "expected charts")
				assert.Equalf(t, test.wantActiveCharts, active, "expected active charts")
			})
		}
	}
}

func genSimpleMetrics(num int) (metrics []string) {
	return genMetrics("netdata_generated_metric_count", 0, num)
}

func genSimpleMetricsFrom(start, num int) (metrics []string) {
	return genMetrics("netdata_generated_metric_count", start, num)
}

func genMetrics(name string, start, end int) (metrics []string) {
	var line string
	for i := start; i < start+end; i++ {
		if i%2 == 0 {
			line = fmt.Sprintf(`%s{number="%d",value="even"} %d`, name, i, i)
		} else {
			line = fmt.Sprintf(`%s{number="%d",value="odd"} %d`, name, i, i)
		}
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

func preparePrometheus(t *testing.T, metrics [][]string) (*Prometheus, func()) {
	t.Helper()
	require.NotZero(t, metrics)

	srv := preparePrometheusEndpoint(metrics)
	prom := New()
	prom.URL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusWithSelector(t *testing.T, metrics [][]string, sr selector.Expr) (*Prometheus, func()) {
	t.Helper()
	require.NotZero(t, metrics)
	srv := preparePrometheusEndpoint(metrics)

	prom := New()
	prom.URL = srv.URL
	prom.Selector = sr
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusWithGrouping(t *testing.T, metrics [][]string, grp []GroupOption) (*Prometheus, func()) {
	t.Helper()
	require.NotZero(t, metrics)
	srv := preparePrometheusEndpoint(metrics)

	prom := New()
	prom.URL = srv.URL
	prom.Grouping = grp
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
	prom.URL = srv.URL
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
	prom.URL = srv.URL
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
	prom.URL = srv.URL
	require.True(t, prom.Init())

	return prom, srv.Close
}

func preparePrometheusConnectionRefused(t *testing.T) (*Prometheus, func()) {
	t.Helper()
	prom := New()
	prom.URL = "http://127.0.0.1:38001/metrics"
	require.True(t, prom.Init())

	return prom, func() {}
}

func preparePrometheusEndpoint(metrics [][]string) *httptest.Server {
	var rv []string
	for _, v := range metrics {
		rv = append(rv, strings.Join(v, "\n"))
	}
	var i int
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if i <= len(metrics)-1 {
				_, _ = w.Write([]byte(rv[i]))
			} else {
				_, _ = w.Write([]byte(rv[len(rv)-1]))
			}
			i++
		}))
}
