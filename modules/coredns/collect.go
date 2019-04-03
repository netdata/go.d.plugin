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

	cd.collectPanic(raw, mx)
	cd.collectRequest(raw, mx)
	cd.collectRequestByType(raw, mx)
	cd.collectResponseByRcode(raw, mx)

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanic(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_panic_count_total"

	mx.Summary.Panic.Set(raw.FindByName(metricName).Max())
}

func (cd *CoreDNS) collectRequest(raw prometheus.Metrics, mx *metrics) {
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

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.Servers[server] == nil {
			mx.Servers[server] = &serverMetrics{}
		}

		srvMX := mx.Servers[server]

		mx.Summary.Request.Total.Add(value)
		srvMX.Request.Total.Add(value)

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

			mx.Summary.Request.ByStatus.Processed.Sub(value)
			srvMX.Request.ByStatus.Processed.Sub(value)
		}
	}
}

func (cd *CoreDNS) collectRequestByType(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_request_type_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			typ    = metric.Labels.Get("type")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		// TODO: server can be empty??
		if typ == "" || zone == "" || server == "" {
			continue
		}

		if zone == "dropped" {
			continue
		}

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.Servers[server] == nil {
			mx.Servers[server] = &serverMetrics{}
		}

		srvMX := mx.Servers[server]

		switch typ {
		default:
			mx.Summary.Request.ByType.Other.Add(value)
			srvMX.Request.ByType.Other.Add(value)
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
		}
	}
}

func (cd *CoreDNS) collectResponseByRcode(raw prometheus.Metrics, mx *metrics) {
	metricName := "coredns_dns_response_rcode_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			rcode  = metric.Labels.Get("rcode")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		// TODO: server can be empty??
		if rcode == "" || zone == "" || server == "" {
			continue
		}

		if zone == "dropped" {
			continue
		}

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.Servers[server] == nil {
			mx.Servers[server] = &serverMetrics{}
		}

		srvMX := mx.Servers[server]

		mx.Summary.Response.Total.Add(value)
		srvMX.Response.Total.Add(value)

		switch rcode {
		default:
			mx.Summary.Response.ByRcode.Other.Add(value)
			srvMX.Response.ByRcode.Other.Add(value)
		case "NOERROR":
			mx.Summary.Response.ByRcode.NOERROR.Add(value)
			srvMX.Response.ByRcode.NOERROR.Add(value)
		case "FORMERR":
			mx.Summary.Response.ByRcode.FORMERR.Add(value)
			srvMX.Response.ByRcode.FORMERR.Add(value)
		case "SERVFAIL":
			mx.Summary.Response.ByRcode.SERVFAIL.Add(value)
			srvMX.Response.ByRcode.SERVFAIL.Add(value)
		case "NXDOMAIN":
			mx.Summary.Response.ByRcode.NXDOMAIN.Add(value)
			srvMX.Response.ByRcode.NXDOMAIN.Add(value)
		case "NOTIMP":
			mx.Summary.Response.ByRcode.NOTIMP.Add(value)
			srvMX.Response.ByRcode.NOTIMP.Add(value)
		case "REFUSED":
			mx.Summary.Response.ByRcode.REFUSED.Add(value)
			srvMX.Response.ByRcode.REFUSED.Add(value)
		case "YXDOMAIN":
			mx.Summary.Response.ByRcode.YXDOMAIN.Add(value)
			srvMX.Response.ByRcode.YXDOMAIN.Add(value)
		case "YXRRSET":
			mx.Summary.Response.ByRcode.YXRRSET.Add(value)
			srvMX.Response.ByRcode.YXRRSET.Add(value)
		case "NXRRSET":
			mx.Summary.Response.ByRcode.NXRRSET.Add(value)
			srvMX.Response.ByRcode.NXRRSET.Add(value)
		case "NOTAUTH":
			mx.Summary.Response.ByRcode.NOTAUTH.Add(value)
			srvMX.Response.ByRcode.NOTAUTH.Add(value)
		case "NOTZONE":
			mx.Summary.Response.ByRcode.NOTZONE.Add(value)
			srvMX.Response.ByRcode.NOTZONE.Add(value)
		case "BADSIG":
			mx.Summary.Response.ByRcode.BADSIG.Add(value)
			srvMX.Response.ByRcode.BADSIG.Add(value)
		case "BADKEY":
			mx.Summary.Response.ByRcode.BADKEY.Add(value)
			srvMX.Response.ByRcode.BADKEY.Add(value)
		case "BADTIME":
			mx.Summary.Response.ByRcode.BADTIME.Add(value)
			srvMX.Response.ByRcode.BADTIME.Add(value)
		case "BADMODE":
			mx.Summary.Response.ByRcode.BADMODE.Add(value)
			srvMX.Response.ByRcode.BADMODE.Add(value)
		case "BADNAME":
			mx.Summary.Response.ByRcode.BADNAME.Add(value)
			srvMX.Response.ByRcode.BADNAME.Add(value)
		case "BADALG":
			mx.Summary.Response.ByRcode.BADALG.Add(value)
			srvMX.Response.ByRcode.BADALG.Add(value)
		case "BADTRUNC":
			mx.Summary.Response.ByRcode.BADTRUNC.Add(value)
			srvMX.Response.ByRcode.BADTRUNC.Add(value)
		case "BADCOOKIE":
			mx.Summary.Response.ByRcode.BADCOOKIE.Add(value)
			srvMX.Response.ByRcode.BADCOOKIE.Add(value)
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
