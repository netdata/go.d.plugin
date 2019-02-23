package nvidia_nvml

//
//import (
//	"sync/atomic"
//
//	"github.com/mindprince/gonvml"
//	"github.com/netdata/go-orchestrator/module"
//)
//
//func init() {
//	creator := module.Creator{
//		//DisabledByDefault: true,
//		Create: func() module.Module { return New() },
//	}
//
//	module.Register("nvidia_nvml", creator)
//}
//
//var once = func() func() bool {
//	var c int64
//	return func() bool {
//		return atomic.AddInt64(&c, 1) == 1
//	}
//}()
//
//// New creates NvidiaNVML with default values.
//func New() *NvidiaNVML {
//	return &NvidiaNVML{
//		charts: &Charts{},
//	}
//}
//
//// NvidiaNVML NvidiaNVML module.
//type NvidiaNVML struct {
//	module.Base
//	charts *Charts
//}
//
//// Cleanup makes cleanup.
//func (NvidiaNVML) Cleanup() {
//	_ = gonvml.Shutdown()
//}
//
//// Init makes initialization.
//func (n NvidiaNVML) Init() bool {
//	if !once() {
//		n.Error("only one job is supported")
//		return false
//	}
//
//	if err := gonvml.Initialize(); err != nil {
//		n.Errorf("error on nvml initialization : %v", err)
//		return false
//	}
//	return true
//}
//
//// Check makes check.
//func (n *NvidiaNVML) Check() bool {
//	gpus, err := getGPUs()
//
//	if err != nil {
//		n.Error(err)
//		return false
//	}
//
//	for _, g := range gpus {
//		_ = n.charts.Add(*createGPUCharts(g)...)
//	}
//
//	return len(charts) > 0
//}
//
//// Charts creates Charts.
//func (n NvidiaNVML) Charts() *Charts {
//	return n.charts
//}
//
//// Collect collects metrics.
//func (n *NvidiaNVML) Collect() map[string]int64 {
//	m := make(map[string]int64)
//
//	gpus, err := getGPUs()
//
//	if err != nil {
//		n.Error(err)
//		return nil
//	}
//
//	for _, gpu := range gpus {
//		for k, v := range gpu.stats.asMap() {
//			m[gpu.uniqName()+"_"+k] = v
//		}
//	}
//
//	return m
//}
