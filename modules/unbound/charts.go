package unbound

import (
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
		makeIncrIf(queriesIPRLChart.Copy(), cumulative),
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
	chart.ID = strings.Replace(chart.ID, "total", thread, 1)
	chart.Title = strings.Replace(chart.Title, "Total", strings.Title(thread), 1)
	chart.Fam = thread + "_stats"
	chart.Ctx = strings.Replace(chart.Ctx, "total", thread, 1)
	chart.Priority = priority
	for _, dim := range chart.Dims {
		dim.ID = strings.Replace(dim.ID, "total", thread, 1)
	}
}

// NOTE: chart id  should start with 'total_', name with 'Total ', ctx should ends in `_total` (convertTotalChartToThread)
var (
	queriesChart = Chart{
		ID:    "total_queries",
		Title: "Total Queries",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries_total",
		Dims: Dims{
			{ID: "total.num.queries", Name: "queries"},
		},
	}
	queriesIPRLChart = Chart{
		ID:    "total_queries_ip_ratelimited",
		Title: "Total Queries IP Rate Limited",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries_ip_ratelimited_total",
		Dims: Dims{
			{ID: "total.num.queries_ip_ratelimited", Name: "queries"},
		},
	}
	cacheChart = Chart{
		ID:    "total_cache",
		Title: "Total Cache",
		Units: "events",
		Fam:   "cache",
		Ctx:   "unbound.cache_total",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "total.num.cachehits", Name: "hits"},
			{ID: "total.num.cachemiss", Name: "miss"},
		},
	}
	cachePercentageChart = Chart{
		ID:    "total_cache_percentage",
		Title: "Total Cache Percentage",
		Units: "percentage",
		Fam:   "cache",
		Ctx:   "unbound.cache_percantage_total",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "total.num.cachehits", Name: "hits", Algo: module.PercentOfAbsolute},
			{ID: "total.num.cachemiss", Name: "miss", Algo: module.PercentOfAbsolute},
		},
	}
	prefetchChart = Chart{
		ID:    "total_queries_prefetch",
		Title: "Total Cache Prefetches",
		Units: "queries",
		Fam:   "cache",
		Ctx:   "unbound.queries_prefetch_total",
		Dims: Dims{
			{ID: "total.num.prefetch", Name: "queries"},
		},
	}
	zeroTTLChart = Chart{
		ID:    "total_zero_ttl_responses",
		Title: "Total Answers Served From Expired Cache",
		Units: "responses",
		Fam:   "cache",
		Ctx:   "unbound.zero_ttl_responses_total",
		Dims: Dims{
			{ID: "total.num.zero_ttl", Name: "responses"},
		},
	}
	dnsCryptChart = Chart{
		ID:    "total_dnscrypt_queries",
		Title: "Total DNSCrypt Queries",
		Units: "queries",
		Fam:   "dnscrypt",
		Ctx:   "unbound.dnscrypt_queries_total",
		Dims: Dims{
			{ID: "total.num.dnscrypt.crypted", Name: "crypted"},
			{ID: "total.num.dnscrypt.cert", Name: "cert"},
			{ID: "total.num.dnscrypt.cleartext", Name: "cleartext"},
			{ID: "total.num.dnscrypt.malformed", Name: "malformed"},
		},
	}
	recurRepliesChart = Chart{
		ID:    "total_recursive_replies",
		Title: "Total number of replies sent to queries that needed recursive processing",
		Units: "responses",
		Fam:   "responses",
		Ctx:   "unbound.recursive_replies_total",
		Dims: Dims{
			{ID: "total.num.recursivereplies", Name: "recursive"},
		},
	}
	recurTimeChart = Chart{
		ID:    "total_recursion_time",
		Title: "Total Time t took to answer queries that needed recursive processing",
		Units: "milliseconds",
		Fam:   "responses",
		Ctx:   "unbound.recursion_time_total",
		Dims: Dims{
			{ID: "total.recursion.time.avg", Name: "avg", Div: 1000},
			{ID: "total.recursion.time.median", Name: "median", Div: 1000},
		},
	}
	reqListUtilChart = Chart{
		ID:    "total_request_list_utilization",
		Title: "Total Request List Utilization",
		Units: "queries",
		Fam:   "request list",
		Ctx:   "unbound.request_list_utilization_total",
		Dims: Dims{
			{ID: "total.requestlist.avg", Name: "avg", Div: 1000000},
			//{ID: "total.requestlist.max", Name: "max"}, //
		},
	}
	reqListCurUtilChart = Chart{
		ID:    "total_current_request_list_utilization",
		Title: "Total Current Request List Utilization",
		Units: "queries",
		Fam:   "request list",
		Ctx:   "unbound.current_request_list_utilization_total",
		Type:  module.Area,
		Dims: Dims{
			{ID: "total.requestlist.current.all", Name: "all"},
			{ID: "total.requestlist.current.user", Name: "user"},
		},
	}
	reqListJostleChart = Chart{
		ID:    "total_request_list_jostle_list",
		Title: "Total Request List Jostle List Events",
		Units: "events",
		Fam:   "request list",
		Ctx:   "unbound.request_list_jostle_list_total",
		Dims: Dims{
			{ID: "total.requestlist.overwritten", Name: "overwritten"},
			{ID: "total.requestlist.exceeded", Name: "dropped"},
		},
	}
	tcpUsageChart = Chart{
		ID:    "total_tcpusage",
		Title: "Total TCP Accept List Usage",
		Units: "events",
		Fam:   "tcpusage",
		Ctx:   "unbound.tcpusage_total",
		Dims: Dims{
			{ID: "total.tcpusage", Name: "tcpusage"},
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
			{ID: "mem.cache.dnscrypt_nonce", Name: "dnscrypt_nonce", Div: 1024},
			{ID: "mem.cache.dnscrypt_shared_secret", Name: "dnscrypt_shared_secret", Div: 1024},
		},
	}
	memModChart = Chart{
		ID:    "mod_memory",
		Title: "Module Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.mod_memory",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "mem.mod.ipsecmod", Name: "ipsec", Div: 1024},
			{ID: "mem.mod.iterator", Name: "iterator", Div: 1024},
			{ID: "mem.mod.respip", Name: "respip", Div: 1024},
			{ID: "mem.mod.subnet", Name: "subnet", Div: 1024},
			{ID: "mem.mod.validator", Name: "validator", Div: 1024},
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
		Title: "Cache Count",
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
		ID:    "query_type",
		Title: "Queries By Type",
		Units: "queries",
		Fam:   "query type",
		Ctx:   "unbound.type_queries",
		Type:  module.Stacked,
	}
	queryClassChart = Chart{
		ID:    "query_class",
		Title: "Queries By Class",
		Units: "queries",
		Fam:   "query class",
		Ctx:   "unbound.class_queries",
		Type:  module.Stacked,
	}
	queryOpCodeChart = Chart{
		ID:    "query_opcode",
		Title: "Queries By OpCode",
		Units: "queries",
		Fam:   "query opcode",
		Ctx:   "unbound.opcode_queries",
		Type:  module.Stacked,
	}
	queryFlagChart = Chart{
		ID:    "query_flag",
		Title: "Queries By Flag",
		Units: "queries",
		Fam:   "query flag",
		Ctx:   "unbound.flag_queries",
		Type:  module.Stacked,
	}
	answerRCodeChart = Chart{
		ID:    "answer_rcode",
		Title: "Answers By Rcode",
		Units: "answers",
		Fam:   "answer rcode",
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
	if hasExtendedData := len(u.curCache.queryFlags) > 0; !hasExtendedData {
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
	for v := range u.curCache.queryFlags {
		if !u.cache.queryFlags[v] {
			u.cache.queryFlags[v] = true
			u.addDimToQueryFlagsChart(v)
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
func (u *Unbound) addDimToQueryFlagsChart(flag string) {
	u.addDimToChart(queryFlagChart.ID, "num.query.flags."+flag, flag)
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
