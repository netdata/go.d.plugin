package unbound

import (
	"fmt"
	"strings"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Charts is an alias for module.Charts
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var threadPriority = orchestrator.DefaultJobPriority + len(*charts(true)) + len(*extendedCharts(true)) + 10

func charts(cumulative bool) *Charts {
	return &Charts{
		makeIncrIf(queriesChart.Copy(), cumulative),
		makeIncrIf(ipRateLimitedQueriesChart.Copy(), cumulative),
		makeIncrIf(cacheChart.Copy(), cumulative),
		makePercOfIncrIf(cachePercentageChart.Copy(), cumulative),
		makeIncrIf(prefetchChart.Copy(), cumulative),
		makeIncrIf(zeroTTLChart.Copy(), cumulative),
		makeIncrIf(dnsCryptChart.Copy(), cumulative),
		makeIncrIf(recurRepliesChart.Copy(), cumulative),
		recurTimeChart.Copy(),
		reqListUtilChart.Copy(),
		reqListCurUtilChart.Copy(),
		makeIncrIf(reqListJostleChart.Copy(), cumulative),
		tcpUsageChart.Copy(),
		uptimeChart.Copy(),
	}
}

func extendedCharts(cumulative bool) *Charts {
	return &Charts{
		memCacheChart.Copy(),
		memModChart.Copy(),
		memStreamWaitChart.Copy(),
		cacheCountChart.Copy(),
		makeIncrIf(queryTypeChart.Copy(), cumulative),
		makeIncrIf(queryClassChart.Copy(), cumulative),
		makeIncrIf(queryOpCodeChart.Copy(), cumulative),
		makeIncrIf(queryFlagChart.Copy(), cumulative),
		makeIncrIf(answerRCodeChart.Copy(), cumulative),
	}
}

func threadCharts(thread string, cumulative bool) *Charts {
	charts := charts(cumulative)
	_ = charts.Remove(uptimeChart.ID)

	for i, chart := range *charts {
		convertTotalChartToThread(chart, thread, threadPriority+i)
	}
	return charts
}

func convertTotalChartToThread(chart *Chart, thread string, priority int) {
	chart.ID = fmt.Sprintf("%s_%s", thread, chart.ID)
	chart.Title = fmt.Sprintf("%s %s", strings.Title(thread), thread)
	chart.Fam = thread + "_stats"
	chart.Ctx = fmt.Sprintf("%s_%s", chart.Ctx, thread)
	chart.Priority = priority
	for _, dim := range chart.Dims {
		dim.ID = strings.Replace(dim.ID, "total", thread, 1)
	}
}

// NOTE: chart id  should start with 'total_', name with 'Total ', ctx should ends in `_total` (convertTotalChartToThread)
var (
	queriesChart = Chart{
		ID:    "queries",
		Title: "Received Queries",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries",
		Dims: Dims{
			{ID: "total.num.queries", Name: "queries"},
		},
	}
	ipRateLimitedQueriesChart = Chart{
		ID:    "queries_ip_ratelimited",
		Title: "Rate Limited Queries",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries_ip_ratelimited",
		Dims: Dims{
			{ID: "total.num.queries_ip_ratelimited", Name: "ratelimited"},
		},
	}
	cacheChart = Chart{
		ID:    "cache",
		Title: "Cache Statistics",
		Units: "events",
		Fam:   "cache",
		Ctx:   "unbound.cache",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "total.num.cachehits", Name: "hits"},
			{ID: "total.num.cachemiss", Name: "miss"},
		},
	}
	cachePercentageChart = Chart{
		ID:    "cache_percentage",
		Title: "Cache Statistics Percentage",
		Units: "percentage",
		Fam:   "cache",
		Ctx:   "unbound.cache_percentage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "total.num.cachehits", Name: "hits", Algo: module.PercentOfAbsolute},
			{ID: "total.num.cachemiss", Name: "miss", Algo: module.PercentOfAbsolute},
		},
	}
	prefetchChart = Chart{
		ID:    "cache_prefetch",
		Title: "Cache Prefetches",
		Units: "prefetches",
		Fam:   "cache",
		Ctx:   "unbound.prefetch",
		Dims: Dims{
			{ID: "total.num.prefetch", Name: "prefetches"},
		},
	}
	zeroTTLChart = Chart{
		ID:    "zero_ttl_replies",
		Title: "Replies Served From Expired Cache",
		Units: "replies",
		Fam:   "cache",
		Ctx:   "unbound.zero_ttl_replies",
		Dims: Dims{
			{ID: "total.num.zero_ttl", Name: "zero_ttl"},
		},
	}
	// ifdef USE_DNSCRYPT
	dnsCryptChart = Chart{
		ID:    "dnscrypt_queries",
		Title: "DNSCrypt Queries",
		Units: "queries",
		Fam:   "dnscrypt queries",
		Ctx:   "unbound.dnscrypt_queries",
		Dims: Dims{
			{ID: "total.num.dnscrypt.crypted", Name: "crypted"},
			{ID: "total.num.dnscrypt.cert", Name: "cert"},
			{ID: "total.num.dnscrypt.cleartext", Name: "cleartext"},
			{ID: "total.num.dnscrypt.malformed", Name: "malformed"},
		},
	}
	recurRepliesChart = Chart{
		ID:    "recursive_replies",
		Title: "Replies That Needed Recursive Processing",
		Units: "replies",
		Fam:   "recursion",
		Ctx:   "unbound.recursive_replies",
		Dims: Dims{
			{ID: "total.num.recursivereplies", Name: "recursive"},
		},
	}
	recurTimeChart = Chart{
		ID:    "recursion_time",
		Title: "Time Spent On Recursive Processing",
		Units: "milliseconds",
		Fam:   "recursion",
		Ctx:   "unbound.recursion_time_total",
		Dims: Dims{
			{ID: "total.recursion.time.avg", Name: "avg"},
			{ID: "total.recursion.time.median", Name: "median"},
		},
	}
	reqListUtilChart = Chart{
		ID:    "request_list_utilization",
		Title: "Request List Utilization",
		Units: "queries",
		Fam:   "request list",
		Ctx:   "unbound.request_list_utilization",
		Dims: Dims{
			{ID: "total.requestlist.avg", Name: "avg", Div: 1000},
			{ID: "total.requestlist.max", Name: "max"}, // all time max in cumulative mode, never resets
		},
	}
	reqListCurUtilChart = Chart{
		ID:    "current_request_list_utilization",
		Title: "Current Request List Utilization",
		Units: "queries",
		Fam:   "request list",
		Ctx:   "unbound.current_request_list_utilization",
		Type:  module.Area,
		Dims: Dims{
			{ID: "total.requestlist.current.all", Name: "all"},
			{ID: "total.requestlist.current.user", Name: "user"},
		},
	}
	reqListJostleChart = Chart{
		ID:    "request_list_jostle_list",
		Title: "Request List Jostle List Events",
		Units: "queries",
		Fam:   "request list",
		Ctx:   "unbound.request_list_jostle_list",
		Dims: Dims{
			{ID: "total.requestlist.overwritten", Name: "overwritten"},
			{ID: "total.requestlist.exceeded", Name: "dropped"},
		},
	}
	tcpUsageChart = Chart{
		ID:    "tcpusage",
		Title: "TCP Handler Buffers",
		Units: "buffers",
		Fam:   "tcp buffers",
		Ctx:   "unbound.tcpusage",
		Dims: Dims{
			{ID: "total.tcpusage", Name: "usage"},
		},
	}
	uptimeChart = Chart{
		ID:    "uptime",
		Title: "Uptime",
		Units: "seconds",
		Fam:   "uptime",
		Ctx:   "unbound.uptime",
		Dims: Dims{
			{ID: "time.up", Name: "time"},
		},
	}
)

var (
	// TODO: do not add dnscrypt stuff by default?
	memCacheChart = Chart{
		ID:    "cache_memory",
		Title: "Cache Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.cache_memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "mem.cache.message", Name: "message", Div: 1024},
			{ID: "mem.cache.rrset", Name: "rrset", Div: 1024},
			{ID: "mem.cache.dnscrypt_nonce", Name: "dnscrypt_nonce", Div: 1024},                 // ifdef USE_DNSCRYPT
			{ID: "mem.cache.dnscrypt_shared_secret", Name: "dnscrypt_shared_secret", Div: 1024}, // ifdef USE_DNSCRYPT
		},
	}
	// TODO: do not add subnet and ipsecmod stuff by default?
	memModChart = Chart{
		ID:    "mod_memory",
		Title: "Module Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.mod_memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "mem.mod.iterator", Name: "iterator", Div: 1024},
			{ID: "mem.mod.respip", Name: "respip", Div: 1024},
			{ID: "mem.mod.validator", Name: "validator", Div: 1024},
			{ID: "mem.mod.subnet", Name: "subnet", Div: 1024},  // ifdef CLIENT_SUBNET
			{ID: "mem.mod.ipsecmod", Name: "ipsec", Div: 1024}, // ifdef USE_IPSECMOD
		},
	}
	memStreamWaitChart = Chart{
		ID:    "mem_stream_wait",
		Title: "TCP and TLS Stream Waif Buffer Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.mem_streamwait",
		Dims: Dims{
			{ID: "mem.streamwait", Name: "streamwait", Div: 1024},
		},
	}
	// NOTE: same family as for cacheChart
	cacheCountChart = Chart{
		ID:    "cache_count",
		Title: "Cache Items Count",
		Units: "items",
		Fam:   "cache",
		Ctx:   "unbound.cache_count",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "infra.cache.count", Name: "infra"},
			{ID: "key.cache.count", Name: "key"},
			{ID: "msg.cache.count", Name: "msg"},
			{ID: "rrset.cache.count", Name: "rrset"},
			{ID: "dnscrypt_nonce.cache.count", Name: "dnscrypt_nonce"},
			{ID: "dnscrypt_shared_secret.cache.count", Name: "shared_secret"},
		},
	}
	queryTypeChart = Chart{
		ID:    "queries_by_type",
		Title: "Queries By Type",
		Units: "queries",
		Fam:   "queries by type",
		Ctx:   "unbound.type_queries",
		Type:  module.Stacked,
	}
	queryClassChart = Chart{
		ID:    "queries_by_class",
		Title: "Queries By Class",
		Units: "queries",
		Fam:   "queries by class",
		Ctx:   "unbound.class_queries",
		Type:  module.Stacked,
	}
	queryOpCodeChart = Chart{
		ID:    "queries_by_opcode",
		Title: "Queries By OpCode",
		Units: "queries",
		Fam:   "queries by opcode",
		Ctx:   "unbound.opcode_queries",
		Type:  module.Stacked,
	}
	queryFlagChart = Chart{
		ID:    "queries_by_flag",
		Title: "Queries By Flag",
		Units: "queries",
		Fam:   "queries by flag",
		Ctx:   "unbound.flag_queries",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "num.query.flags.QR", Name: "QR"},
			{ID: "num.query.flags.AA", Name: "AA"},
			{ID: "num.query.flags.TC", Name: "TC"},
			{ID: "num.query.flags.RD", Name: "RD"},
			{ID: "num.query.flags.RA", Name: "RA"},
			{ID: "num.query.flags.Z", Name: "Z"},
			{ID: "num.query.flags.AD", Name: "AD"},
			{ID: "num.query.flags.CD", Name: "CD"},
		},
	}
	answerRCodeChart = Chart{
		ID:    "replies_by_rcode",
		Title: "Replies By Rcode",
		Units: "replies",
		Fam:   "replies by rcode",
		Ctx:   "unbound.rcode_answers",
		Type:  module.Stacked,
	}
)

func (u *Unbound) updateCharts() {
	if len(u.curCache.threads) > 1 {
		for v := range u.curCache.threads {
			if !u.cache.threads[v] {
				u.cache.threads[v] = true
				u.addThreadCharts(v)
			}
		}
	}
	// 0-6 rcodes always included
	if hasExtendedData := len(u.curCache.answerRCode) > 0; !hasExtendedData {
		return
	}

	if !u.extChartsCreated {
		charts := extendedCharts(u.Cumulative)
		if err := u.Charts().Add(*charts...); err != nil {
			u.Warningf("add extended charts: %v", err)
		}
		u.extChartsCreated = true
	}

	for v := range u.curCache.queryType {
		if !u.cache.queryType[v] {
			u.cache.queryType[v] = true
			u.addDimToQueryTypeChart(v)
		}
	}
	for v := range u.curCache.queryClass {
		if !u.cache.queryClass[v] {
			u.cache.queryClass[v] = true
			u.addDimToQueryClassChart(v)
		}
	}
	for v := range u.curCache.queryOpCode {
		if !u.cache.queryOpCode[v] {
			u.cache.queryOpCode[v] = true
			u.addDimToQueryOpCodeChart(v)
		}
	}
	for v := range u.curCache.answerRCode {
		if !u.cache.answerRCode[v] {
			u.cache.answerRCode[v] = true
			u.addDimToAnswerRcodeChart(v)
		}
	}
}

func (u *Unbound) addThreadCharts(thread string) {
	charts := threadCharts(thread, u.Cumulative)
	if err := u.Charts().Add(*charts...); err != nil {
		u.Warningf("add '%s' charts: %v", thread, err)
	}
}

func (u *Unbound) addDimToQueryTypeChart(typ string) {
	u.addDimToChart(queryTypeChart.ID, "num.query.type."+typ, typ)
}
func (u *Unbound) addDimToQueryClassChart(class string) {
	u.addDimToChart(queryClassChart.ID, "num.query.class."+class, class)
}
func (u *Unbound) addDimToQueryOpCodeChart(opcode string) {
	u.addDimToChart(queryOpCodeChart.ID, "num.query.opcode."+opcode, opcode)
}
func (u *Unbound) addDimToAnswerRcodeChart(rcode string) {
	u.addDimToChart(answerRCodeChart.ID, "num.answer.rcode."+rcode, rcode)
}

func (u *Unbound) addDimToChart(chartID, dimID, dimName string) {
	chart := u.Charts().Get(chartID)
	if chart == nil {
		u.Warningf("add '%s' dim: couldn't find '%s' chart", dimID, chartID)
		return
	}
	dim := &Dim{ID: dimID, Name: dimName}
	if u.Cumulative {
		dim.Algo = module.Incremental
	}
	if err := chart.AddDim(dim); err != nil {
		u.Warningf("add '%s' dim: %v", dimID, err)
		return
	}
	chart.MarkNotCreated()
}

func makeIncrIf(chart *Chart, do bool) *Chart {
	if !do {
		return chart
	}
	chart.Units += "/s"
	for _, d := range chart.Dims {
		d.Algo = module.Incremental
	}
	return chart
}

func makePercOfIncrIf(chart *Chart, do bool) *Chart {
	if !do {
		return chart
	}
	for _, d := range chart.Dims {
		d.Algo = module.PercentOfIncremental
	}
	return chart
}
