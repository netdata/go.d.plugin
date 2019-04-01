package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

type metrics struct {
	PanicCountTotal    mtx.Gauge `stm:"panic_count_total"`
	RequestsCountTotal mtx.Gauge `stm:"request_count_total"`
}
