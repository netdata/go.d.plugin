package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

type metrics struct {
	PanicCount mtx.Gauge `stm:"panic_count"`
}
