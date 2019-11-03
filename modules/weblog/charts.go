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
	defaultPriority     = orchestrator.DefaultJobPriority
	prioReqTotal        = defaultPriority
	prioReqUnreported   = defaultPriority + 1
	prioRespStatuses    = defaultPriority + 2
	prioRespCodesGroups = defaultPriority + 3
	prioRespCodes       = defaultPriority + 4
	prioRespCodes1xx    = defaultPriority + 5
	prioRespCodes2xx    = defaultPriority + 6
	prioRespCodes3xx    = defaultPriority + 7
	prioRespCodes4xx    = defaultPriority + 8
	prioRespCodes5xx    = defaultPriority + 9
	prioBandwidth       = defaultPriority + 10
	prioReqProcTime     = defaultPriority + 11
	prioRespTimeHist    = defaultPriority + 12
	prioUpsRespTime     = defaultPriority + 13
	prioUpsRespTimeHist = defaultPriority + 14
	prioUniqIP          = defaultPriority + 15
	prioReqVhost        = defaultPriority + 16
	prioReqPort         = defaultPriority + 17
	prioReqScheme       = defaultPriority + 18
	prioReqMethod       = defaultPriority + 19
	prioReqVersion      = defaultPriority + 20
	prioReqIPProto      = defaultPriority + 21
	prioReqCustom       = defaultPriority + 22
	prioReqURL          = defaultPriority + 23
	prioURLStats        = defaultPriority + 25 // 3 charts per URL TODO: order?
)

// NOTE: inconsistency between contexts with python web_log
// TODO: current histogram charts are misleading in netdata

// Total Requests       [requests]
// Unreported Requests  [requests]
// Resp Statuses        [responses]
// Resp Codes Per Group [responses]
// Resp Codes           [responses]
// Bandwidth            [bandwidth]
// Resp Time            [timings]
// Resp Time Hist       [timings]
// Resp Time Ups        [upstream]
// Resp Time Hist Ups   [upstream]
// Uniq IPs             [clients]
// Req Per Vhost        [req vhost]
// Req Per Port         [req port]
// Req Per Scheme       [req scheme]
// Req Per Method       [req method]
// Req Per Version      [req version]
// Req Per IP Proto     [req ip proto]
// Req Per Custom       [req custom]
// Req Per URL          [req url]
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

	reqUnreported = Chart{
		ID:       "requests_unreported",
		Title:    "Unreported Requests",
		Units:    "requests/s",
		Fam:      "requests",
		Ctx:      "web_log.requests_unreported",
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
	respCodesGroups = Chart{
		ID:       "response_codes_group",
		Title:    "Response Codes Per Group",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_group",
		Type:     module.Stacked,
		Priority: prioRespCodesGroups,
		Dims: Dims{
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
		},
	}
	// netdata specific grouping
	respStatuses = Chart{
		ID:       "response_statuses",
		Title:    "Response Statuses",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_statuses",
		Type:     module.Stacked,
		Priority: prioRespStatuses,
		Dims: Dims{
			{ID: "resp_successful", Name: "success", Algo: module.Incremental},
			{ID: "resp_client_error", Name: "bad", Algo: module.Incremental},
			{ID: "resp_redirect", Name: "redirect", Algo: module.Incremental},
			{ID: "resp_server_error", Name: "error", Algo: module.Incremental},
		},
	}
	respCodes = Chart{
		ID:       "response_codes",
		Title:    "Response Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes",
		Type:     module.Stacked,
		Priority: prioRespCodes,
	}
	respCodes1xx = Chart{
		ID:       "response_codes_1xx",
		Title:    "Informational Response Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_1xx",
		Type:     module.Stacked,
		Priority: prioRespCodes1xx,
	}
	respCodes2xx = Chart{
		ID:       "response_codes_2xx",
		Title:    "Successful Response Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_2xx",
		Type:     module.Stacked,
		Priority: prioRespCodes2xx,
	}
	respCodes3xx = Chart{
		ID:       "response_codes_3xx",
		Title:    "Redirects Response Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_3xx",
		Type:     module.Stacked,
		Priority: prioRespCodes3xx,
	}
	respCodes4xx = Chart{
		ID:       "response_codes_4xx",
		Title:    "Client Errors Response Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_4xx",
		Type:     module.Stacked,
		Priority: prioRespCodes4xx,
	}
	respCodes5xx = Chart{
		ID:       "response_codes_5xx",
		Title:    "Server Errors Codes 5xx",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_5xx",
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
	respTime = Chart{
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
	respTimeHist = Chart{
		ID:       "request_processing_time_histogram",
		Title:    "Request Processing Time Histogram",
		Units:    "requests/s",
		Fam:      "timings",
		Ctx:      "web_log.request_processing_time_histogram",
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
			{ID: "upstream_req_proc_time_min", Name: "min", Div: 1000},
			{ID: "upstream_req_proc_time_max", Name: "max", Div: 1000},
			{ID: "upstream_req_proc_time_avg", Name: "avg", Div: 1000},
		},
	}
	upsRespTimeHist = Chart{
		ID:       "upstream_response_time_histogram",
		Title:    "Upstream Response Time Histogram",
		Units:    "requests/s",
		Fam:      "timings",
		Ctx:      "web_log.upstream_response_time_histogram",
		Priority: prioUpsRespTimeHist,
	}
)

// Clients
var (
	uniqIPsCurPoll = Chart{
		ID:       "uniq_clients_current_poll",
		Title:    "Unique Clients Current Poll",
		Units:    "clients",
		Fam:      "clients",
		Ctx:      "web_log.uniq_clients_current_poll",
		Type:     module.Stacked,
		Priority: prioUniqIP,
		Dims: Dims{
			{ID: "uniq_ipv4", Name: "ipv4", Algo: module.Absolute},
			{ID: "uniq_ipv6", Name: "ipv6", Algo: module.Absolute},
		},
	}
)

// Requester Per N
var (
	reqPerVhost = Chart{
		ID:       "requests_vhost",
		Title:    "Requests Per Vhost",
		Units:    "requests/s",
		Fam:      "req vhost",
		Ctx:      "web_log.requests_vhost",
		Type:     module.Stacked,
		Priority: prioReqVhost,
	}
	reqPerPort = Chart{
		ID:       "requests_port",
		Title:    "Requests Per Port",
		Units:    "requests/s",
		Fam:      "req port",
		Ctx:      "web_log.requests_port",
		Type:     module.Stacked,
		Priority: prioReqPort,
	}
	reqPerScheme = Chart{
		ID:       "requests_scheme",
		Title:    "Requests Per Scheme",
		Units:    "requests/s",
		Fam:      "req scheme",
		Ctx:      "web_log.requests_scheme",
		Type:     module.Stacked,
		Priority: prioReqScheme,
		Dims: Dims{
			{ID: "req_http_scheme", Name: "http", Algo: module.Incremental},
			{ID: "req_https_scheme", Name: "https", Algo: module.Incremental},
		},
	}
	reqPerMethod = Chart{
		ID:       "requests_http_method",
		Title:    "Requests Per HTTP Method",
		Units:    "requests/s",
		Fam:      "req method",
		Ctx:      "web_log.requests_http_method",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	reqPerVersion = Chart{
		ID:       "requests_http_version",
		Title:    "Requests Per HTTP Version",
		Units:    "requests/s",
		Fam:      "req version",
		Ctx:      "web_log.requests_http_version",
		Type:     module.Stacked,
		Priority: prioReqVersion,
	}
	reqPerIPProto = Chart{
		ID:       "requests_ip_proto",
		Title:    "Requests Per IP Protocol",
		Units:    "requests/s",
		Fam:      "req ip protocol",
		Ctx:      "web_log.requests_ip_proto",
		Type:     module.Stacked,
		Priority: prioReqIPProto,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: module.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: module.Incremental},
		},
	}
	reqPerCustom = Chart{
		ID:       "requests_custom",
		Title:    "Requests Per User Defined Custom Categories",
		Units:    "requests/s",
		Fam:      "req custom",
		Ctx:      "web_log.requests_custom",
		Type:     module.Stacked,
		Priority: prioReqCustom,
	}
	reqPerURL = Chart{
		ID:       "requests_url",
		Title:    "Requests Per Url",
		Units:    "requests/s",
		Fam:      "req url",
		Ctx:      "web_log.requests_url",
		Type:     module.Stacked,
		Priority: prioReqURL,
	}
)

func newRespTimeHistChart(histogram []float64) *Chart {
	chart := respTimeHist.Copy()
	for i, v := range histogram {
		dimID := fmt.Sprintf("resp_time_hist_bucket_%d", i+1)
		name := fmt.Sprintf("%.3f", v)

		dim := &Dim{ID: dimID, Name: name, Algo: module.Incremental}
		check(chart.AddDim(dim))
	}
	check(chart.AddDim(&Dim{ID: "resp_time_hist_count", Name: "+Inf", Algo: module.Incremental}))
	return chart
}

func newUpsRespTimeHistChart(histogram []float64) *Chart {
	chart := upsRespTimeHist.Copy()
	for i, v := range histogram {
		dimID := fmt.Sprintf("upstream_resp_time_hist_bucket_%d", i+1)
		name := fmt.Sprintf("%.3f", v)

		dim := &Dim{ID: dimID, Name: name, Algo: module.Incremental}
		check(chart.AddDim(dim))
	}
	check(chart.AddDim(&Dim{
		ID:   "upstream_resp_time_hist_count",
		Name: "+Inf",
		Algo: module.Incremental,
	}))
	return chart
}

func newReqPerURLChart(cs []*category) *Chart {
	chart := reqPerURL.Copy()
	for _, c := range cs {
		dim := &Dim{ID: "req_url_" + c.name, Name: c.name, Algo: module.Incremental}
		check(chart.AddDim(dim))
	}
	return chart
}

func newReqPerCustomChart(cs []*category) *Chart {
	chart := reqPerCustom.Copy()
	for _, c := range cs {
		dim := &Dim{ID: "req_custom_" + c.name, Name: c.name, Algo: module.Incremental}
		check(chart.AddDim(dim))
	}
	return chart
}

func newURLRespCodesChart(name string) *Chart {
	return &Chart{
		ID:       respCodes.ID + "_" + name,
		Title:    "Response Codes",
		Units:    "responses/s",
		Fam:      "url " + name,
		Ctx:      "web_log.response_codes_per_url",
		Type:     module.Stacked,
		Priority: prioURLStats,
	}
}

func newURLBandwidthChart(name string) *Chart {
	return &Chart{
		ID:       bandwidth.ID + "_" + name,
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "url " + name,
		Ctx:      "web_log.bandwidth_per_url",
		Type:     module.Area,
		Priority: prioURLStats + 1,
		Dims: Dims{
			{ID: name + "_bytes_received", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: name + "_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}
}

func newURLRespTimeChart(name string) *Chart {
	return &Chart{
		ID:       respTime.ID + "_" + name,
		Title:    "Request Processing Time",
		Units:    "milliseconds",
		Fam:      "url " + name,
		Ctx:      "web_log.request_processing_time_per_url",
		Type:     module.Area,
		Priority: prioURLStats + 2,
		Dims: Dims{
			{ID: name + "_req_proc_time_min", Name: "min", Algo: module.Incremental, Div: 1000},
			{ID: name + "_req_proc_time_max", Name: "max", Algo: module.Incremental, Div: 1000},
			{ID: name + "_req_proc_time_avg", Name: "avg", Algo: module.Incremental, Div: 1000},
		},
	}
}

func (w *WebLog) createCharts(line *logLine) *Charts {
	charts := Charts{
		reqTotal.Copy(),
		reqUnreported.Copy(),
		respCodesGroups.Copy(),
		respStatuses.Copy(),
	}
	if w.GroupRespCodes {
		check(charts.Add(respCodes.Copy()))
	} else {
		// NOTE: per fam resp codes charts are added during runtime
	}
	if line.hasVhost() {
		check(charts.Add(reqPerVhost.Copy()))
	}
	if line.hasPort() {
		check(charts.Add(reqPerPort.Copy()))
	}
	if line.hasReqScheme() {
		check(charts.Add(reqPerScheme.Copy()))
	}
	if line.hasReqClient() {
		check(charts.Add(reqPerIPProto.Copy()))
		check(charts.Add(uniqIPsCurPoll.Copy()))
	}
	if line.hasReqMethod() {
		check(charts.Add(reqPerMethod.Copy()))
	}
	if line.hasReqURL() && len(w.catURL) > 0 {
		chart := newReqPerURLChart(w.catURL)
		check(charts.Add(chart))

		for _, c := range w.catURL {
			chart := newURLRespCodesChart(c.name)
			check(charts.Add(chart))
		}
	}
	if line.hasReqProto() {
		check(charts.Add(reqPerVersion.Copy()))
	}
	if line.hasReqSize() || line.hasRespSize() {
		check(charts.Add(bandwidth.Copy()))

		for _, c := range w.catURL {
			chart := newURLBandwidthChart(c.name)
			check(charts.Add(chart))
		}
	}
	if line.hasReqProcTime() {
		check(charts.Add(respTime.Copy()))
		if len(w.Histogram) != 0 {
			chart := newRespTimeHistChart(w.Histogram)
			check(charts.Add(chart))
		}

		for _, c := range w.catURL {
			chart := newURLRespTimeChart(c.name)
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
	if line.hasCustom() && len(w.catCustom) > 0 {
		chart := newReqPerCustomChart(w.catCustom)
		check(charts.Add(chart))
	}

	return &charts
}

func (w *WebLog) addDimToVhostChart(vhost string) {
	chart := w.Charts().Get(reqPerVhost.ID)

	dim := &Dim{
		ID:   "req_vhost_" + vhost,
		Name: vhost,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToPortChart(port string) {
	chart := w.Charts().Get(reqPerPort.ID)

	dim := &Dim{
		ID:   "req_port_" + port,
		Name: port,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToReqMethodChart(method string) {
	chart := w.Charts().Get(reqPerMethod.ID)

	dim := &Dim{
		ID:   "req_method_" + method,
		Name: method,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToReqVersionChart(version string) {
	chart := w.Charts().Get(reqPerVersion.ID)

	dim := &Dim{
		ID:   "req_version_" + version,
		Name: version,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToRespCodesChart(code string) {
	chart := w.findRespCodesChartByCode(code)
	if chart == nil {
		return
	}

	dim := &Dim{
		ID:   "resp_code_" + code,
		Name: code,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) addDimToURLRespCodesChart(name, code string) {
	chart := w.Charts().Get(respCodes.ID + "_" + name)

	dim := &Dim{
		ID:   name + "_resp_code_" + code,
		Name: code,
		Algo: module.Incremental,
	}

	check(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func (w *WebLog) findRespCodesChartByCode(code string) *Chart {
	if !w.GroupRespCodes {
		return w.Charts().Get(respCodes.ID)
	}

	var chart Chart
	switch v := code[:1]; v {
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

	if !w.Charts().Has(chart.ID) {
		check(w.Charts().Add(chart.Copy()))
	}
	return w.Charts().Get(chart.ID)
}

// TODO: get rid of
func check(err error) {
	if err != nil {
		panic(err)
	}
}
