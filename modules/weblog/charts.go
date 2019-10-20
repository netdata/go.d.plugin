package weblog

import (
	"fmt"

	"github.com/netdata/go-orchestrator"
	"github.com/netdata/go-orchestrator/module"
)

type (
	Charts = module.Charts
	Chart  = module.Chart
	Dims   = module.Dims
	Dim    = module.Dim
)

const (
	// TODO: chart priorities
	defaultPriority = orchestrator.DefaultJobPriority
)

// NOTE: inconsistency between contexts with python web_log
// TODO: current histogram charts are misleading in netdata

var charts = Charts{
	{
		ID:    "requests",
		Title: "Requests",
		Units: "requests/s",
		Fam:   "requests",
		Ctx:   "web_log.requests",
		Type:  module.Area,
		Dims: Dims{
			{ID: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "response_statuses",
		Title: "Response Statuses",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_statuses",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "resp_successful", Name: "success", Algo: module.Incremental},
			{ID: "resp_client_error", Name: "bad", Algo: module.Incremental},
			{ID: "resp_redirect", Name: "redirect", Algo: module.Incremental},
			{ID: "resp_server_error", Name: "error", Algo: module.Incremental},
		},
	},
	{
		ID:    "response_codes",
		Title: "Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
		},
	},
}

var (
	bandwidth = Chart{
		ID:    "bandwidth",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "web_log.bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_received", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}

	requestsPerHTTPMethod = Chart{
		ID:    "requests_per_http_method",
		Title: "Requests Per HTTP ReqHTTPMethod",
		Units: "requests/s",
		Fam:   "http methods",
		Ctx:   "web_log.requests_per_http_method",
		Type:  module.Stacked,
		Dims:  Dims{},
	}

	requestsPerHTTPVersion = Chart{
		ID:    "requests_per_http_version",
		Title: "Requests Per HTTP ReqHTTPVersion",
		Units: "requests/s",
		Fam:   "http versions",
		Ctx:   "web_log.requests_per_http_version",
		Type:  module.Stacked,
	}

	requestsPerIPProto = Chart{
		ID:    "requests_per_ip_proto",
		Title: "Requests Per IP Protocol",
		Units: "requests/s",
		Fam:   "ip protocols",
		Ctx:   "web_log.requests_per_ip_proto",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}

	uniqueReqPerIPCurPoll = Chart{
		ID:    "clients_current",
		Title: "Current Poll Unique ClientAddr IPs",
		Units: "unique ips",
		Fam:   "clients",
		Ctx:   "web_log.current_poll_ips",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "unique_current_poll_ipv4", Name: "ipv4", Algo: module.Absolute},
			{ID: "unique_current_poll_ipv6", Name: "ipv6", Algo: module.Absolute},
		},
	}

	responseCodesDetailed = Chart{
		ID:    "detailed_response_codes",
		Title: "Detailed Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes_detailed",
		Type:  module.Stacked,
	}

	responseTime = Chart{
		ID:    "response_time",
		Title: "Processing Time",
		Units: "milliseconds",
		Fam:   "timings",
		Ctx:   "web_log.response_time",
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Div: 1000},
			{ID: "resp_time_max", Name: "max", Div: 1000},
			{ID: "resp_time_avg", Name: "avg", Div: 1000},
		},
	}

	responseTimeHistogram = Chart{
		ID:    "response_time_histogram",
		Title: "Processing Time Histogram",
		Units: "requests/s",
		Fam:   "timings",
		Ctx:   "web_log.response_time_histogram",
	}

	responseTimeUpstream = Chart{
		ID:    "response_time_upstream",
		Title: "Processing Time Upstream",
		Units: "milliseconds",
		Fam:   "timings",
		Ctx:   "web_log.response_time_upstream",
		Dims: Dims{
			{ID: "resp_time_upstream_min", Name: "min", Div: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Div: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Div: 1000},
		},
	}

	responseTimeUpstreamHistogram = Chart{
		ID:    "response_time_upstream_histogram",
		Title: "Processing Time Upstream Histogram",
		Units: "requests/s",
		Fam:   "timings",
		Ctx:   "web_log.response_time_upstream_histogram",
	}

	requestsPerURL = Chart{
		ID:    "requests_per_url",
		Title: "Requests Per Url",
		Units: "requests/s",
		Fam:   "urls",
		Ctx:   "web_log.requests_per_url",
		Type:  module.Stacked,
	}

	requestsPerUserDefined = Chart{
		ID:    "requests_per_user_defined",
		Title: "Requests Per User Defined Pattern",
		Units: "requests/s",
		Fam:   "user defined",
		Ctx:   "web_log.requests_per_user_defined",
		Type:  module.Stacked,
	}

	requestsPerVhost = Chart{
		ID:    "requests_per_vhost",
		Title: "Requests Per Vhost",
		Units: "requests/s",
		Fam:   "vhost",
		Ctx:   "web_log.requests_per_vhost",
		Type:  module.Stacked,
	}
)

func newResponseCodesDetailedPerFamilyCharts() []*Chart {
	return []*Chart{
		{
			ID:    responseCodesDetailed.ID + "_1xx",
			Title: "Detailed Response Codes 1xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_1xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_2xx",
			Title: "Detailed Response Codes 2xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_2xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_3xx",
			Title: "Detailed Response Codes 3xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_3xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_4xx",
			Title: "Detailed Response Codes 4xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_4xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_5xx",
			Title: "Detailed Response Codes 5xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_5xx",
			Type:  module.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_other",
			Title: "Detailed Response Codes Other",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_other",
			Type:  module.Stacked,
		},
	}
}

func newRequestsPerURLCategoriesChart(cats []*category) *Chart {
	chart := requestsPerURL.Copy()
	for _, c := range cats {
		dim := &Dim{
			ID:   "req_uri_" + c.name,
			Name: c.name,
			Algo: module.Incremental,
		}
		panicIfErr(chart.AddDim(dim))
	}
	return chart
}

func newRequestsPerCustomCategoriesChart(cats []*category) *Chart {
	chart := requestsPerUserDefined.Copy()
	for _, c := range cats {
		dim := &Dim{
			ID:   "req_custom_" + c.name,
			Name: c.name,
			Algo: module.Incremental,
		}
		panicIfErr(chart.AddDim(dim))
	}
	return chart
}

func newResponseTimeHistogramChart(histogram []float64) *Chart {
	chart := responseTimeHistogram.Copy()
	for i, v := range histogram {
		dimID := fmt.Sprintf("resp_time_hist_bucket_%d", i+1)
		name := fmt.Sprintf("%.3f", v)
		dim := &Dim{
			ID:   dimID,
			Name: name,
			Algo: module.Incremental,
		}
		panicIfErr(chart.AddDim(dim))
	}
	panicIfErr(chart.AddDim(&Dim{
		ID:   "resp_time_hist_count",
		Name: "+Inf",
		Algo: module.Incremental,
	}))
	return chart
}

func newResponseTimeUpstreamHistogramChart(histogram []float64) *Chart {
	chart := responseTimeUpstreamHistogram.Copy()
	for i, v := range histogram {
		dimID := fmt.Sprintf("resp_time_upstream_hist_bucket_%d", i+1)
		name := fmt.Sprintf("%.3f", v)
		dim := &Dim{
			ID:   dimID,
			Name: name,
			Algo: module.Incremental,
		}
		panicIfErr(chart.AddDim(dim))
	}
	panicIfErr(chart.AddDim(&Dim{
		ID:   "resp_time_upstream_hist_count",
		Name: "+Inf",
		Algo: module.Incremental,
	}))
	return chart
}

func (w *WebLog) updateCharts() {
	if w.col.vhost {
		w.addVhostChart()
		w.updateVhostChart()
	}
	if w.col.client {
		w.addClientCharts()
	}
	if w.col.method {
		w.addHTTPMethodChart()
		w.updateHTTPMethodChart()
	}
	if w.col.uri {
		w.addURIChart()
	}
	if w.col.version {
		w.addHTTPVersionChart()
		w.updateHTTPVersionChart()
	}
	if w.col.status {
		w.addRespCodesDetailedChart()
		w.updateRespCodesDetailedChart()
	}
	if w.col.reqSize || w.col.respSize {
		w.addBandwidthChart()
	}
	if w.col.respTime {
		w.addRespTimeCharts()
	}
	if w.col.upRespTime {
		w.addUpstreamRespTimeCharts()
	}
	if w.col.custom {
		w.addCustomChart()
	}
}

func (w *WebLog) addVhostChart() {
	if w.chartsCache.created.addIfNotExist(requestsPerVhost.ID) {
		return
	}
	panicIfErr(w.Charts().Add(requestsPerVhost.Copy()))
}

func (w *WebLog) addClientCharts() {
	if w.chartsCache.created.addIfNotExist(requestsPerIPProto.ID) {
		return
	}

	panicIfErr(w.Charts().Add(requestsPerIPProto.Copy()))
	panicIfErr(w.Charts().Add(uniqueReqPerIPCurPoll.Copy()))
}

func (w *WebLog) addHTTPMethodChart() {
	if w.chartsCache.created.addIfNotExist(requestsPerHTTPMethod.ID) {
		return
	}

	panicIfErr(w.Charts().Add(requestsPerHTTPMethod.Copy()))
}

func (w *WebLog) addURIChart() {
	if w.chartsCache.created.addIfNotExist(requestsPerURL.ID) {
		return
	}

	chart := newRequestsPerURLCategoriesChart(w.urlCats)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addHTTPVersionChart() {
	if w.chartsCache.created.addIfNotExist(requestsPerHTTPVersion.ID) {
		return
	}

	panicIfErr(w.Charts().Add(requestsPerHTTPVersion.Copy()))
}

func (w *WebLog) addRespCodesDetailedChart() {
	if w.chartsCache.created.addIfNotExist(responseCodesDetailed.ID) {
		return
	}

	if w.AggregateResponseCodes {
		panicIfErr(w.Charts().Add(responseCodesDetailed.Copy()))
		return
	}

	// TODO: do not create charts for all families
	panicIfErr(w.Charts().Add(newResponseCodesDetailedPerFamilyCharts()...))
}

func (w *WebLog) addBandwidthChart() {
	if w.chartsCache.created.addIfNotExist(bandwidth.ID) {
		return
	}

	panicIfErr(w.Charts().Add(bandwidth.Copy()))
}

func (w *WebLog) addRespTimeCharts() {
	if w.chartsCache.created.addIfNotExist(responseTime.ID) {
		return
	}

	panicIfErr(w.Charts().Add(responseTime.Copy()))

	if len(w.Histogram) == 0 {
		return
	}

	chart := newResponseTimeHistogramChart(w.Histogram)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addUpstreamRespTimeCharts() {
	if w.chartsCache.created.addIfNotExist(responseTimeUpstream.ID) {
		return
	}

	panicIfErr(w.Charts().Add(responseTimeUpstream.Copy()))

	if len(w.Histogram) == 0 {
		return
	}

	chart := newResponseTimeUpstreamHistogramChart(w.Histogram)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addCustomChart() {
	if w.chartsCache.created.addIfNotExist(requestsPerUserDefined.ID) {
		return
	}

	chart := newRequestsPerCustomCategoriesChart(w.userCats)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) updateVhostChart() {
	chart := w.Charts().Get(requestsPerVhost.ID)

	for v := range w.mx.ReqVhost {
		if w.chartsCache.vhosts.addIfNotExist(v) {
			continue
		}
		addDimToVhostChart(chart, v)
	}
}

func (w *WebLog) updateHTTPMethodChart() {
	chart := w.Charts().Get(requestsPerHTTPMethod.ID)

	for v := range w.mx.ReqMethod {
		if w.chartsCache.methods.addIfNotExist(v) {
			continue
		}
		addDimToHTTPMethodChart(chart, v)
	}
}

func (w *WebLog) updateHTTPVersionChart() {
	chart := w.Charts().Get(requestsPerHTTPVersion.ID)

	for v := range w.mx.ReqVersion {
		if w.chartsCache.versions.addIfNotExist(v) {
			continue
		}
		addDimToHTTPVersionChart(chart, v)
	}
}

func (w *WebLog) updateRespCodesDetailedChart() {
	var chart *Chart
	for v := range w.mx.RespCode {
		if w.chartsCache.codes.addIfNotExist(v) {
			continue
		}
		chart = w.respCodesDetailedChartByCode(v)
		addDimToRespCodesDetailedChart(chart, v)
	}
}

func (w *WebLog) respCodesDetailedChartByCode(code string) *Chart {
	if w.AggregateResponseCodes {
		return w.Charts().Get(responseCodesDetailed.ID)
	}

	var id string
	switch v := code[:1]; v {
	case "1", "2", "3", "4", "5":
		id = fmt.Sprintf("%s_%sxx", responseCodesDetailed.ID, v)
	default:
		id = fmt.Sprintf("%s_other", responseCodesDetailed.ID)
	}
	return w.Charts().Get(id)
}

func addDimToRespCodesDetailedChart(chart *Chart, code string) {
	dimID := "req_code_" + code
	dim := &Dim{
		ID:   dimID,
		Name: code,
		Algo: module.Incremental,
	}
	panicIfErr(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func addDimToHTTPVersionChart(chart *Chart, version string) {
	dimID := "req_version_" + version
	dim := &Dim{
		ID:   dimID,
		Name: version,
		Algo: module.Incremental,
	}
	panicIfErr(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func addDimToHTTPMethodChart(chart *Chart, method string) {
	dimID := "req_method_" + method
	dim := &Dim{
		ID:   dimID,
		Name: method,
		Algo: module.Incremental,
	}
	panicIfErr(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func addDimToVhostChart(chart *Chart, vhost string) {
	dimID := "req_vhost_" + vhost
	dim := &Dim{
		ID:   dimID,
		Name: vhost,
		Algo: module.Incremental,
	}
	panicIfErr(chart.AddDim(dim))
	chart.MarkNotCreated()
}

// TODO: get rid of
func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: per ReqURI category charts
//func perCategoryStats(id string) []*Chart {
//	return []*Chart{
//		{
//			ID:    responseCodesDetailed.ID + "_" + id,
//			Title: "Detailed Response Codes",
//			Units: "requests/s",
//			Fam:   "url " + id,
//			Ctx:   "web_log.response_codes_detailed_per_url",
//			Type:  module.Stacked,
//		},
//		{
//			ID:    bandwidth.ID + "_" + id,
//			Title: "Bandwidth",
//			Units: "kilobits/s",
//			Fam:   "url " + id,
//			Ctx:   "web_log.bandwidth_per_url",
//			Type:  module.Area,
//			Dims: Dims{
//				{ID: id + "_resp_length", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
//				{ID: id + "_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
//			},
//		},
//		{
//			ID:    responseTime.ID + "_" + id,
//			Title: "Processing Time",
//			Units: "milliseconds",
//			Fam:   "url " + id,
//			Ctx:   "web_log.response_time_per_url",
//			Type:  module.Area,
//			Dims: Dims{
//				{ID: id + "_resp_time_min", Name: "min", Algo: module.Incremental, Div: 1000},
//				{ID: id + "_resp_time_max", Name: "max", Algo: module.Incremental, Div: 1000},
//				{ID: id + "_resp_time_avg", Name: "avg", Algo: module.Incremental, Div: 1000},
//			},
//		},
//	}
//}

type (
	cache map[string]struct{}

	chartsCache struct {
		created  cache
		vhosts   cache
		methods  cache
		versions cache
		codes    cache
	}
)

func (c cache) has(v string) bool {
	_, ok := c[v]
	return ok
}

func (c cache) add(v string) {
	c[v] = struct{}{}
}

func (c cache) addIfNotExist(v string) (exist bool) {
	if c.has(v) {
		return true
	}
	c.add(v)
	return
}
