package unbound

import (
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

var (
	queriesChart = Chart{
		ID:    "queries",
		Title: "Total Queries",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries",
		Dims: Dims{
			{ID: "total.num.queries", Name: "queries"},
		},
	}
	queriesIPRLChart = Chart{
		ID:    "queries_ip_ratelimited",
		Title: "Queries IP Rate Limited",
		Units: "queries",
		Fam:   "queries",
		Ctx:   "unbound.queries_ip_ratelimited",
		Dims: Dims{
			{ID: "total.num.queries_ip_ratelimited", Name: "queries"},
		},
	}
	cacheChart = Chart{
		ID:    "cache",
		Title: "Cache",
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
		Title: "Cache Percentage",
		Units: "percentage",
		Fam:   "cache",
		Ctx:   "unbound.cache_percantage",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "total.num.cachehits", Name: "hits", Algo: module.PercentOfAbsolute},
			{ID: "total.num.cachemiss", Name: "miss", Algo: module.PercentOfAbsolute},
		},
	}
	prefetchChart = Chart{
		ID:    "queries_prefetch",
		Title: "Prefetch Queries",
		Units: "queries",
		Fam:   "cache",
		Ctx:   "unbound.queries_prefetch",
		Dims: Dims{
			{ID: "total.num.prefetch", Name: "queries"},
		},
	}
	zeroTTLChart = Chart{
		ID:    "zero_ttl_responses",
		Title: "Answers Served From Expired Cache",
		Units: "responses",
		Fam:   "cache",
		Ctx:   "unbound.zero_ttl_responses",
		Dims: Dims{
			{ID: "total.num.zero_ttl", Name: "responses"},
		},
	}
	dnsCryptChart = Chart{
		ID:    "dnscrypt_queries",
		Title: "DNSCrypt Queries",
		Units: "queries",
		Fam:   "dnscrypt",
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
		Title: "The number of replies sent to queries that needed recursive processing",
		Units: "responses",
		Fam:   "responses",
		Ctx:   "unbound.recursive_replies",
		Dims: Dims{
			{ID: "total.num.recursivereplies", Name: "recursive"},
		},
	}
	recurTimeChart = Chart{
		ID:    "recursion_time",
		Title: "Time t took to answer queries that needed recursive processing",
		Units: "milliseconds",
		Fam:   "responses",
		Ctx:   "unbound.recursion_time",
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
			{ID: "total.requestlist.avg", Name: "avg"},
			//{ID: "total.requestlist.max", Name: "max"}, //
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
		Units: "events",
		Fam:   "request list",
		Ctx:   "unbound.request_list_jostle_list",
		Dims: Dims{
			{ID: "total.requestlist.overwritten", Name: "overwritten"},
			{ID: "total.requestlist.exceeded", Name: "dropped"},
		},
	}
	tcpUsageChart = Chart{
		ID:    "tcpusage",
		Title: "TCP Accept List Usage",
		Units: "events",
		Fam:   "tcpusage",
		Ctx:   "unbound.tcpusage",
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
		Dims: Dims{
			{ID: "mem_cache_message", Name: "message"},
			{ID: "mem_cache_rrset", Name: "rrset"},
			{ID: "mem_cache_dnscrypt_nonce", Name: "dnscrypt_nonce"},
			{ID: "mem_cache_dnscrypt_shared_secret", Name: "dnscrypt_shared_secret"},
		},
	}
	memModChart = Chart{
		ID:    "mod_memory",
		Title: "Module Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.mod_memory",
		Dims: Dims{
			{ID: "mem_mod_ipsecmod", Name: "ipsec"},
			{ID: "mem_mod_iterator", Name: "iterator"},
			{ID: "mem_mod_respip", Name: "respip"},
			{ID: "mem_mod_subnet", Name: "subnet"},
			{ID: "mem_mod_validator", Name: "validator"},
		},
	}
	memStreamWaitChart = Chart{
		ID:    "mem_stream_wait",
		Title: "TCP and TLS Stream Waif Buffer Memory",
		Units: "KB",
		Fam:   "mem",
		Ctx:   "unbound.mem_stream_wait",
		Dims: Dims{
			{ID: "mem_stream_wait", Name: "stream_wait"},
		},
	}
	cacheCountChart = Chart{
		ID:    "cache_count",
		Title: "Cache Count",
		Units: "items",
		Fam:   "cache count",
		Ctx:   "unbound.cache_count",
		Dims: Dims{
			{ID: "cache_count_infra", Name: "infra"},
			{ID: "cache_count_key", Name: "key"},
			{ID: "cache_count_msg", Name: "msg"},
			{ID: "cache_count_rrset", Name: "rrset"},
			{ID: "cache_count_dnscrypt_nonce", Name: "dnscrypt_nonce"},
			{ID: "cache_count_dnscrypt_shared_secret", Name: "shared_secret"},
		},
	}
	queryTypeChart = Chart{
		ID:    "query_type",
		Title: "Queries By Type",
		Units: "queries",
		Fam:   "query type",
		Ctx:   "unbound.type_queries",
	}
	queryClassChart = Chart{
		ID:    "query_class",
		Title: "Queries By Class",
		Units: "queries",
		Fam:   "query class",
		Ctx:   "unbound.class_queries",
	}
	queryOpCodeChart = Chart{
		ID:    "query_opcode",
		Title: "Queries By OpCode",
		Units: "queries",
		Fam:   "query opcode",
		Ctx:   "unbound.opcode_queries",
	}
	queryFlagChart = Chart{
		ID:    "query_flag",
		Title: "Queries By Flag",
		Units: "queries",
		Fam:   "query flag",
		Ctx:   "unbound.flag_queries",
	}
	answerRCodeChart = Chart{
		ID:    "answer_rcode",
		Title: "Answers By Rcode",
		Units: "answers",
		Fam:   "answer rcode",
		Ctx:   "unbound.rcode_answers",
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
	if hasExtendedData := len(u.curCache.queryFlags) > 0; hasExtendedData {
		return
	}

	if !u.hasExtCharts {
		charts := extendedCharts(u.Cumulative)
		if err := u.Charts().Add(*charts...); err != nil {
			u.Warning(err)
		}
		u.hasExtCharts = true
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

func newThreadCharts(id string) Charts {
	return nil
}

func (u *Unbound) addThreadCharts(id string) {
	//charts := newThreadCharts(id)`
	//if err := u.Charts().Add(charts...); err != nil {
	//	return
	//}
}

func (u *Unbound) addDimToQueryTypeChart(typ string) {
	u.addDimToChart(queryTypeChart.ID, "num.query.type"+typ, typ)
}
func (u *Unbound) addDimToQueryClassChart(class string) {
	u.addDimToChart(queryClassChart.ID, "num.query.class."+class, class)
}
func (u *Unbound) addDimToQueryOpCodeChart(opcode string) {
	u.addDimToChart(queryOpCodeChart.ID, "num.query.opcode."+opcode, opcode)
}
func (u *Unbound) addDimToQueryFlagsChart(flag string) {
	u.addDimToChart(queryFlagChart.ID, "num.query.flag."+flag, flag)
}
func (u *Unbound) addDimToAnswerRcodeChart(rcode string) {
	u.addDimToChart(answerRCodeChart.ID, "num.answer.rcode."+rcode, rcode)
}

func (u *Unbound) addDimToChart(chartID, dimID, dimName string) {
	chart := u.Charts().Get(chartID)
	if chart == nil {
		return
	}
	dim := &Dim{ID: dimID, Name: dimName}
	if u.Cumulative {
		dim.Algo = module.Incremental
	}
	if err := chart.AddDim(dim); err != nil {
		return
	}
	chart.MarkNotCreated()
}

func makeIncrIf(chart *Chart, do bool) *Chart {
	if !do {
		return chart
	}
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
