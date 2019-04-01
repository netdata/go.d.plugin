package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

func newMetrics() *metrics {
	mx := metrics{}
	mx.Request.Count.ByTypeTotal = make(map[string]mtx.Gauge)
	mx.Response.Count.ByRcodeTotal = make(map[string]mtx.Gauge)

	return &mx
}

type metrics struct {
	Panic struct {
		Count struct {
			Total mtx.Gauge `stm:"total"`
		} `stm:"count"`
	} `stm:"panic"`
	Request struct {
		Count struct {
			Total       mtx.Gauge            `stm:"total"`
			ByTypeTotal map[string]mtx.Gauge `stm:"by_type_total"`
		} `stm:"count"`
	} `stm:"request"`
	Response struct {
		Count struct {
			ByRcodeTotal map[string]mtx.Gauge `stm:"by_rcode_total"`
		} `stm:"count"`
	} `stm:"response"`
}
