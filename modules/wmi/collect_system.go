package wmi

import (
	"fmt"
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricSysContextSwitchesTotal     = "wmi_system_context_switches_total"
	metricSysExceptionDispatchesTotal = "wmi_system_exception_dispatches_total"
	metricSysProcessorQueueLength     = "wmi_system_processor_queue_length"
	metricSysSystemCallsTotal         = "wmi_system_system_calls_total"
	metricSysSystemUpTime             = "wmi_system_system_up_time"
	metricSysThreads                  = "wmi_system_threads"
)

func (w *WMI) collectSystem(mx *metrics, pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorSystem)
	if !(enabled && success) {
		return false
	}
	mx.System = &system{}

	names := []string{
		metricSysContextSwitchesTotal,
		metricSysExceptionDispatchesTotal,
		metricSysProcessorQueueLength,
		metricSysSystemCallsTotal,
		metricSysSystemUpTime,
		metricSysThreads,
	}

	for _, name := range names {
		collectSystemAny(mx, pms, name)
	}

	mx.System.SystemUpTime = time.Now().Unix() - int64(mx.System.SystemBootTime)

	return true
}

func collectSystemAny(mx *metrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()

	switch name {
	default:
		panic(fmt.Sprintf("unknown metric name during system collection : %s", name))
	case metricSysContextSwitchesTotal:
		mx.System.ContextSwitchesTotal = value
	case metricSysExceptionDispatchesTotal:
		mx.System.ExceptionDispatchesTotal = value
	case metricSysProcessorQueueLength:
		mx.System.ProcessorQueueLength = value
	case metricSysSystemCallsTotal:
		mx.System.SystemCallsTotal = value
	case metricSysSystemUpTime:
		mx.System.SystemBootTime = value
	case metricSysThreads:
		mx.System.Threads = value
	}
}
