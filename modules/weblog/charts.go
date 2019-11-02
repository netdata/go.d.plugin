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
	prioRespStatuses    = defaultPriority
	prioRespCodesPerFam = defaultPriority + 1
	prioRespCodes       = defaultPriority + 2 // 5 charts
	prioBandwidth       = defaultPriority + 10
	prioRespTime        = defaultPriority + 11
	prioRespTimeHist    = defaultPriority + 12
	prioUpsRespTime     = defaultPriority + 13
	prioUpsRespTimeHist = defaultPriority + 14
	prioUniqIP          = defaultPriority + 15
	prioRequests        = defaultPriority + 16
	prioReqVhost        = defaultPriority + 17
	prioReqPort         = defaultPriority + 18
	prioReqScheme       = defaultPriority + 19
	prioReqMethod       = defaultPriority + 20
	prioReqVersion      = defaultPriority + 21
	prioReqIPProto      = defaultPriority + 22
	prioReqCustom       = defaultPriority + 23
	prioReqURL          = defaultPriority + 24
	prioURLStats        = defaultPriority + 25 // 3 charts per URL TODO: order?
)

// NOTE: inconsistency between contexts with python web_log
// TODO: current histogram charts are misleading in netdata

// Resp Statuses       [responses]
// Resp Codes Per Fam  [responses]
// Resp Codes          [responses]
// Bandwidth           [bandwidth]
// Resp Time           [timings]
// Resp Time Hist      [timings]
// Resp Time Ups       [upstream]
// Resp Time Hist Ups  [upstream]
// Uniq IPs            [clients]
// Requests            [requests]
// Req Per Vhost       [requests]
// Req Per Port        [requests]
// Req Per Scheme      [requests]
// Req Per Method      [requests]
// Req Per Version     [requests]
// Req Per IP Proto    [requests]
// Req Per Custom      [requests]
// Req Per URL         [requests]
// URL Stats           [requests]

var charts = Charts{
	respCodesFam.Copy(),
	respStatuses.Copy(),
}

// Responses
var (
	respCodesFam = Chart{
		ID:       "response_codes_family",
		Title:    "Response Codes Per Family",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_codes_family",
		Type:     module.Stacked,
		Priority: prioRespCodesPerFam,
		Dims: Dims{
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
		},
	}
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
		ID:       "response_status_codes",
		Title:    "Response Status Codes",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "web_log.response_status_codes",
		Type:     module.Stacked,
		Priority: prioRespCodes,
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
		ID:       "response_time",
		Title:    "Processing Time",
		Units:    "milliseconds",
		Fam:      "timings",
		Ctx:      "web_log.response_time",
		Priority: prioRespTime,
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Div: 1000},
			{ID: "resp_time_max", Name: "max", Div: 1000},
			{ID: "resp_time_avg", Name: "avg", Div: 1000},
		},
	}
	respTimeHist = Chart{
		ID:       "response_time_histogram",
		Title:    "Processing Time Histogram",
		Units:    "requests/s",
		Fam:      "timings",
		Ctx:      "web_log.response_time_histogram",
		Priority: prioRespTimeHist,
	}
)

// Upstream
var (
	upsRespTime = Chart{
		ID:       "upstream_response_time",
		Title:    "Upstream Processing Time",
		Units:    "milliseconds",
		Fam:      "timings",
		Ctx:      "web_log.upstream_response_time",
		Priority: prioUpsRespTime,
		Dims: Dims{
			{ID: "resp_time_upstream_min", Name: "min", Div: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Div: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Div: 1000},
		},
	}
	upsRespTimeHist = Chart{
		ID:       "upstream_response_time_histogram",
		Title:    "Upstream Processing Time Histogram",
		Units:    "requests/s",
		Fam:      "timings",
		Ctx:      "web_log.upstream_response_time_histogram",
		Priority: prioUpsRespTimeHist,
	}
)

// Clients
var (
	uniqIPsCurPoll = Chart{
		ID:       "uniq_ips_current_poll",
		Title:    "Unique Clients Current Poll",
		Units:    "ips",
		Fam:      "clients",
		Ctx:      "web_log.uniq_ips_current_poll",
		Type:     module.Stacked,
		Priority: prioUniqIP,
		Dims: Dims{
			{ID: "unique_current_poll_ipv4", Name: "ipv4", Algo: module.Absolute},
			{ID: "unique_current_poll_ipv6", Name: "ipv6", Algo: module.Absolute},
		},
	}
)

// Requests
var (
	reqTotal = Charts{
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
	}
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
		Fam:      "req methods",
		Ctx:      "web_log.requests_http_method",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	reqPerVersion = Chart{
		ID:       "requests_http_version",
		Title:    "Requests Per HTTP Version",
		Units:    "requests/s",
		Fam:      "req versions",
		Ctx:      "web_log.requests_http_version",
		Type:     module.Stacked,
		Priority: prioReqVersion,
	}
	reqPerIPProto = Chart{
		ID:       "requests_ip_proto",
		Title:    "Requests Per IP Protocol",
		Units:    "requests/s",
		Fam:      "req ip protocols",
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
		Title:    "Requests Per User Defined Custom Category",
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
		Fam:      "req urls",
		Ctx:      "web_log.requests_url",
		Type:     module.Stacked,
		Priority: prioReqURL,
	}
)

func newURLRespCodesChart(urlCat string) *Chart {
	return &Chart{
		ID:       respCodes.ID + "_" + urlCat,
		Title:    "Response Status Codes",
		Units:    "responses/s",
		Fam:      "url " + urlCat,
		Ctx:      "web_log.response_status_codes_per_url",
		Type:     module.Stacked,
		Priority: prioURLStats,
	}
}

func newURLBandwidthChart(urlCat string) *Chart {
	return &Chart{
		ID:       bandwidth.ID + "_" + urlCat,
		Title:    "Bandwidth",
		Units:    "kilobits/s",
		Fam:      "url " + urlCat,
		Ctx:      "web_log.bandwidth_per_url",
		Type:     module.Area,
		Priority: prioURLStats + 1,
		Dims: Dims{
			{ID: urlCat + "_resp_length", Name: "received", Algo: module.Incremental, Mul: 8, Div: 1000},
			{ID: urlCat + "_bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}
}

func newURLRespTimeChart(urlCat string) *Chart {
	return &Chart{
		ID:       respTime.ID + "_" + urlCat,
		Title:    "Processing Time",
		Units:    "milliseconds",
		Fam:      "url " + urlCat,
		Ctx:      "web_log.response_time_per_url",
		Type:     module.Area,
		Priority: prioURLStats + 2,
		Dims: Dims{
			{ID: urlCat + "_resp_time_min", Name: "min", Algo: module.Incremental, Div: 1000},
			{ID: urlCat + "_resp_time_max", Name: "max", Algo: module.Incremental, Div: 1000},
			{ID: urlCat + "_resp_time_avg", Name: "avg", Algo: module.Incremental, Div: 1000},
		},
	}
}

func newRespCodesDetailedPerFamCharts() []*Chart {
	return []*Chart{
		{
			ID:       respCodes.ID + "_1xx",
			Title:    "Response Status Codes 1xx",
			Units:    "requests/s",
			Fam:      "responses",
			Ctx:      "web_log.response_status_codes_1xx",
			Type:     module.Stacked,
			Priority: prioRespCodes,
		},
		{
			ID:       respCodes.ID + "_2xx",
			Title:    "Response Status Codes 2xx",
			Units:    "requests/s",
			Fam:      "responses",
			Ctx:      "web_log.response_status_codes_2xx",
			Type:     module.Stacked,
			Priority: prioRespCodes + 1,
		},
		{
			ID:       respCodes.ID + "_3xx",
			Title:    "Response Status Codes 3xx",
			Units:    "requests/s",
			Fam:      "responses",
			Ctx:      "web_log.response_status_codes_3xx",
			Type:     module.Stacked,
			Priority: prioRespCodes + 2,
		},
		{
			ID:       respCodes.ID + "_4xx",
			Title:    "Response Status Codes 4xx",
			Units:    "requests/s",
			Fam:      "responses",
			Ctx:      "web_log.response_status_codes_4xx",
			Type:     module.Stacked,
			Priority: prioRespCodes + 3,
		},
		{
			ID:       respCodes.ID + "_5xx",
			Title:    "Response Status Codes 5xx",
			Units:    "requests/s",
			Fam:      "responses",
			Ctx:      "web_log.response_status_codes_5xx",
			Type:     module.Stacked,
			Priority: prioRespCodes + 4,
		},
	}
}

func newReqPerURLCatsChart(cats []*category) *Chart {
	chart := reqPerURL.Copy()
	for _, c := range cats {
		dim := &Dim{
			ID:   "req_url_" + c.name,
			Name: c.name,
			Algo: module.Incremental,
		}
		panicIfErr(chart.AddDim(dim))
	}
	return chart
}

func newReqPerCustomCatsChart(cats []*category) *Chart {
	chart := reqPerCustom.Copy()
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

func newRespTimeHistogramChart(histogram []float64) *Chart {
	chart := respTimeHist.Copy()
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

func newUpsRespTimeHistogramChart(histogram []float64) *Chart {
	chart := upsRespTimeHist.Copy()
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
	if w.col.port {
		w.addPortChart()
		w.updatePortChart()
	}
	if w.col.scheme {
		w.addSchemeChart()
	}
	if w.col.client {
		w.addClientCharts()
	}
	if w.col.method {
		w.addHTTPMethodChart()
		w.updateReqMethodChart()
	}
	if w.col.url {
		w.addURLChart()
	}
	if w.col.version {
		w.addHTTPVersionChart()
		w.updateReqVersionChart()
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

	if w.col.url {
		w.addPerURLRespCodesCharts()
		if w.col.reqSize || w.col.respSize {
			w.addPerURLBandwidthCharts()
		}
		if w.col.respTime {
			w.addPerURLRespTimeCharts()
		}
	}
}

func (w *WebLog) addVhostChart() {
	if w.chartsCache.created.addIfNotExist(reqPerVhost.ID) {
		return
	}
	panicIfErr(w.Charts().Add(reqPerVhost.Copy()))
}

func (w *WebLog) addPortChart() {
	if w.chartsCache.created.addIfNotExist(reqPerPort.ID) {
		return
	}
	panicIfErr(w.Charts().Add(reqPerPort.Copy()))
}

func (w *WebLog) addSchemeChart() {
	if w.chartsCache.created.addIfNotExist(reqPerScheme.ID) {
		return
	}
	panicIfErr(w.Charts().Add(reqPerScheme.Copy()))
}

func (w *WebLog) addClientCharts() {
	if w.chartsCache.created.addIfNotExist(reqPerIPProto.ID) {
		return
	}

	panicIfErr(w.Charts().Add(reqPerIPProto.Copy()))
	panicIfErr(w.Charts().Add(uniqIPsCurPoll.Copy()))
}

func (w *WebLog) addHTTPMethodChart() {
	if w.chartsCache.created.addIfNotExist(reqPerMethod.ID) {
		return
	}

	panicIfErr(w.Charts().Add(reqPerMethod.Copy()))
}

func (w *WebLog) addURLChart() {
	if w.chartsCache.created.addIfNotExist(reqPerURL.ID) {
		return
	}

	chart := newReqPerURLCatsChart(w.urlCats)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addPerURLRespCodesCharts() {
	if w.chartsCache.created.addIfNotExist(respCodes.ID + "_per_url") {
		return
	}

	for _, c := range w.urlCats {
		panicIfErr(w.Charts().Add(newURLRespCodesChart(c.name)))
	}
}

func (w *WebLog) addPerURLBandwidthCharts() {
	if w.chartsCache.created.addIfNotExist(bandwidth.ID + "_per_url") {
		return
	}

	for _, c := range w.urlCats {
		panicIfErr(w.Charts().Add(newURLBandwidthChart(c.name)))
	}
}

func (w *WebLog) addPerURLRespTimeCharts() {
	if w.chartsCache.created.addIfNotExist(respTime.ID + "_per_url") {
		return
	}

	for _, c := range w.urlCats {
		panicIfErr(w.Charts().Add(newURLRespTimeChart(c.name)))
	}
}

func (w *WebLog) addHTTPVersionChart() {
	if w.chartsCache.created.addIfNotExist(reqPerVersion.ID) {
		return
	}

	panicIfErr(w.Charts().Add(reqPerVersion.Copy()))
}

func (w *WebLog) addRespCodesDetailedChart() {
	if w.chartsCache.created.addIfNotExist(respCodes.ID) {
		return
	}

	if w.AggregateResponseCodes {
		panicIfErr(w.Charts().Add(respCodes.Copy()))
		return
	}

	// TODO: do not create charts for all families
	panicIfErr(w.Charts().Add(newRespCodesDetailedPerFamCharts()...))
}

func (w *WebLog) addBandwidthChart() {
	if w.chartsCache.created.addIfNotExist(bandwidth.ID) {
		return
	}

	panicIfErr(w.Charts().Add(bandwidth.Copy()))
}

func (w *WebLog) addRespTimeCharts() {
	if w.chartsCache.created.addIfNotExist(respTime.ID) {
		return
	}

	panicIfErr(w.Charts().Add(respTime.Copy()))

	if len(w.Histogram) == 0 {
		return
	}

	chart := newRespTimeHistogramChart(w.Histogram)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addUpstreamRespTimeCharts() {
	if w.chartsCache.created.addIfNotExist(upsRespTime.ID) {
		return
	}

	panicIfErr(w.Charts().Add(upsRespTime.Copy()))

	if len(w.Histogram) == 0 {
		return
	}

	chart := newUpsRespTimeHistogramChart(w.Histogram)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) addCustomChart() {
	if w.chartsCache.created.addIfNotExist(reqPerCustom.ID) {
		return
	}

	chart := newReqPerCustomCatsChart(w.userCats)
	panicIfErr(w.Charts().Add(chart))
}

func (w *WebLog) updateVhostChart() {
	chart := w.Charts().Get(reqPerVhost.ID)

	for v := range w.mx.ReqVhost {
		if w.chartsCache.vhosts.addIfNotExist(v) {
			continue
		}
		addDimToVhostChart(chart, v)
	}
}

func (w *WebLog) updatePortChart() {
	chart := w.Charts().Get(reqPerPort.ID)

	for v := range w.mx.ReqPort {
		if w.chartsCache.ports.addIfNotExist(v) {
			continue
		}
		addDimToPortChart(chart, v)
	}
}

func (w *WebLog) updateReqMethodChart() {
	chart := w.Charts().Get(reqPerMethod.ID)

	for v := range w.mx.ReqMethod {
		if w.chartsCache.methods.addIfNotExist(v) {
			continue
		}
		addDimToReqMethodChart(chart, v)
	}
}

func (w *WebLog) updateReqVersionChart() {
	chart := w.Charts().Get(reqPerVersion.ID)

	for v := range w.mx.ReqVersion {
		if w.chartsCache.versions.addIfNotExist(v) {
			continue
		}
		addDimToReqVersionChart(chart, v)
	}
}

func (w *WebLog) updateRespCodesDetailedChart() {
	for v := range w.mx.RespCode {
		if w.chartsCache.codes.addIfNotExist(v) {
			continue
		}
		chart := w.respCodesDetailedChartByCode(v)
		addDimToRespCodesDetailedChart(chart, v)
	}
}

func (w *WebLog) respCodesDetailedChartByCode(code string) *Chart {
	if w.AggregateResponseCodes {
		return w.Charts().Get(respCodes.ID)
	}

	var id string
	switch v := code[:1]; v {
	case "1", "2", "3", "4", "5":
		id = fmt.Sprintf("%s_%sxx", respCodes.ID, v)
	default:
		// TODO: delete or add
		id = fmt.Sprintf("%s_other", respCodes.ID)
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

func addDimToReqVersionChart(chart *Chart, version string) {
	dimID := "req_version_" + version
	dim := &Dim{
		ID:   dimID,
		Name: version,
		Algo: module.Incremental,
	}
	panicIfErr(chart.AddDim(dim))
	chart.MarkNotCreated()
}

func addDimToReqMethodChart(chart *Chart, method string) {
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

func addDimToPortChart(chart *Chart, port string) {
	dimID := "req_port_" + port
	dim := &Dim{
		ID:   dimID,
		Name: port,
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

type (
	cache map[string]struct{}

	chartsCache struct {
		created  cache
		vhosts   cache
		ports    cache
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
