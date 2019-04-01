package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"

	"github.com/netdata/go-orchestrator/module"
)

func (cd *CoreDNS) collect() (map[string]int64, error) {
	raw, err := cd.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	cd.collectPanicTotal(raw, mx)
	cd.collectRequestTotal(raw, mx)
	cd.collectRequestsByTypeTotal(raw, mx)
	cd.collectResponsesByRcodeTotal(raw, mx)

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanicTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_panic_count_total"

	mx.Panic.Count.Total.Set(raw.FindByName(metricName).Max())
}

func (cd CoreDNS) collectRequestTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_request_count_total"

	for _, metric := range raw.FindByName(metricName) {
		mx.Request.Count.Total.Add(metric.Value)
	}
}

func (cd *CoreDNS) collectRequestsByTypeTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_request_type_count_total"
	chartName := "request_type_count_total"
	chart := cd.charts.Get(chartName)

	for _, metric := range raw.FindByName(metricName) {
		typ := metric.Labels.Get("type")
		if typ == "" {
			continue
		}
		if chart == nil {
			_ = cd.charts.Add(chartReqByTypeTotal.Copy())
			chart = cd.charts.Get(chartName)
		}
		dimID := "request_count_by_type_total_" + typ
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: typ, Algo: module.Incremental})
		}

		current := mx.Request.Count.ByTypeTotal[typ].Value()
		mx.Request.Count.ByTypeTotal[typ] = mtx.Gauge(metric.Value + current)
	}
}

func (cd *CoreDNS) collectResponsesByRcodeTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_response_rcode_count_total"
	chartName := "response_rcode_count_total"
	chart := cd.charts.Get(chartName)

	for _, metric := range raw.FindByName(metricName) {
		rcode := metric.Labels.Get("rcode")
		if rcode == "" {
			continue
		}
		if chart == nil {
			_ = cd.charts.Add(chartRespByRcodeTotal.Copy())
			chart = cd.charts.Get(chartName)
		}
		dimID := "response_count_by_rcode_total_" + rcode
		if !chart.HasDim(dimID) {
			_ = chart.AddDim(&Dim{ID: dimID, Name: rcode, Algo: module.Incremental})
		}

		current := mx.Response.Count.ByRcodeTotal[rcode].Value()
		mx.Response.Count.ByRcodeTotal[rcode] = mtx.Gauge(metric.Value + current)
	}
}
