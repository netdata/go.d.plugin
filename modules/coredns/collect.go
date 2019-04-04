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

	cd.collectPanic(mx, raw)
	cd.collectSummaryRequests(mx, raw)
	cd.collectSummaryRequestByType(mx, raw)
	cd.collectSummaryResponseByRcode(mx, raw)

	cd.collectPerServerRequests(mx, raw)
	cd.collectPerServerRequestByType(mx, raw)
	cd.collectPerServerResponseByRcode(mx, raw)

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanic(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_panic_count_total"

	mx.Panic.Set(raw.FindByName(metricName).Max())
}

func (cd *CoreDNS) collectSummaryRequests(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_request_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			family = metric.Labels.Get("family")
			proto  = metric.Labels.Get("proto")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if family == "" || proto == "" || zone == "" {
			continue
		}

		if server == "" {
			mx.NoServerDropped.Add(value)
			continue
		}

		setRequestTotal(&mx.Summary.Request, value, family, proto, zone)
	}
}

func (cd *CoreDNS) collectPerServerRequests(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_request_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			family = metric.Labels.Get("family")
			proto  = metric.Labels.Get("proto")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if family == "" || proto == "" || zone == "" || server == "" {
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

		setRequestTotal(&srvMX.Request, value, family, proto, zone)
	}
}

func (cd *CoreDNS) collectSummaryRequestByType(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_request_type_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			typ    = metric.Labels.Get("type")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if typ == "" || zone == "" || zone == "dropped" || server == "" {
			continue
		}

		setRequestByType(&mx.Summary.Request, value, typ)
	}
}

func (cd *CoreDNS) collectPerServerRequestByType(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_request_type_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			typ    = metric.Labels.Get("type")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if typ == "" || zone == "" || zone == "dropped" || server == "" {
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

		setRequestByType(&srvMX.Request, value, typ)
	}
}

func (cd *CoreDNS) collectSummaryResponseByRcode(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_response_rcode_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			rcode  = metric.Labels.Get("rcode")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if rcode == "" || zone == "" || zone == "dropped" || server == "" {
			continue
		}

		setResponseByRcode(&mx.Summary.Response, value, rcode)
	}
}

func (cd *CoreDNS) collectPerServerResponseByRcode(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_response_rcode_count_total"

	for _, metric := range raw.FindByName(metricName) {
		var (
			rcode  = metric.Labels.Get("rcode")
			zone   = metric.Labels.Get("zone")
			server = metric.Labels.Get("server")
			value  = metric.Value
		)

		if rcode == "" || zone == "" || zone == "dropped" || server == "" {
			continue
		}

		if !cd.activeServers[server] {
			cd.addNewServerCharts(server)
			cd.activeServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &serverMetrics{}
		}

		setResponseByRcode(&mx.PerServer[server].Response, value, rcode)
	}
}

func setRequestTotal(mx *requestMetrics, value float64, family, proto, zone string) {
	mx.Total.Add(value)

	switch family {
	case "1":
		mx.ByIPFamily.IPv4.Add(value)
	case "2":
		mx.ByIPFamily.IPv6.Add(value)
	}
	switch proto {
	case "udp":
		mx.ByProto.UDP.Add(value)
	case "tcp":
		mx.ByProto.TCP.Add(value)
	}
	switch zone {
	default:
		mx.ByStatus.Processed.Add(value)
	case "dropped":
		mx.ByStatus.Dropped.Add(value)
		mx.ByStatus.Processed.Sub(value)
	}
}

func setRequestByType(mx *requestMetrics, value float64, typ string) {
	switch typ {
	default:
		mx.ByType.Other.Add(value)
	case "A":
		mx.ByType.A.Add(value)
	case "AAAA":
		mx.ByType.AAAA.Add(value)
	case "MX":
		mx.ByType.MX.Add(value)
	case "SOA":
		mx.ByType.SOA.Add(value)
	case "CNAME":
		mx.ByType.CNAME.Add(value)
	case "PTR":
		mx.ByType.PTR.Add(value)
	case "TXT":
		mx.ByType.TXT.Add(value)
	case "NS":
		mx.ByType.NS.Add(value)
	case "DS":
		mx.ByType.DS.Add(value)
	case "DNSKEY":
		mx.ByType.DNSKEY.Add(value)
	case "RRSIG":
		mx.ByType.RRSIG.Add(value)
	case "NSEC":
		mx.ByType.NSEC.Add(value)
	case "NSEC3":
		mx.ByType.NSEC3.Add(value)
	case "IXFR":
		mx.ByType.IXFR.Add(value)
	case "ANY":
		mx.ByType.ANY.Add(value)
	}
}

func setResponseByRcode(mx *responseMetrics, value float64, rcode string) {
	mx.Total.Add(value)

	switch rcode {
	default:
		mx.ByRcode.Other.Add(value)
	case "NOERROR":
		mx.ByRcode.NOERROR.Add(value)
	case "FORMERR":
		mx.ByRcode.FORMERR.Add(value)
	case "SERVFAIL":
		mx.ByRcode.SERVFAIL.Add(value)
	case "NXDOMAIN":
		mx.ByRcode.NXDOMAIN.Add(value)
	case "NOTIMP":
		mx.ByRcode.NOTIMP.Add(value)
	case "REFUSED":
		mx.ByRcode.REFUSED.Add(value)
	case "YXDOMAIN":
		mx.ByRcode.YXDOMAIN.Add(value)
	case "YXRRSET":
		mx.ByRcode.YXRRSET.Add(value)
	case "NXRRSET":
		mx.ByRcode.NXRRSET.Add(value)
	case "NOTAUTH":
		mx.ByRcode.NOTAUTH.Add(value)
	case "NOTZONE":
		mx.ByRcode.NOTZONE.Add(value)
	case "BADSIG":
		mx.ByRcode.BADSIG.Add(value)
	case "BADKEY":
		mx.ByRcode.BADKEY.Add(value)
	case "BADTIME":
		mx.ByRcode.BADTIME.Add(value)
	case "BADMODE":
		mx.ByRcode.BADMODE.Add(value)
	case "BADNAME":
		mx.ByRcode.BADNAME.Add(value)
	case "BADALG":
		mx.ByRcode.BADALG.Add(value)
	case "BADTRUNC":
		mx.ByRcode.BADTRUNC.Add(value)
	case "BADCOOKIE":
		mx.ByRcode.BADCOOKIE.Add(value)
	}
}

func (cd *CoreDNS) addNewServerCharts(name string) {
	charts := serverCharts.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, name)
		chart.Fam = fmt.Sprintf(chart.Fam, name)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, name)
		}
	}
	_ = cd.charts.Add(*charts...)
}
