package wmi

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
)

const (
	metricNetBytesReceivedTotal       = "wmi_net_bytes_received_total"
	metricNetBytesSentTotal           = "wmi_net_bytes_sent_total"
	metricNetBytesTotal               = "wmi_net_bytes_total"
	metricNetPacketsOutboundDiscarded = "wmi_net_packets_outbound_discarded"
	metricNetPacketsOutboundErrors    = "wmi_net_packets_outbound_errors"
	metricNetPacketsReceivedDiscarded = "wmi_net_packets_received_discarded"
	metricNetPacketsReceivedErrors    = "wmi_net_packets_received_errors"
	metricNetPacketsReceivedTotal     = "wmi_net_packets_received_total"
	metricNetPacketsReceivedUnknown   = "wmi_net_packets_received_unknown"
	metricNetPacketsTotal             = "wmi_net_packets_total"
	metricNetPacketsSentTotal         = "wmi_net_packets_sent_total"
	metricNetCurrentBandwidth         = "wmi_net_current_bandwidth"
)

func (w *WMI) collectNet(mx *metrics, pms prometheus.Metrics) {
	names := []string{
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

	for _, name := range names {
		collectNetAny(mx, pms, name)
	}
}

func collectNetAny(mx *metrics, pms prometheus.Metrics, name string) {
	nic := newNIC("")

	for _, pm := range pms.FindByName(name) {
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
		switch name {
		default:
			panic(fmt.Sprintf("unknown metric name during net collection : %s", name))
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
}
