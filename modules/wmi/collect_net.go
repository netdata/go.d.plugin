// SPDX-License-Identifier: GPL-3.0-or-later

package wmi

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorNet = "net"

	metricNetCurrentBandwidth = "windows_net_current_bandwidth"

	metricNetBytesTotal         = "windows_net_bytes_total"
	metricNetBytesReceivedTotal = "windows_net_bytes_received_total"
	metricNetBytesSentTotal     = "windows_net_bytes_sent_total"

	metricNetPacketsTotal         = "windows_net_packets_total"
	metricNetPacketsReceivedTotal = "windows_net_packets_received_total"
	metricNetPacketsSentTotal     = "windows_net_packets_sent_total"

	// up to v0.15.0 (TODO: should be removed)
	metricNetPacketsReceivedDiscarded = "windows_net_packets_received_discarded"
	metricNetPacketsOutboundDiscarded = "windows_net_packets_outbound_discarded"
	metricNetPacketsReceivedErrors    = "windows_net_packets_received_errors"
	metricNetPacketsOutboundErrors    = "windows_net_packets_outbound_errors"
	metricNetPacketsReceivedUnknown   = "windows_net_packets_received_unknown"

	// v0.16.0+
	metricNetPacketsReceivedDiscardedTotal = "windows_net_packets_received_discarded_total"
	metricNetPacketsOutboundDiscardedTotal = "windows_net_packets_outbound_discarded_total"
	metricNetPacketsReceivedErrorsTotal    = "windows_net_packets_received_errors_total"
	metricNetPacketsOutboundErrorsTotal    = "windows_net_packets_outbound_errors_total"
	metricNetPacketsReceivedUnknownTotal   = "windows_net_packets_received_unknown_total"
)

var netMetricNames = []string{
	metricNetCurrentBandwidth,

	metricNetBytesTotal,
	metricNetBytesReceivedTotal,
	metricNetBytesSentTotal,

	metricNetPacketsTotal,
	metricNetPacketsReceivedTotal,
	metricNetPacketsSentTotal,

	metricNetPacketsReceivedDiscarded,
	metricNetPacketsOutboundDiscarded,
	metricNetPacketsReceivedErrors,
	metricNetPacketsOutboundErrors,
	metricNetPacketsReceivedUnknown,

	metricNetPacketsReceivedDiscardedTotal,
	metricNetPacketsOutboundDiscardedTotal,
	metricNetPacketsReceivedErrorsTotal,
	metricNetPacketsOutboundErrorsTotal,
	metricNetPacketsReceivedUnknownTotal,
}

func doCollectNet(pms prometheus.Metrics) bool {
	enabled, success := checkCollector(pms, collectorNet)
	return enabled && success
}

func collectNet(pms prometheus.Metrics) *networkMetrics {
	if !doCollectNet(pms) {
		return nil
	}

	nm := &networkMetrics{}
	for _, name := range netMetricNames {
		collectNetMetric(nm, pms, name)
	}
	return nm
}

func collectNetMetric(nm *networkMetrics, pms prometheus.Metrics, name string) {
	var nic *netNIC

	for _, pm := range pms.FindByName(name) {
		nicID := pm.Labels.Get("nic")
		if nicID == "" {
			continue
		}

		nicID = cleanNICID(nicID)
		if nic == nil || nic.ID != nicID {
			nic = nm.NICs.get(nicID)
		}

		assignNICMetric(nic, name, pm.Value)
	}
}

func assignNICMetric(nic *netNIC, name string, value float64) {
	switch name {
	case metricNetBytesReceivedTotal:
		nic.BytesReceivedTotal = value
	case metricNetBytesSentTotal:
		nic.BytesSentTotal = value
	case metricNetBytesTotal:
		nic.BytesTotal = value
	case metricNetPacketsOutboundDiscarded, metricNetPacketsOutboundDiscardedTotal:
		nic.PacketsOutboundDiscarded = value
	case metricNetPacketsOutboundErrors, metricNetPacketsOutboundErrorsTotal:
		nic.PacketsOutboundErrors = value
	case metricNetPacketsReceivedDiscarded, metricNetPacketsReceivedDiscardedTotal:
		nic.PacketsReceivedDiscarded = value
	case metricNetPacketsReceivedErrors, metricNetPacketsReceivedErrorsTotal:
		nic.PacketsReceivedErrors = value
	case metricNetPacketsReceivedTotal:
		nic.PacketsReceivedTotal = value
	case metricNetPacketsReceivedUnknown, metricNetPacketsReceivedUnknownTotal:
		nic.PacketsReceivedUnknown = value
	case metricNetPacketsTotal:
		nic.PacketsTotal = value
	case metricNetPacketsSentTotal:
		nic.PacketsSentTotal = value
	case metricNetCurrentBandwidth:
		nic.CurrentBandwidth = value
	}
}

func cleanNICID(id string) string {
	return strings.Replace(id, "__", "_", -1)
}
