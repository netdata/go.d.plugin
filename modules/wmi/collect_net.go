package wmi

import (
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	collectorNet = "net"

	metricNetBytesReceivedTotal       = "windows_net_bytes_received_total"
	metricNetBytesSentTotal           = "windows_net_bytes_sent_total"
	metricNetBytesTotal               = "windows_net_bytes_total"
	metricNetPacketsOutboundDiscarded = "windows_net_packets_outbound_discarded"
	metricNetPacketsOutboundErrors    = "windows_net_packets_outbound_errors"
	metricNetPacketsReceivedDiscarded = "windows_net_packets_received_discarded"
	metricNetPacketsReceivedErrors    = "windows_net_packets_received_errors"
	metricNetPacketsReceivedTotal     = "windows_net_packets_received_total"
	metricNetPacketsReceivedUnknown   = "windows_net_packets_received_unknown"
	metricNetPacketsTotal             = "windows_net_packets_total"
	metricNetPacketsSentTotal         = "windows_net_packets_sent_total"
	metricNetCurrentBandwidth         = "windows_net_current_bandwidth"
)

var netMetricNames = []string{
	metricNetBytesReceivedTotal,
	metricNetBytesSentTotal,
	metricNetBytesTotal,
	metricNetPacketsOutboundDiscarded,
	metricNetPacketsOutboundErrors,
	metricNetPacketsReceivedDiscarded,
	metricNetPacketsReceivedErrors,
	metricNetPacketsReceivedTotal,
	metricNetPacketsReceivedUnknown,
	metricNetPacketsTotal,
	metricNetPacketsSentTotal,
	metricNetCurrentBandwidth,
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
	case metricNetPacketsOutboundDiscarded:
		nic.PacketsOutboundDiscarded = value
	case metricNetPacketsOutboundErrors:
		nic.PacketsOutboundErrors = value
	case metricNetPacketsReceivedDiscarded:
		nic.PacketsReceivedDiscarded = value
	case metricNetPacketsReceivedErrors:
		nic.PacketsReceivedErrors = value
	case metricNetPacketsReceivedTotal:
		nic.PacketsReceivedTotal = value
	case metricNetPacketsReceivedUnknown:
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
