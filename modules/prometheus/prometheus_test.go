package prometheus

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrometheus_Collect(t *testing.T) {
	tests := map[string]struct {
		input         [][]string
		wantCharts    int
		wantCollected map[string]int64
	}{
		"success on Gauge with metadata": {
			input: [][]string{
				{
					`# HELP test_gauge_value Test gauge value.`,
					`# TYPE test_gauge_value gauge`,
					`test_gauge_value{labelA="valueA01",labelB="valueB01"} 1`,
					`test_gauge_value{labelA="valueA02",labelB="valueB02"} 2`,
					`test_gauge_value{labelA="valueA03",labelB="valueB03"} 3`,
					`test_gauge_value{labelA="valueA04",labelB="valueB04"} 4`,
				},
			},
			wantCharts: 4,
			wantCollected: map[string]int64{
				"test_gauge_value-labelA=valueA01-labelB=valueB01": 1000,
				"test_gauge_value-labelA=valueA02-labelB=valueB02": 2000,
				"test_gauge_value-labelA=valueA03-labelB=valueB03": 3000,
				"test_gauge_value-labelA=valueA04-labelB=valueB04": 4000,
			},
		},
		"success on Counter with metadata": {
			input: [][]string{
				{
					`# HELP test_counter_value_total Test gauge value.`,
					`# TYPE test_counter_value_total counter`,
					`test_counter_value_total{labelA="valueA01",labelB="valueB01"} 1`,
					`test_counter_value_total{labelA="valueA02",labelB="valueB02"} 2`,
					`test_counter_value_total{labelA="valueA03",labelB="valueB03"} 3`,
					`test_counter_value_total{labelA="valueA04",labelB="valueB04"} 4`,
				},
			},
			wantCharts: 4,
			wantCollected: map[string]int64{
				"test_counter_value_total-labelA=valueA01-labelB=valueB01": 1000,
				"test_counter_value_total-labelA=valueA02-labelB=valueB02": 2000,
				"test_counter_value_total-labelA=valueA03-labelB=valueB03": 3000,
				"test_counter_value_total-labelA=valueA04-labelB=valueB04": 4000,
			},
		},
		"success on Summary with metadata": {
			input: [][]string{
				{
					`# HELP test_summary_seconds Test summary value.`,
					`# TYPE test_summary_seconds summary`,
					`test_summary_seconds{interval="15s",quantile="0.01"} 14.999892842`,
					`test_summary_seconds{interval="15s",quantile="0.05"} 14.999933467`,
					`test_summary_seconds{interval="15s",quantile="0.5"} 15.000030499`,
					`test_summary_seconds{interval="15s",quantile="0.9"} 15.000099345`,
					`test_summary_seconds{interval="15s",quantile="0.99"} 15.000169848`,
					`test_summary_seconds_sum{interval="15s"} 314505.6192476938`,
					`test_summary_seconds_count{interval="15s"} 20967`,
					`test_summary_seconds{interval="30s",quantile="0.01"} 14.999892842`,
					`test_summary_seconds{interval="30s",quantile="0.05"} 14.999933467`,
					`test_summary_seconds{interval="30s",quantile="0.5"} 15.000030499`,
					`test_summary_seconds{interval="30s",quantile="0.9"} 15.000099345`,
					`test_summary_seconds{interval="30s",quantile="0.99"} 15.000169848`,
					`test_summary_seconds_sum{interval="30s"} 314505.6192476938`,
					`test_summary_seconds_count{interval="30s"} 20967`,
				},
			},
			wantCharts: 2,
			wantCollected: map[string]int64{
				"test_summary_seconds-interval=15s_quantile=0.01": 14999,
				"test_summary_seconds-interval=15s_quantile=0.05": 14999,
				"test_summary_seconds-interval=15s_quantile=0.5":  15000,
				"test_summary_seconds-interval=15s_quantile=0.9":  15000,
				"test_summary_seconds-interval=15s_quantile=0.99": 15000,
				"test_summary_seconds-interval=30s_quantile=0.01": 14999,
				"test_summary_seconds-interval=30s_quantile=0.05": 14999,
				"test_summary_seconds-interval=30s_quantile=0.5":  15000,
				"test_summary_seconds-interval=30s_quantile=0.9":  15000,
				"test_summary_seconds-interval=30s_quantile=0.99": 15000,
			},
		},
		"success on Histogram with metadata": {
			input: [][]string{
				{
					`# HELP test_histogram_seconds Test histogram value.`,
					`# TYPE test_histogram_seconds histogram`,
					`test_histogram_seconds_bucket{verb="GET",le="0.001"} 0`,
					`test_histogram_seconds_bucket{verb="GET",le="0.002"} 0`,
					`test_histogram_seconds_bucket{verb="GET",le="0.004"} 0`,
					`test_histogram_seconds_bucket{verb="GET",le="0.008"} 0`,
					`test_histogram_seconds_bucket{verb="GET",le="0.016"} 0`,
					`test_histogram_seconds_bucket{verb="GET",le="0.032"} 2`,
					`test_histogram_seconds_bucket{verb="GET",le="0.064"} 2`,
					`test_histogram_seconds_bucket{verb="GET",le="0.128"} 2`,
					`test_histogram_seconds_bucket{verb="GET",le="0.256"} 3`,
					`test_histogram_seconds_bucket{verb="GET",le="0.512"} 3`,
					`test_histogram_seconds_bucket{verb="GET",le="+Inf"} 3`,
					`test_histogram_seconds_sum{verb="GET"} 0.28126861`,
					`test_histogram_seconds_count{verb="GET"} 3`,
					`test_histogram_seconds_bucket{verb="POST",le="0.001"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.002"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.004"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.008"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.016"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.032"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.064"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.128"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.256"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="0.512"} 0`,
					`test_histogram_seconds_bucket{verb="POST",le="+Inf"} 1`,
					`test_histogram_seconds_sum{verb="POST"} 4.008446017`,
					`test_histogram_seconds_count{verb="POST"} 1`,
				},
			},
			wantCharts: 2,
			wantCollected: map[string]int64{
				"test_histogram_seconds-verb=GET_bucket=+Inf":   3,
				"test_histogram_seconds-verb=GET_bucket=0.001":  0,
				"test_histogram_seconds-verb=GET_bucket=0.002":  0,
				"test_histogram_seconds-verb=GET_bucket=0.004":  0,
				"test_histogram_seconds-verb=GET_bucket=0.008":  0,
				"test_histogram_seconds-verb=GET_bucket=0.016":  0,
				"test_histogram_seconds-verb=GET_bucket=0.032":  2,
				"test_histogram_seconds-verb=GET_bucket=0.064":  2,
				"test_histogram_seconds-verb=GET_bucket=0.128":  2,
				"test_histogram_seconds-verb=GET_bucket=0.256":  3,
				"test_histogram_seconds-verb=GET_bucket=0.512":  3,
				"test_histogram_seconds-verb=POST_bucket=+Inf":  1,
				"test_histogram_seconds-verb=POST_bucket=0.001": 0,
				"test_histogram_seconds-verb=POST_bucket=0.002": 0,
				"test_histogram_seconds-verb=POST_bucket=0.004": 0,
				"test_histogram_seconds-verb=POST_bucket=0.008": 0,
				"test_histogram_seconds-verb=POST_bucket=0.016": 0,
				"test_histogram_seconds-verb=POST_bucket=0.032": 0,
				"test_histogram_seconds-verb=POST_bucket=0.064": 0,
				"test_histogram_seconds-verb=POST_bucket=0.128": 0,
				"test_histogram_seconds-verb=POST_bucket=0.256": 0,
				"test_histogram_seconds-verb=POST_bucket=0.512": 0,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prom, cleanup := preparePrometheus(t, test.input)
			defer cleanup()

			mx := prom.Collect()
			//fmt.Println(mx)
			//
			//m := mx
			//l := make([]string, 0)
			//for k := range m {
			//	l = append(l, k)
			//}
			//sort.Strings(l)
			//for _, value := range l {
			//	fmt.Println(fmt.Sprintf("\"%s\": %d,", value, m[value]))
			//}
			//return

			require.Equal(t, test.wantCollected, mx)
			if len(test.wantCollected) > 0 {
				assert.Equal(t, test.wantCharts, len(*prom.Charts()))
			}
		})
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

func preparePrometheusEndpoint(metrics [][]string) *httptest.Server {
	var rv []string
	for _, v := range metrics {
		rv = append(rv, strings.Join(v, "\n")+"\n")
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

//func genSimpleMetrics(num int) (metrics []string) {
//	return genMetrics("netdata_generated_metric_count", 0, num)
//}
//
//func genSimpleMetricsFrom(start, num int) (metrics []string) {
//	return genMetrics("netdata_generated_metric_count", start, num)
//}
//
//func genMetrics(name string, start, end int) (metrics []string) {
//	var line string
//	for i := start; i < start+end; i++ {
//		if i%2 == 0 {
//			line = fmt.Sprintf(`%s{number="%d",value="even"} %d`, name, i, i)
//		} else {
//			line = fmt.Sprintf(`%s{number="%d",value="odd"} %d`, name, i, i)
//		}
//		metrics = append(metrics, line)
//	}
//	return metrics
//}
