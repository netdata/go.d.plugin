package wmi

import (
	"fmt"
	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricNetBytesReceived            = "wmi_net_bytes_received_total"
	metricNetBytesSent                = "wmi_net_bytes_sent_total"
	metricNetBytes                    = "wmi_net_bytes_total"
	metricNetPacketsOutboundDiscarded = "wmi_net_packets_outbound_discarded"
	metricNetPacketsOutboundErrors    = "wmi_net_packets_outbound_errors"
	metricNetPacketsReceivedDiscarded = "wmi_net_packets_received_discarded"
	metricNetPacketsReceivedErrors    = "wmi_net_packets_received_errors"
	metricNetPacketsReceived          = "wmi_net_packets_received_total"
	metricNetPacketsReceivedUnknown   = "wmi_net_packets_received_unknown"
	metricNetPackets                  = "wmi_net_packets_total"
	metricNetPacketsSent              = "wmi_net_packets_sent_total"
	metricNetCurrentBandwidth         = "wmi_net_current_bandwidth"
)

func (w *WMI) collectNet(mx *metrics, pms prometheus.Metrics) {
	names := []string{
		metricNetBytesReceived,
		metricNetBytesSent,
		metricNetBytes,
		metricNetPacketsOutboundDiscarded,
		metricNetPacketsOutboundErrors,
		metricNetPacketsReceivedDiscarded,
		metricNetPacketsReceivedErrors,
		metricNetPacketsReceived,
		metricNetPacketsReceivedUnknown,
		metricNetPackets,
		metricNetPacketsSent,
		metricNetCurrentBandwidth,
	}

	for _, name := range names {
		collectNetAny(mx, pms, name)
	}
}

func collectNetAny(mx *metrics, pms prometheus.Metrics, metricName string) {
	n := newNIC("")

	for _, pm := range pms.FindByName(metricName) {
		var (
			nicID = pm.Labels.Get("nic")
			value = pm.Value
		)
		if nicID == "" {
			continue
		}
		if n.ID != nicID {
			n = mx.Net.NICs.get(nicID, true)
		}
		switch metricName {
		default:
			panic(fmt.Sprintf("unknown metric name during net collection : %s", metricName))
		case metricNetBytesReceived:
			n.BytesReceivedTotal.Set(value)
		case metricNetBytesSent:
			n.BytesSentTotal.Set(value)
		case metricNetBytes:
			n.BytesTotal.Set(value)
		case metricNetPacketsOutboundDiscarded:
			n.PacketsOutboundDiscarded.Set(value)
		case metricNetPacketsOutboundErrors:
			n.PacketsOutboundErrors.Set(value)
		case metricNetPacketsReceivedDiscarded:
			n.PacketsReceivedDiscarded.Set(value)
		case metricNetPacketsReceivedErrors:
			n.PacketsReceivedErrors.Set(value)
		case metricNetPacketsReceived:
			n.PacketsReceivedTotal.Set(value)
		case metricNetPacketsReceivedUnknown:
			n.PacketsReceivedUnknown.Set(value)
		case metricNetPackets:
			n.PacketsTotal.Set(value)
		case metricNetPacketsSent:
			n.PacketsSentTotal.Set(value)
		case metricNetCurrentBandwidth:
			n.CurrentBandwidth.Set(value)
		}
	}
}
