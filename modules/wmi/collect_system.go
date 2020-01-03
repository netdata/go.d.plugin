package wmi

import (
	"time"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorSystem = "system"

	metricSysContextSwitchesTotal     = "wmi_system_context_switches_total"
	metricSysExceptionDispatchesTotal = "wmi_system_exception_dispatches_total"
	metricSysProcessorQueueLength     = "wmi_system_processor_queue_length"
	metricSysSystemCallsTotal         = "wmi_system_system_calls_total"
	metricSysSystemUpTime             = "wmi_system_system_up_time"
	metricSysThreads                  = "wmi_system_threads"
)

var systemMetricsNames = []string{
	metricSysContextSwitchesTotal,
	metricSysExceptionDispatchesTotal,
	metricSysProcessorQueueLength,
	metricSysSystemCallsTotal,
	metricSysSystemUpTime,
	metricSysThreads,
}

func doCollectSystem(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorSystem)
	return enabled && success
}

func collectSystem(pms prometheus.Metrics) *systemMetrics {
	if !doCollectSystem(pms) {
		return nil
	}

	sm := &systemMetrics{}
	for _, name := range systemMetricsNames {
		collectSystemMetric(sm, pms, name)
	}
	sm.SystemUpTime = time.Now().Unix() - int64(sm.SystemBootTime)
	return sm
}

func collectSystemMetric(sm *systemMetrics, pms prometheus.Metrics, name string) {
	value := pms.FindByName(name).Max()
	assignSystemMetric(sm, name, value)
}

func assignSystemMetric(sm *systemMetrics, name string, value float64) {
	switch name {
	case metricSysContextSwitchesTotal:
		sm.ContextSwitchesTotal = value
	case metricSysExceptionDispatchesTotal:
		sm.ExceptionDispatchesTotal = value
	case metricSysProcessorQueueLength:
		sm.ProcessorQueueLength = value
	case metricSysSystemCallsTotal:
		sm.SystemCallsTotal = value
	case metricSysSystemUpTime:
		sm.SystemBootTime = value
	case metricSysThreads:
		sm.Threads = value
	}
}
