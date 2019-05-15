package wmi

import (
	"fmt"
	"strings"

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
	nic := newNIC("")

	for _, pm := range pms.FindByName(metricName) {
		var (
			nicID = pm.Labels.Get("nic")
			value = pm.Value
		)
		if nicID == "" {
			continue
		}
		nicID = strings.Replace(nicID, "__", "_", -1)
		if nic.ID != nicID {
			nic = mx.Net.NICs.get(nicID, true)
		}
		switch metricName {
		default:
			panic(fmt.Sprintf("unknown metric name during net collection : %s", metricName))
		case metricNetBytesReceived:
			nic.BytesReceivedTotal = value
		case metricNetBytesSent:
			nic.BytesSentTotal = value
		case metricNetBytes:
			nic.BytesTotal = value
		case metricNetPacketsOutboundDiscarded:
			nic.PacketsOutboundDiscarded = value
		case metricNetPacketsOutboundErrors:
			nic.PacketsOutboundErrors = value
		case metricNetPacketsReceivedDiscarded:
			nic.PacketsReceivedDiscarded = value
		case metricNetPacketsReceivedErrors:
			nic.PacketsReceivedErrors = value
		case metricNetPacketsReceived:
			nic.PacketsReceivedTotal = value
		case metricNetPacketsReceivedUnknown:
			nic.PacketsReceivedUnknown = value
		case metricNetPackets:
			nic.PacketsTotal = value
		case metricNetPacketsSent:
			nic.PacketsSentTotal = value
		case metricNetCurrentBandwidth:
			nic.CurrentBandwidth = value
		}
	}
}
