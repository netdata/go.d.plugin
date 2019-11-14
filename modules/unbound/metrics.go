package unbound

import (
	"github.com/netdata/go.d.plugin/pkg/metrics"
)

/*
total.num.queries=0
total.num.queries_ip_ratelimited=0
total.num.cachehits=0
total.num.cachemiss=0
total.num.prefetch=0
total.num.zero_ttl=0
total.num.recursivereplies=0
total.num.dnscrypt.crypted=0
total.num.dnscrypt.cert=0
total.num.dnscrypt.cleartext=0
total.num.dnscrypt.malformed=0
total.requestlist.avg=0
total.requestlist.max=0
total.requestlist.overwritten=0
total.requestlist.exceeded=0
total.requestlist.current.all=0
total.requestlist.current.user=0
total.recursion.time.avg=0.000000
total.recursion.time.median=0
total.tcpusage=0
*/

type metricsData struct {
	Total struct {
		Queries     metrics.Gauge `stm:"queries"`
		QueriesIPRL metrics.Gauge `stm:"queries_ip_ratelimited"`
		Cache       struct {
			Hits metrics.Gauge `stm:"hits"`
			Miss metrics.Gauge `stm:"miss"`
		} `stm:"cache"`
		Prefetch         metrics.Gauge `stm:"prefetch"`
		ZeroTTL          metrics.Gauge `stm:"zero_ttl"`
		RecursiveReplies metrics.Gauge `stm:"recursive_replies"`
		DNSCrypt         struct {
			Crypted   metrics.Gauge `stm:"crypted"`
			Cert      metrics.Gauge `stm:"cert"`
			ClearText metrics.Gauge `stm:"clear_text"`
			Malformed metrics.Gauge `stm:"_malformed"`
		} `stm:"dns_crypt"`
		RequestList struct {
			Avg         metrics.Gauge `stm:"avg"`
			Max         metrics.Gauge `stm:"max"`
			Overwritten metrics.Gauge `stm:"overwritten"`
			Exceeded    metrics.Gauge `stm:"exceeded"`
			CurrentAll  metrics.Gauge `stm:"current_all"`
			CurrentUser metrics.Gauge `stm:"current_user"`
		} `stm:"request_list"`
		RecursionTime struct {
			Avg    metrics.Gauge `stm:"avg"`
			Median metrics.Gauge `stm:"median"`
		} `stm:"recursion_time"`
		TCPUsage metrics.Gauge `stm:"tcp_usage"`
	}
	Uptime metrics.Gauge `stm:"uptime"`
}
