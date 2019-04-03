package coredns

import (
	mtx "github.com/netdata/go.d.plugin/pkg/metrics"
)

func newMetrics() *metrics {
	return &metrics{
		PerServer: make(map[string]*serverMetrics),
	}
}

type metrics struct {
	Summary struct {
		Panic   mtx.Gauge      `stm:"panic_total"`
		Request requestMetrics `stm:"request"`
	} `stm:""`
	PerServer map[string]*serverMetrics `stm:""`
}

type serverMetrics struct {
	Request requestMetrics `stm:"request"`
}

type requestMetrics struct {
	Total    mtx.Gauge `stm:"total"`
	ByStatus struct {
		Processed mtx.Gauge `stm:"processed"`
		Dropped   mtx.Gauge `stm:"dropped"`
	} `stm:"by_status"`
	ByProto struct {
		UDP mtx.Gauge `stm:"udp"`
		TCP mtx.Gauge `stm:"tcp"`
	} `stm:"by_proto"`
	ByIPFamily struct {
		IPv4 mtx.Gauge `stm:"v4"`
		IPv6 mtx.Gauge `stm:"v6"`
	} `stm:"by_ip_family"`
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
	} `stm:"by_type"`
}
