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
	defaultPriority           = orchestrator.DefaultJobPriority
	prioReqTotal              = defaultPriority
	prioReqUnreported         = defaultPriority + 1
	prioRespStatuses          = defaultPriority + 2
	prioRespCodesClass        = defaultPriority + 3
	prioRespCodes             = defaultPriority + 4
	prioRespCodes1xx          = defaultPriority + 5
	prioRespCodes2xx          = defaultPriority + 6
	prioRespCodes3xx          = defaultPriority + 7
	prioRespCodes4xx          = defaultPriority + 8
	prioRespCodes5xx          = defaultPriority + 9
	prioBandwidth             = defaultPriority + 10
	prioReqProcTime           = defaultPriority + 11
	prioRespTimeHist          = defaultPriority + 12
	prioUpsRespTime           = defaultPriority + 13
	prioUpsRespTimeHist       = defaultPriority + 14
	prioUniqIP                = defaultPriority + 15
	prioReqVhost              = defaultPriority + 16
	prioReqPort               = defaultPriority + 17
	prioReqScheme             = defaultPriority + 18
	prioReqMethod             = defaultPriority + 19
	prioReqVersion            = defaultPriority + 20
	prioReqIPProto            = defaultPriority + 21
	prioReqSSLProto           = defaultPriority + 22
	prioReqSSLCipherSuite     = defaultPriority + 23
	prioReqCustomFieldPattern = defaultPriority + 40 // chart per custom field, alphabetical order
	prioReqURLPattern         = defaultPriority + 41
	prioURLPatternStats       = defaultPriority + 42 // 3 charts per url pattern, alphabetical order
)

// NOTE: inconsistency between contexts with python web_log
// TODO: current histogram charts are misleading in netdata

// Total Requests       [requests]
// Unreported Requests  [requests]
// Resp Statuses        [responses]
// Resp Codes By Class  [responses]
// Resp Codes           [responses]
// Bandwidth            [bandwidth]
// Resp Time            [timings]
// Resp Time Hist       [timings]
// Resp Time Ups        [upstream]
// Resp Time Hist Ups   [upstream]
// Uniq IPs             [clients]
// Req By Vhost        [req vhost]
// Req By Port         [req port]
// Req By Scheme       [req scheme]
// Req By Method       [req method]
// Req By Version      [req version]
// Req By IP Proto     [req ip proto]
// Req By Custom       [req custom]
// Req By URL          [req url]
// URL Stats            [url <name>]

// Requests
var (
	reqTotal = Chart{
		ID:       "requests",
		Title:    "Total Requests",
		Units:    "requests/s",
		Fam:      "requests",
		Ctx:      "web_log.requests",
		Priority: prioReqTotal,
		Dims: Dims{
			{ID: "requests", Algo: module.Incremental},
		},
	}

	reqExcluded = Chart{
		ID:       "excluded_requests",
		Title:    "Excluded Requests",
		Units:    "requests/s",
		Fam:      "requests",
		Ctx:      "web_log.excluded_requests",
		Type:     module.Stacked,
		Priority: prioReqUnreported,
		Dims: Dims{
			{ID: "req_filtered", Name: "filtered", Algo: module.Incremental},
			{ID: "req_unmatched", Name: "unmatched", Algo: module.Incremental},
		},
	}
)

// Responses
var (
	// netdata specific grouping
	respStatuses = Chart{
		ID:       "type_responses",
		Title:    "Responses By Type",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.type_responses",
		Type:     module.Stacked,
		Priority: prioRespStatuses,
		Dims: Dims{
			{ID: "resp_successful", Name: "success", Algo: module.Incremental},
			{ID: "resp_client_error", Name: "bad", Algo: module.Incremental},
			{ID: "resp_redirect", Name: "redirect", Algo: module.Incremental},
			{ID: "resp_server_error", Name: "error", Algo: module.Incremental},
		},
	}
	respCodeClass = Chart{
		ID:       "status_code_class_responses",
		Title:    "Responses By Status Code Class",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_responses",
		Type:     module.Stacked,
		Priority: prioRespCodesClass,
		Dims: Dims{
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
		},
	}
	respCodes = Chart{
		ID:       "status_code_responses",
		Title:    "Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes,
	}
	respCodes1xx = Chart{
		ID:       "status_code_class_1xx_responses",
		Title:    "Informational Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_1xx_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes1xx,
	}
	respCodes2xx = Chart{
		ID:       "status_code_class_2xx_responses",
		Title:    "Successful Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_2xx_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes2xx,
	}
	respCodes3xx = Chart{
		ID:       "status_code_class_3xx_responses",
		Title:    "Redirects Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_3xx_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes3xx,
	}
	respCodes4xx = Chart{
		ID:       "status_code_class_4xx_responses",
		Title:    "Client Errors Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_4xx_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes4xx,
	}
	respCodes5xx = Chart{
		ID:       "status_code_class_5xx_responses",
		Title:    "Server Errors Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.status_code_class_5xx_responses",
		Type:     module.Stacked,
		Priority: prioRespCodes5xx,
	}
)

// Bandwidth
var (
	bandwidth = Chart{
		ID:       "bandwidth",
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "bandwidth",
		Ctx:      "web_log.bandwidth",
		Type:     module.Area,
		Priority: prioBandwidth,
		Dims: Dims{
			{ID: "bytes_received", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}
)

// Timings
var (
	reqProcTime = Chart{
		ID:       "request_processing_time",
		Title:    "Request Processing Time",
		Units:    "milliseconds",
		Fam:      "timings",
		Ctx:      "web_log.request_processing_time",
		Priority: prioReqProcTime,
		Dims: Dims{
			{ID: "req_proc_time_min", Name: "min", Div: 1000},
			{ID: "req_proc_time_max", Name: "max", Div: 1000},
			{ID: "req_proc_time_avg", Name: "avg", Div: 1000},
		},
	}
	reqProcTimeHist = Chart{
		ID:       "requests_processing_time_histogram",
		Title:    "Requests Processing Time Histogram",
		Units:    "requests/s",
		Fam:      "timings",
		Ctx:      "web_log.requests_processing_time_histogram",
		Priority: prioRespTimeHist,
	}
)

// Upstream
var (
	upsRespTime = Chart{
		ID:       "upstream_response_time",
		Title:    "Upstream Response Time",
		Units:    "milliseconds",
		Fam:      "timings",
		Ctx:      "web_log.upstream_response_time",
		Priority: prioUpsRespTime,
		Dims: Dims{
			{ID: "upstream_resp_time_min", Name: "min", Div: 1000},
			{ID: "upstream_resp_time_max", Name: "max", Div: 1000},
			{ID: "upstream_resp_time_avg", Name: "avg", Div: 1000},
		},
	}
	upsRespTimeHist = Chart{
		ID:       "upstream_responses_time_histogram",
		Title:    "Upstream Responses Time Histogram",
		Units:    "responses/s",
		Fam:      "timings",
		Ctx:      "web_log.upstream_responses_time_histogram",
		Priority: prioUpsRespTimeHist,
	}
)

// Clients
var (
	uniqIPsCurPoll = Chart{
		ID:       "current_poll_uniq_clients",
		Title:    "Current Poll Unique Clients",
		Units:    "clients",
		Fam:      "clients",
		Ctx:      "web_log.current_poll_uniq_clients",
		Type:     module.Stacked,
		Priority: prioUniqIP,
		Dims: Dims{
			{ID: "uniq_ipv4", Name: "ipv4", Algo: module.Absolute},
			{ID: "uniq_ipv6", Name: "ipv6", Algo: module.Absolute},
		},
	}
)

// Request By N
var (
	reqByVhost = Chart{
		ID:       "vhost_requests",
		Title:    "Requests By Vhost",
		Units:    "requests/s",
		Fam:      "vhost",
		Ctx:      "web_log.vhost_requests",
		Type:     module.Stacked,
		Priority: prioReqVhost,
	}
	reqByPort = Chart{
		ID:       "port_requests",
		Title:    "Requests By Port",
		Units:    "requests/s",
		Fam:      "port",
		Ctx:      "web_log.port_requests",
		Type:     module.Stacked,
		Priority: prioReqPort,
	}
	reqByScheme = Chart{
		ID:       "scheme_requests",
		Title:    "Requests By Scheme",
		Units:    "requests/s",
		Fam:      "scheme",
		Ctx:      "web_log.scheme_requests",
		Type:     module.Stacked,
		Priority: prioReqScheme,
		Dims: Dims{
			{ID: "req_http_scheme", Name: "http", Algo: module.Incremental},
			{ID: "req_https_scheme", Name: "https", Algo: module.Incremental},
		},
	}
	reqByMethod = Chart{
		ID:       "http_method_requests",
		Title:    "Requests By HTTP Method",
		Units:    "requests/s",
		Fam:      "http method",
		Ctx:      "web_log.http_method_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	reqByVersion = Chart{
		ID:       "http_version_requests",
		Title:    "Requests By HTTP Version",
		Units:    "requests/s",
		Fam:      "http version",
		Ctx:      "web_log.http_version_requests",
		Type:     module.Stacked,
		Priority: prioReqVersion,
	}
	reqByIPProto = Chart{
		ID:       "ip_proto_requests",
		Title:    "Requests By IP Protocol",
		Units:    "requests/s",
		Fam:      "ip proto",
		Ctx:      "web_log.ip_proto_requests",
		Type:     module.Stacked,
		Priority: prioReqIPProto,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	reqBySSLProto = Chart{
		ID:       "ssl_proto_requests",
		Title:    "Requests By SSL Connection Protocol",
		Units:    "requests/s",
		Fam:      "ssl conn",
		Ctx:      "web_log.ssl_proto_requests",
		Type:     module.Stacked,
		Priority: prioReqSSLProto,
	}
	reqBySSLCipherSuite = Chart{
		ID:       "ssl_cipher_suite_requests",
		Title:    "Requests By SSL Connection Cipher Suite",
		Units:    "requests/s",
		Fam:      "ssl conn",
		Ctx:      "web_log.ssl_cipher_suite_requests",
		Type:     module.Stacked,
		Priority: prioReqSSLCipherSuite,
	}
)

// Request By N Patterns
var (
	reqByURLPattern = Chart{
		ID:       "url_pattern_requests",
		Title:    "URL Field, Requests By Pattern",
		Units:    "requests/s",
		Fam:      "url ptn",
		Ctx:      "web_log.url_pattern_requests",
		Type:     module.Stacked,
		Priority: prioReqURLPattern,
	}
	reqByCustomPattern = Chart{
		ID:       "custom_field_%s_pattern_requests",
		Title:    "Custom Field %s, Requests By Pattern",
		Units:    "requests/s",
		Fam:      "custom field",
		Ctx:      "web_log.custom_field_%s_pattern_requests",
		Type:     module.Stacked,
		Priority: prioReqCustomFieldPattern,
	}
)

// URL pattern stats
var (
	urlPatternRespCodes = Chart{
		ID:       "url_pattern_%s_status_code_responses",
		Title:    "Responses By Status Code",
		Units:    "responses/s",
		Fam:      "url ptn %s",
		Ctx:      "web_log.url_pattern_%s_status_code_responses",
		Type:     module.Stacked,
		Priority: prioURLPatternStats,
	}
	urlPatternBandwidth = Chart{
		ID:       "url_pattern_%s_bandwidth",
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "url ptn %s",
		Ctx:      "web_log.url_pattern_%s_bandwidth",
		Type:     module.Area,
		Priority: prioURLPatternStats + 1,
		Dims: Dims{
			{ID: "url_ptn_%s_bytes_received", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: "url_ptn_%s_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}
	urlPatternReqProcTime = Chart{
		ID:       "url_pattern_%s_request_processing_time",
		Title:    "Request Processing Time",
		Units:    "milliseconds",
		Fam:      "url ptn %s",
		Ctx:      "web_log.url_pattern_%s_request_processing_time",
		Type:     module.Area,
		Priority: prioURLPatternStats + 2,
		Dims: Dims{
			{ID: "url_ptn_%s_req_proc_time_min", Name: "min", Algo: module.Incremental, Div: 1000},
			{ID: "url_ptn_%s_req_proc_time_max", Name: "max", Algo: module.Incremental, Div: 1000},
			{ID: "url_ptn_%s_req_proc_time_avg", Name: "avg", Algo: module.Incremental, Div: 1000},
		},
	}
)

func newReqProcTimeHistChart(histogram []float64) *Chart {
	chart := reqProcTimeHist.Copy()
	for i, v := range histogram {
		dim := &Dim{
			ID:   fmt.Sprintf("req_proc_time_hist_bucket_%d", i+1),
			Name: fmt.Sprintf("%.3f", v),
			Algo: module.Incremental,
		}
		check(chart.AddDim(dim))
	}
	check(chart.AddDim(&Dim{
		ID:   "req_proc_time_hist_count",
		Name: "+Inf",
		Algo: module.Incremental,
	}))
	return chart
}

func newUpsRespTimeHistChart(histogram []float64) *Chart {
	chart := upsRespTimeHist.Copy()
	for i, v := range histogram {
		dim := &Dim{
			ID:   fmt.Sprintf("upstream_resp_time_hist_bucket_%d", i+1),
			Name: fmt.Sprintf("%.3f", v),
			Algo: module.Incremental,
		}
		check(chart.AddDim(dim))
	}
	check(chart.AddDim(&Dim{
		ID:   "upstream_resp_time_hist_count",
		Name: "+Inf",
		Algo: module.Incremental,
	}))
	return chart
}

func newReqByURLPatternChart(patterns []userPattern) *Chart {
	chart := reqByURLPattern.Copy()
	for _, p := range patterns {
		dim := &Dim{
			ID:   "req_url_ptn_" + p.Name,
			Name: p.Name,
			Algo: module.Incremental,
		}
		check(chart.AddDim(dim))
	}
	return chart
}

func newReqByCustomPatternCharts(fields []customField) []*Chart {
	charts := Charts{}
	for _, f := range fields {
		chart := newReqByCustomPatternChart(f)
		check(charts.Add(chart))
	}
	return charts
}

func newReqByCustomPatternChart(f customField) *Chart {
	chart := reqByCustomPattern.Copy()
	chart.ID = fmt.Sprintf(chart.ID, f.Name)
	chart.Title = fmt.Sprintf(chart.Title, f.Name)
	chart.Ctx = fmt.Sprintf(chart.Ctx, f.Name)
	for _, p := range f.Patterns {
		dim := &Dim{
			ID:   fmt.Sprintf("req_custom_field_%s_%s", f.Name, p.Name),
			Name: p.Name,
			Algo: module.Incremental,
		}
		check(chart.AddDim(dim))
	}
	return chart
}

func newURLPatternRespStatusCodeChart(name string) *Chart {
	chart := urlPatternRespCodes.Copy()
	chart.ID = fmt.Sprintf(chart.ID, name)
	chart.Fam = fmt.Sprintf(chart.Fam, name)
	chart.Ctx = fmt.Sprintf(chart.Ctx, name)
	return chart
}

func newURLPatternBandwidthChart(name string) *Chart {
	chart := urlPatternBandwidth.Copy()
	chart.ID = fmt.Sprintf(chart.ID, name)
	chart.Fam = fmt.Sprintf(chart.Fam, name)
	chart.Ctx = fmt.Sprintf(chart.Ctx, name)
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, name)
	}
	return chart
}

func newURLPatternReqProcTimeChart(name string) *Chart {
	chart := urlPatternReqProcTime.Copy()
	chart.ID = fmt.Sprintf(chart.ID, name)
	chart.Fam = fmt.Sprintf(chart.Fam, name)
	chart.Ctx = fmt.Sprintf(chart.Ctx, name)
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, name)
	}
	return chart
}

// TODO: this method is hard to read, should be refactored/simplified
func (w WebLog) createCharts(line *logLine) *Charts {
	// Following charts are created during runtime:
	//   - reqBySSLProto, reqBySSLCipherSuite - it is likely line has no SSL stuff at this moment
	charts := Charts{
		reqTotal.Copy(),
		reqExcluded.Copy(),
		respCodeClass.Copy(),
		respStatuses.Copy(),
	}
	if !w.GroupRespCodes {
		check(charts.Add(respCodes.Copy()))
	} else {
		check(charts.Add(respCodes1xx.Copy()))
		check(charts.Add(respCodes2xx.Copy()))
		check(charts.Add(respCodes3xx.Copy()))
		check(charts.Add(respCodes4xx.Copy()))
		check(charts.Add(respCodes5xx.Copy()))
	}
	if line.hasVhost() {
		check(charts.Add(reqByVhost.Copy()))
	}
	if line.hasPort() {
		check(charts.Add(reqByPort.Copy()))
	}
	if line.hasReqScheme() {
		check(charts.Add(reqByScheme.Copy()))
	}
	if line.hasReqClient() {
		check(charts.Add(reqByIPProto.Copy()))
		check(charts.Add(uniqIPsCurPoll.Copy()))
	}
	if line.hasReqMethod() {
		check(charts.Add(reqByMethod.Copy()))
	}
	if line.hasReqURL() && len(w.URLPatterns) > 0 {
		chart := newReqByURLPatternChart(w.URLPatterns)
		check(charts.Add(chart))

		for _, p := range w.URLPatterns {
			chart := newURLPatternRespStatusCodeChart(p.Name)
			check(charts.Add(chart))
		}
	}
	if line.hasReqProto() {
		check(charts.Add(reqByVersion.Copy()))
	}
	if line.hasReqSize() || line.hasRespSize() {
		check(charts.Add(bandwidth.Copy()))

		for _, p := range w.URLPatterns {
			chart := newURLPatternBandwidthChart(p.Name)
			check(charts.Add(chart))
		}
	}
	if line.hasReqProcTime() {
		check(charts.Add(reqProcTime.Copy()))
		if len(w.Histogram) != 0 {
			chart := newReqProcTimeHistChart(w.Histogram)
			check(charts.Add(chart))
		}

		for _, p := range w.URLPatterns {
			chart := newURLPatternReqProcTimeChart(p.Name)
			check(charts.Add(chart))
		}
	}
	if line.hasUpstreamRespTime() {
		check(charts.Add(upsRespTime.Copy()))
		if len(w.Histogram) != 0 {
			chart := newUpsRespTimeHistChart(w.Histogram)
			check(charts.Add(chart))
		}
	}
	if line.hasCustom() && len(w.CustomFields) > 0 {
		cs := newReqByCustomPatternCharts(w.CustomFields)
		check(charts.Add(cs...))
	}

	return &charts
}

func (w *WebLog) addDimToVhostChart(vhost string) {
	chart := w.Charts().Get(reqByVhost.ID)
	dim := &Dim{
		ID:   "req_vhost_" + vhost,
		Name: vhost,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToPortChart(port string) {
	chart := w.Charts().Get(reqByPort.ID)
	dim := &Dim{
		ID:   "req_port_" + port,
		Name: port,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToReqMethodChart(method string) {
	chart := w.Charts().Get(reqByMethod.ID)
	dim := &Dim{
		ID:   "req_method_" + method,
		Name: method,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToReqVersionChart(version string) {
	chart := w.Charts().Get(reqByVersion.ID)
	dim := &Dim{
		ID:   "req_version_" + version,
		Name: version,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToSSLProtoChart(proto string) {
	chart := w.Charts().Get(reqBySSLProto.ID)
	if chart == nil {
		chart = reqBySSLProto.Copy()
		check(w.Charts().Add(chart))
	}
	dim := &Dim{
		ID:   "req_ssl_proto_" + proto,
		Name: proto,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToSSLCipherSuiteChart(cipher string) {
	chart := w.Charts().Get(reqBySSLCipherSuite.ID)
	if chart == nil {
		chart = reqBySSLCipherSuite.Copy()
		check(w.Charts().Add(chart))
	}
	dim := &Dim{
		ID:   "req_ssl_cipher_suite_" + cipher,
		Name: cipher,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToRespCodesChart(code string) {
	chart := w.findRespCodesChart(code)
	if chart == nil {
		return
	}
	dim := &Dim{
		ID:   "resp_status_code_" + code,
		Name: code,
		Algo: module.Incremental,
	}
	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToURLPatternRespCodesChart(name, code string) {
	id := fmt.Sprintf(urlPatternRespCodes.ID, name)
	chart := w.Charts().Get(id)
	dim := &Dim{
		ID:   fmt.Sprintf("url_ptn_%s_resp_status_code_%s", name, code),
		Name: code,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) findRespCodesChart(code string) *Chart {
	if !w.GroupRespCodes {
		return w.Charts().Get(respCodes.ID)
	}

	var chart Chart
	switch class := code[:1]; class {
	case "1":
		chart = respCodes1xx
	case "2":
		chart = respCodes2xx
	case "3":
		chart = respCodes3xx
	case "4":
		chart = respCodes4xx
	case "5":
		chart = respCodes5xx
	default:
		return nil
	}
	return w.Charts().Get(chart.ID)
}

// TODO: get rid of
func check(err error) {
	if err != nil {
		panic(err)
	}
}
