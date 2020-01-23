package cockroachdb

import "github.com/netdata/go.d.plugin/pkg/prometheus"

func collectRuntime(pms prometheus.Metrics) runtimeMetrics {
	var rm runtimeMetrics
	rm.LiveNodes = pms.FindByName("liveness_livenodes").Max()
	rm.SysUptime = pms.FindByName("sys_uptime").Max()
	collectRuntimeMemory(&rm, pms)
	return rm
}

func collectRuntimeMemory(rm *runtimeMetrics, pms prometheus.Metrics) {
	rm.Memory.RSS = pms.FindByName("sys_rss").Max()

	rm.Memory.GoAllocBytes = pms.FindByName("sys_go_allocbytes").Max()
	rm.Memory.GoTotalBytes = pms.FindByName("sys_go_totalbytes").Max()
	rm.Memory.CGoAllocBytes = pms.FindByName("sys_cgo_allocbytes").Max()
	rm.Memory.CGoTotalBytes = pms.FindByName("sys_cgo_totalbytes").Max()
}
