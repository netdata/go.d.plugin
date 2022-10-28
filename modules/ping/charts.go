// SPDX-License-Identifier: GPL-3.0-or-later

package ping

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioHostRTTLatency = module.Priority + iota
	prioHostPingPacketLoss
	prioHostPingPackets
)

var hostChartsTmpl = module.Charts{
	hostRTTChartTmpl.Copy(),
	hostPacketLossChartTmpl.Copy(),
	hostPacketsChartTmpl.Copy(),
}

var hostRTTChartTmpl = module.Chart{
	ID:       "ping_host_%s_rtt",
	Title:    "Ping round-trip time",
	Units:    "seconds",
	Fam:      "latency",
	Ctx:      "ping.host_rtt",
	Priority: prioHostRTTLatency,
	Type:     module.Area,
	Dims: module.Dims{
		{ID: "host_%s_min_rtt", Name: "min", Div: 1e6},
		{ID: "host_%s_max_rtt", Name: "max", Div: 1e6},
		{ID: "host_%s_avg_rtt", Name: "avg", Div: 1e6},
	},
}

var hostPacketLossChartTmpl = module.Chart{
	ID:       "host_host_%s_packet_loss",
	Title:    "Ping packet loss",
	Units:    "percentage",
	Fam:      "packet loss",
	Ctx:      "ping.host_packet_loss",
	Priority: prioHostPingPacketLoss,
	Dims: module.Dims{
		{ID: "host_%s_packet_loss", Name: "loss", Div: 1000},
	},
}

var hostPacketsChartTmpl = module.Chart{
	ID:       "host_host_%s_packets",
	Title:    "Ping packets transferred",
	Units:    "packets",
	Fam:      "packets",
	Ctx:      "ping.host_packets",
	Priority: prioHostPingPackets,
	Dims: module.Dims{
		{ID: "host_%s_packets_recv", Name: "received"},
		{ID: "host_%s_packets_sent", Name: "sent"},
	},
}

func newHostCharts(host string) *module.Charts {
	charts := hostChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, strings.ReplaceAll(host, ".", "_"))
		chart.Labels = []module.Label{
			{Key: "host", Value: host},
		}
		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, host)
		}
	}

	return charts
}

func (p *Ping) addHostCharts(host string) {
	charts := newHostCharts(host)

	if err := p.Charts().Add(*charts...); err != nil {
		p.Warning(err)
	}
}
