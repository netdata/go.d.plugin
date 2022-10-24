// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import "github.com/netdata/go.d.plugin/pkg/prometheus"

const (
	collectorTCP = "tcp"

	afIPV4 = "ipv4"
	afIPV6 = "ipv6"

	metricTCPConnectionFailure               = "windows_tcp_connection_failures_total"
	metricTCPConnectionActive                = "windows_tcp_connections_active_total"
	metricTCPConnectionEstablished           = "windows_tcp_connections_established"
	metricTCPConnectionPassive               = "windows_tcp_connections_passive_total"
	metricTCPConnectionReset                 = "windows_tcp_connections_reset_total"
	metricTCPConnectionSegmentsReceived      = "windows_tcp_segments_received_total"
	metricTCPConnectionSegmentsRetransmitted = "windows_tcp_segments_retransmitted_total"
	metricTCPConnectionSegmentsSent          = "windows_tcp_segments_sent_total"
)

var tcpMetricNames = []string{
	metricTCPConnectionFailure,
	metricTCPConnectionActive,
	metricTCPConnectionEstablished,
	metricTCPConnectionPassive,
	metricTCPConnectionReset,
	metricTCPConnectionSegmentsReceived,
	metricTCPConnectionSegmentsRetransmitted,
	metricTCPConnectionSegmentsSent,
}

func doCollectTCP(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorTCP)
	return enabled && success
}

func collectTCP(pms prometheus.Metrics) *tcpMetrics {
	if !doCollectTCP(pms) {
		return nil
	}

	tcpm := &tcpMetrics{}
	for _, name := range tcpMetricNames {
		collectTCPMetric(tcpm, pms, name)
	}

	return tcpm
}

func collectTCPMetric(tcpm *tcpMetrics, pms prometheus.Metrics, name string) {
	for _, pm := range pms.FindByName(name) {
		af := pm.Labels.Get("af")

		assignTCPMetric(tcpm, af, name, pm.Value)
	}
}

func assignTCPMetric(tcpm *tcpMetrics, af string, name string, value float64) {
	switch name {
	case metricTCPConnectionFailure:
		assignTCPConnection(&tcpm.failures, af, value)
	case metricTCPConnectionActive:
		assignTCPConnection(&tcpm.active, af, value)
	case metricTCPConnectionEstablished:
		assignTCPConnection(&tcpm.established, af, value)
	case metricTCPConnectionPassive:
		assignTCPConnection(&tcpm.passive, af, value)
	case metricTCPConnectionReset:
		assignTCPConnection(&tcpm.reset, af, value)
	case metricTCPConnectionSegmentsReceived:
		assignTCPConnection(&tcpm.segmentsReceived, af, value)
	case metricTCPConnectionSegmentsRetransmitted:
		assignTCPConnection(&tcpm.segmentsRetransmitted, af, value)
	case metricTCPConnectionSegmentsSent:
		assignTCPConnection(&tcpm.segmentsSent, af, value)
	}
}

func assignTCPConnection(c *tcpConnectionAF, af string, value float64) {
	switch af {
	case afIPV4:
		c.ipv4 = value
	case afIPV6:
		c.ipv6 = value
	}
}
