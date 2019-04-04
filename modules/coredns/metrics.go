package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

func newMetrics() *metrics {
	mx := &metrics{}
	mx.PerServer = make(map[string]*requestResponse)

	return mx
}

type metrics struct {
	Panic         mtx.Gauge                   `stm:"panic_total"`
	NoZoneDropped mtx.Gauge                   `stm:"no_matching_zone_dropped_total"`
	Summary       requestResponse             `stm:""`
	PerServer     map[string]*requestResponse `stm:""`
}

type requestResponse struct {
	Request  request  `stm:"request"`
	Response response `stm:"response"`
}

type request struct {
	Total    mtx.Gauge `stm:"total"`
	ByStatus struct {
		Processed mtx.Gauge `stm:"processed"`
		Dropped   mtx.Gauge `stm:"dropped"`
	} `stm:"per_status"`
	ByProto struct {
		UDP mtx.Gauge `stm:"udp"`
		TCP mtx.Gauge `stm:"tcp"`
	} `stm:"per_proto"`
	ByIPFamily struct {
		IPv4 mtx.Gauge `stm:"v4"`
		IPv6 mtx.Gauge `stm:"v6"`
	} `stm:"per_ip_family"`
	// https://github.com/coredns/coredns/blob/master/plugin/metrics/vars/report.go
	ByType struct {
		A      mtx.Gauge `stm:"A"`
		AAAA   mtx.Gauge `stm:"AAAA"`
		MX     mtx.Gauge `stm:"MX"`
		SOA    mtx.Gauge `stm:"SOA"`
		CNAME  mtx.Gauge `stm:"CNAME"`
		PTR    mtx.Gauge `stm:"PTR"`
		TXT    mtx.Gauge `stm:"TXT"`
		NS     mtx.Gauge `stm:"NS"`
		SRV    mtx.Gauge `stm:"SRV"`
		DS     mtx.Gauge `stm:"DS"`
		DNSKEY mtx.Gauge `stm:"DNSKEY"`
		RRSIG  mtx.Gauge `stm:"RRSIG"`
		NSEC   mtx.Gauge `stm:"NSEC"`
		NSEC3  mtx.Gauge `stm:"NSEC3"`
		IXFR   mtx.Gauge `stm:"IXFR"`
		ANY    mtx.Gauge `stm:"ANY"`
		Other  mtx.Gauge `stm:"other"`
	} `stm:"per_type"`
}

// https://github.com/miekg/dns/blob/master/types.go
// https://github.com/miekg/dns/blob/master/msg.go#L169
type response struct {
	Total   mtx.Gauge `stm:"total"`
	ByRcode struct {
		NOERROR   mtx.Gauge `stm:"NOERROR"`
		FORMERR   mtx.Gauge `stm:"FORMERR"`
		SERVFAIL  mtx.Gauge `stm:"SERVFAIL"`
		NXDOMAIN  mtx.Gauge `stm:"NXDOMAIN"`
		NOTIMP    mtx.Gauge `stm:"NOTIMP"`
		REFUSED   mtx.Gauge `stm:"REFUSED"`
		YXDOMAIN  mtx.Gauge `stm:"YXDOMAIN"`
		YXRRSET   mtx.Gauge `stm:"YXRRSET"`
		NXRRSET   mtx.Gauge `stm:"NXRRSET"`
		NOTAUTH   mtx.Gauge `stm:"NOTAUTH"`
		NOTZONE   mtx.Gauge `stm:"NOTZONE"`
		BADSIG    mtx.Gauge `stm:"BADSIG"`
		BADKEY    mtx.Gauge `stm:"BADKEY"`
		BADTIME   mtx.Gauge `stm:"BADTIME"`
		BADMODE   mtx.Gauge `stm:"BADMODE"`
		BADNAME   mtx.Gauge `stm:"BADNAME"`
		BADALG    mtx.Gauge `stm:"BADALG"`
		BADTRUNC  mtx.Gauge `stm:"BADTRUNC"`
		BADCOOKIE mtx.Gauge `stm:"BADCOOKIE"`
		Other     mtx.Gauge `stm:"other"`
	} `stm:"per_rcode"`
}
