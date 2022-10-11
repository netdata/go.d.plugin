// SPDX-License-Identifier: GPL-3.0-or-later

package dnsquery

import (
	"fmt"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"
)

const (
	prioDNSQueryTime = module.Priority + iota
)

var (
	dnsChartsTmpl = module.Charts{
		dnsQueryTimeChartTmpl.Copy(),
	}
	dnsQueryTimeChartTmpl = module.Chart{
		ID:       "server_%s_record_%s_query_time",
		Title:    "DNS Query Time",
		Units:    "seconds",
		Fam:      "query time",
		Ctx:      "dns_query.query_time",
		Priority: prioDNSQueryTime,
		Dims: module.Dims{
			{ID: "server_%s_record_%s_query_time", Name: "query_time", Div: 1e9},
		},
	}
)

func newDNSServerCharts(server, network, rtype string) *module.Charts {
	charts := dnsChartsTmpl.Copy()

	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, strings.ReplaceAll(server, ".", "_"), rtype)
		chart.Labels = []module.Label{
			{Key: "server", Value: server},
			{Key: "network", Value: network},
			{Key: "record_type", Value: rtype},
		}
		for _, d := range chart.Dims {
			d.ID = fmt.Sprintf(d.ID, server, rtype)
		}
	}

	return charts
}
