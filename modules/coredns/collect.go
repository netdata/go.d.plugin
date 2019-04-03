package coredns

import (
	"fmt"

	"github.com/netdata/go.d.plugin/pkg/prometheus"
	"github.com/netdata/go.d.plugin/pkg/stm"
)

func (cd *CoreDNS) collect() (map[string]int64, error) {
	raw, err := cd.prom.Scrape()

	if err != nil {
		return nil, err
	}

	mx := newMetrics()

	cd.collectPanicTotal(raw, mx)
	cd.collectRequestTotal(raw, mx)
	cd.collectRequestByTypeTotal(raw, mx)

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanicTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_panic_count_total"

	mx.Summary.Panic.Set(raw.FindByName(metricName).Max())
}

func (cd *CoreDNS) collectRequestTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_request_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			family = metric.Labels.Get("family")
			proto  = metric.Labels.Get("proto")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		// TODO: server can be empty
		if family == "" || proto == "" || zone == "" || server == "" {
			continue
		}

		mx.Summary.Request.Total.Add(value)

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &serverMetrics{}
		}

		srvMX := mx.PerServer[server]

		switch family {
		case "1":
			mx.Summary.Request.ByIPFamily.IPv4.Add(value)
			srvMX.Request.ByIPFamily.IPv4.Add(value)
		case "2":
			mx.Summary.Request.ByIPFamily.IPv6.Add(value)
			srvMX.Request.ByIPFamily.IPv6.Add(value)
		}

		switch proto {
		case "udp":
			mx.Summary.Request.ByProto.UDP.Add(value)
			srvMX.Request.ByProto.UDP.Add(value)
		case "tcp":
			mx.Summary.Request.ByProto.UDP.Add(value)
			srvMX.Request.ByProto.TCP.Add(value)
		}

		switch zone {
		default:
			mx.Summary.Request.ByStatus.Processed.Add(value)
			srvMX.Request.ByStatus.Processed.Add(value)
		case "dropped":
			mx.Summary.Request.ByStatus.Dropped.Add(value)
			srvMX.Request.ByStatus.Dropped.Add(value)
		}
	}
}

func (cd *CoreDNS) collectRequestByTypeTotal(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_request_type_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			typ    = metric.Labels.Get("type")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		// TODO: server can be empty??
		if typ == "" || server == "" {
			continue
		}

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &serverMetrics{}
		}

		srvMX := mx.PerServer[server]

		switch typ {
		case "A":
			mx.Summary.Request.ByType.A.Add(value)
			srvMX.Request.ByType.A.Add(value)
		case "AAAA":
			mx.Summary.Request.ByType.AAAA.Add(value)
			srvMX.Request.ByType.AAAA.Add(value)
		case "MX":
			mx.Summary.Request.ByType.MX.Add(value)
			srvMX.Request.ByType.MX.Add(value)
		case "SOA":
			mx.Summary.Request.ByType.SOA.Add(value)
			srvMX.Request.ByType.SOA.Add(value)
		case "CNAME":
			mx.Summary.Request.ByType.CNAME.Add(value)
			srvMX.Request.ByType.CNAME.Add(value)
		case "PTR":
			mx.Summary.Request.ByType.PTR.Add(value)
			srvMX.Request.ByType.PTR.Add(value)
		case "TXT":
			mx.Summary.Request.ByType.TXT.Add(value)
			srvMX.Request.ByType.TXT.Add(value)
		case "NS":
			mx.Summary.Request.ByType.NS.Add(value)
			srvMX.Request.ByType.NS.Add(value)
		case "DS":
			mx.Summary.Request.ByType.DS.Add(value)
			srvMX.Request.ByType.DS.Add(value)
		case "DNSKEY":
			mx.Summary.Request.ByType.DNSKEY.Add(value)
			srvMX.Request.ByType.DNSKEY.Add(value)
		case "RRSIG":
			mx.Summary.Request.ByType.RRSIG.Add(value)
			srvMX.Request.ByType.RRSIG.Add(value)
		case "NSEC":
			mx.Summary.Request.ByType.NSEC.Add(value)
			srvMX.Request.ByType.NSEC.Add(value)
		case "NSEC3":
			mx.Summary.Request.ByType.NSEC3.Add(value)
			srvMX.Request.ByType.NSEC3.Add(value)
		case "IXFR":
			mx.Summary.Request.ByType.IXFR.Add(value)
			srvMX.Request.ByType.IXFR.Add(value)
		case "ANY":
			mx.Summary.Request.ByType.ANY.Add(value)
			srvMX.Request.ByType.ANY.Add(value)
		case "other":
			mx.Summary.Request.ByType.Other.Add(value)
			srvMX.Request.ByType.Other.Add(value)
		}
	}
}

func (cd *CoreDNS) addNewServerCharts(serverName string) {
	charts := serverCharts.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, serverName)
		chart.Fam = fmt.Sprintf(chart.Fam, serverName)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, serverName)
		}
	}
	_ = cd.charts.Add(*charts...)
}

//func (cd CoreDNS) collectRequestTotalPerFamily(raw prometheus.Metrics, mx *metrics) {
//	metricName := "coredns_dns_request_count_total"
//
//	for _, metric := range raw.FindByName(metricName) {
//		mx.Request.Count.Total.Add(metric.Value)
//	}
//}
//
//func (cd *CoreDNS) collectRequestsByTypeTotal(raw prometheus.Metrics, mx *metrics) {
//	metricName := "coredns_dns_request_type_count_total"
//	chartName := "request_type_count_total"
//	chart := cd.charts.Get(chartName)
//
//	for _, metric := range raw.FindByName(metricName) {
//		typ := metric.Labels.Get("type")
//		if typ == "" {
//			continue
//		}
//		if chart == nil {
//			_ = cd.charts.Add(chartReqByTypeTotal.Copy())
//			chart = cd.charts.Get(chartName)
//		}
//		dimID := "request_count_by_type_total_" + typ
//		if !chart.HasDim(dimID) {
//			_ = chart.AddDim(&Dim{ID: dimID, Name: typ, Algo: module.Incremental})
//		}
//
//		current := mx.Request.Count.ByTypeTotal[typ].Value()
//		mx.Request.Count.ByTypeTotal[typ] = mtx.Gauge(metric.Value + current)
//	}
//}
//
//func (cd *CoreDNS) collectResponsesByRcodeTotal(raw prometheus.Metrics, mx *metrics) {
//	metricName := "coredns_dns_response_rcode_count_total"
//	chartName := "response_rcode_count_total"
//	chart := cd.charts.Get(chartName)
//
//	for _, metric := range raw.FindByName(metricName) {
//		rcode := metric.Labels.Get("rcode")
//		if rcode == "" {
//			continue
//		}
//		if chart == nil {
//			_ = cd.charts.Add(chartRespByRcodeTotal.Copy())
//			chart = cd.charts.Get(chartName)
//		}
//		dimID := "response_count_by_rcode_total_" + rcode
//		if !chart.HasDim(dimID) {
//			_ = chart.AddDim(&Dim{ID: dimID, Name: rcode, Algo: module.Incremental})
//		}
//
//		current := mx.Response.Count.ByRcodeTotal[rcode].Value()
//		mx.Response.Count.ByRcodeTotal[rcode] = mtx.Gauge(metric.Value + current)
//	}
//}
