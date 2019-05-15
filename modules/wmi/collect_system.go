package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	metricSysContextSwitchesTotal     = "wmi_system_context_switches_total"
	metricSysExceptionDispatchesTotal = "wmi_system_exception_dispatches_total"
	metricSysProcessorQueueLength     = "wmi_system_processor_queue_length"
	metricSysSystemCallsTotal         = "wmi_system_system_calls_total"
	metricSysSystemUpTime             = "wmi_system_system_up_time"
	metricSysThreads                  = "wmi_system_threads"
)

func (w *WMI) collectSystem(mx *metrics, pms prometheus.Metrics) {
	mx.System.ContextSwitchesTotal = pms.FindByName(metricSysContextSwitchesTotal).Max()
	mx.System.ExceptionDispatchesTotal = pms.FindByName(metricSysExceptionDispatchesTotal).Max()
	mx.System.ProcessorQueueLength = pms.FindByName(metricSysProcessorQueueLength).Max()
	mx.System.SystemCallsTotal = pms.FindByName(metricSysSystemCallsTotal).Max()
	mx.System.SystemUpTime = pms.FindByName(metricSysSystemUpTime).Max()
	mx.System.Threads = pms.FindByName(metricSysThreads).Max()
}
