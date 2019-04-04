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
	cd.collectSummaryRequestsDuration(mx, raw)

	if cd.perServerMatcher != nil {
		cd.collectPerServerRequests(mx, raw)
		cd.collectPerServerRequestByType(mx, raw)
		cd.collectPerServerResponseByRcode(mx, raw)
	}

	return stm.ToMap(mx), nil
}

func (cd CoreDNS) collectPanic(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_panic_count_total"

	mx.Panic.Set(raw.FindByName(metricName).Max())
}

func (cd *CoreDNS) collectSummaryRequestsDuration(mx *metrics, raw prometheus.Metrics) {
	metricName := "coredns_dns_request_duration_seconds_bucket"

	for _, metric := range raw.FindByName(metricName) {
		var (
			server = metric.Labels.Get("server")
			zone   = metric.Labels.Get("zone")
			le     = metric.Labels.Get("le")
			value  = metric.Value
		)

		if server == "" || zone == "" || le == "" {
			continue
		}

		setRequestDuration(&mx.Summary.Request, value, le)
	}

	processRequestDuration(&mx.Summary.Request)
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
			mx.NoZoneDropped.Add(value)
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

		if !cd.perServerMatcher.MatchString(server) {
			continue
		}

		if !cd.collectedServers[server] {
			cd.addNewServerCharts(server)
			cd.collectedServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &requestResponse{}
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

		if !cd.perServerMatcher.MatchString(server) {
			continue
		}

		if !cd.collectedServers[server] {
			cd.addNewServerCharts(server)
			cd.collectedServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &requestResponse{}
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

		if !cd.perServerMatcher.MatchString(server) {
			continue
		}

		if !cd.collectedServers[server] {
			cd.addNewServerCharts(server)
			cd.collectedServers[server] = true
		}

		if mx.PerServer[server] == nil {
			mx.PerServer[server] = &requestResponse{}
		}

		setResponseByRcode(&mx.PerServer[server].Response, value, rcode)
	}
}

func setRequestTotal(mx *request, value float64, family, proto, zone string) {
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

func setRequestByType(mx *request, value float64, typ string) {
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

func setResponseByRcode(mx *response, value float64, rcode string) {
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

func setRequestDuration(mx *request, value float64, le string) {
	switch le {
	case "0.00025":
		mx.Duration.LE000025.Add(value)
	case "0.0005":
		mx.Duration.LE00005.Add(value)
	case "0.001":
		mx.Duration.LE0001.Add(value)
	case "0.002":
		mx.Duration.LE0002.Add(value)
	case "0.004":
		mx.Duration.LE0004.Add(value)
	case "0.008":
		mx.Duration.LE0008.Add(value)
	case "0.016":
		mx.Duration.LE0016.Add(value)
	case "0.032":
		mx.Duration.LE0032.Add(value)
	case "0.064":
		mx.Duration.LE0064.Add(value)
	case "0.128":
		mx.Duration.LE0128.Add(value)
	case "0.256":
		mx.Duration.LE0256.Add(value)
	case "0.512":
		mx.Duration.LE0512.Add(value)
	case "1.024":
		mx.Duration.LE1024.Add(value)
	case "2.048":
		mx.Duration.LE2048.Add(value)
	case "4.096":
		mx.Duration.LE4096.Add(value)
	case "8.192":
		mx.Duration.LE8192.Add(value)
	case "+Inf":
		mx.Duration.LEInf.Add(value)
	}
}

func processRequestDuration(mx *request) {
	mx.Duration.LEInf.Sub(mx.Duration.LE8192.Value())
	mx.Duration.LE8192.Sub(mx.Duration.LE4096.Value())
	mx.Duration.LE4096.Sub(mx.Duration.LE2048.Value())
	mx.Duration.LE2048.Sub(mx.Duration.LE1024.Value())
	mx.Duration.LE1024.Sub(mx.Duration.LE0512.Value())
	mx.Duration.LE0512.Sub(mx.Duration.LE0256.Value())
	mx.Duration.LE0256.Sub(mx.Duration.LE0128.Value())
	mx.Duration.LE0128.Sub(mx.Duration.LE0064.Value())
	mx.Duration.LE0064.Sub(mx.Duration.LE0032.Value())
	mx.Duration.LE0032.Sub(mx.Duration.LE0016.Value())
	mx.Duration.LE0016.Sub(mx.Duration.LE0008.Value())
	mx.Duration.LE0008.Sub(mx.Duration.LE0004.Value())
	mx.Duration.LE0004.Sub(mx.Duration.LE0002.Value())
	mx.Duration.LE0002.Sub(mx.Duration.LE0001.Value())
	mx.Duration.LE0001.Sub(mx.Duration.LE00005.Value())
	mx.Duration.LE00005.Sub(mx.Duration.LE000025.Value())
}

func (cd *CoreDNS) addNewServerCharts(name string) {
	charts := serverCharts.Copy()
	for _, chart := range *charts {
		chart.ID = fmt.Sprintf(chart.ID, name)
		chart.Title = fmt.Sprintf(chart.Title, name)
		chart.Fam = fmt.Sprintf(chart.Fam, name)

		for _, dim := range chart.Dims {
			dim.ID = fmt.Sprintf(dim.ID, name)
		}
	}
	_ = cd.charts.Add(*charts...)
}
