package squidlog

import (
	"errors"

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
	prioReqTotal = orchestrator.DefaultJobPriority + iota
	prioReqExcluded
	prioReqType

	prioRespCodesClass
	prioRespCodes

	prioBandwidth
	prioRespTime

	prioUniqClients
	prioReqMethod
)

// Requests
var (
	reqTotal = Chart{
		ID:       "requests",
		Title:    "Total Requests",
		Units:    "requests/s",
		Fam:      "requests",
		Ctx:      "squid.requests",
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
		Ctx:      "squid.excluded_requests",
		Type:     module.Stacked,
		Priority: prioReqExcluded,
		Dims: Dims{
			{ID: "req_unmatched", Name: "unmatched", Algo: module.Incremental},
		},
	}
	reqTypes = Chart{
		ID:       "requests_by_type",
		Title:    "Requests By Type",
		Units:    "requests/s",
		Fam:      "requests",
		Ctx:      "squid.type_requests",
		Type:     module.Stacked,
		Priority: prioReqType,
		Dims: Dims{
			{ID: "req_type_success", Name: "success", Algo: module.Incremental},
			{ID: "req_type_bad", Name: "bad", Algo: module.Incremental},
			{ID: "req_type_redirect", Name: "redirect", Algo: module.Incremental},
			{ID: "req_type_error", Name: "error", Algo: module.Incremental},
		},
	}
)

// Responses
var (
	respCodeClass = Chart{
		ID:       "responses_by_status_code_class",
		Title:    "Responses By Status Code Class",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "squid.status_code_class_responses",
		Type:     module.Stacked,
		Priority: prioRespCodesClass,
		Dims: Dims{
			{ID: "resp_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "resp_1xx", Name: "1xx", Algo: module.Incremental},
			{ID: "resp_0xx", Name: "0xx", Algo: module.Incremental},
			{ID: "resp_6xx", Name: "6xx", Algo: module.Incremental},
		},
	}
	respCodes = Chart{
		ID:       "responses_by_status_code",
		Title:    "Responses By Status Code",
		Units:    "responses/s",
		Fam:      "responses",
		Ctx:      "squid.status_code_responses",
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
		Ctx:      "squid.bandwidth",
		Type:     module.Area,
		Priority: prioBandwidth,
		Dims: Dims{
			{ID: "bytes_sent", Name: "sent", Algo: module.Incremental, Mul: -8, Div: 1000},
		},
	}
)

// Clients
var (
	uniqClientsCurPoll = Chart{
		ID:       "uniq_clients",
		Title:    "Unique Clients",
		Units:    "clients/s",
		Fam:      "client",
		Ctx:      "squid.uniq_clients",
		Priority: prioUniqClients,
		Dims: Dims{
			{ID: "uniq_clients", Name: "clients"},
		},
	}
)

var (
	respTime = Chart{
		ID:       "response_time",
		Title:    "Response Time",
		Units:    "milliseconds",
		Fam:      "timings",
		Ctx:      "squid.request_processing_time",
		Priority: prioRespTime,
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Div: 1000},
			{ID: "resp_time_max", Name: "max", Div: 1000},
			{ID: "resp_time_avg", Name: "avg", Div: 1000},
		},
	}
)

var (
	reqByMethod = Chart{
		ID:       "requests_by_http_method",
		Title:    "Requests By HTTP Method",
		Units:    "requests/s",
		Fam:      "http method",
		Ctx:      "squid.http_method_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	reqByMimeType = Chart{
		ID:       "requests_by_mime_type",
		Title:    "Requests By MIME Type",
		Units:    "requests/s",
		Fam:      "mime type",
		Ctx:      "squid.mime_type_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	reqByHierCode = Chart{
		ID:       "requests_by_hier_code",
		Title:    "Requests By HIER Code",
		Units:    "requests/s",
		Fam:      "hier code",
		Ctx:      "squid.hier_code_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}

	reqByServer = Chart{
		ID:       "requests_by_server_address",
		Title:    "Requests By Server Address",
		Units:    "requests/s",
		Fam:      "http method",
		Ctx:      "squid.http_method_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
)

var (
	cacheCode = Chart{
		ID:       "cache_result_code",
		Title:    "Requests Cache Result Code",
		Units:    "result/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_result_code",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	cacheCodeTransport = Chart{
		ID:       "requests_by_cache_code_transport",
		Title:    "Requests By Cache Code Transport",
		Units:    "requests/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_code_transport_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	cacheCodeHandling = Chart{
		ID:       "requests_by_cache_code_handling",
		Title:    "Requests By Cache Code Handling",
		Units:    "requests/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_code_handling_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	cacheCodeObject = Chart{
		ID:       "requests_by_cache_code_object",
		Title:    "Requests By Cache Code Object",
		Units:    "requests/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_code_object_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	cacheCodeLoadSource = Chart{
		ID:       "requests_by_cache_code_load_source",
		Title:    "Requests By Cache Code Load Source",
		Units:    "requests/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_code_load_source_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
	cacheCodeError = Chart{
		ID:       "requests_by_cache_code_error",
		Title:    "Requests By Cache Code Error",
		Units:    "requests/s",
		Fam:      "cache code",
		Ctx:      "squid.cache_code_error_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
)

func (s *SquidLog) createCharts(line *logLine) error {
	if line.empty() {
		return errors.New("empty line")
	}
	s.charts = nil
	// Following charts are created during runtime:
	//   - reqBySSLProto, reqBySSLCipherSuite - it is likely line has no SSL stuff at this moment
	charts := &Charts{
		reqTotal.Copy(),
		reqExcluded.Copy(),
	}
	if line.hasRespTime() {
		if err := addRespTimeCharts(charts); err != nil {
			return err
		}
	}
	if line.hasClientAddress() {
		if err := addClientAddressCharts(charts); err != nil {
			return err
		}
	}
	if line.hasCacheCode() {
		if err := addCacheCodeCharts(charts); err != nil {
			return err
		}
	}
	if line.hasHTTPCode() {
		if err := addHTTPCodeCharts(charts); err != nil {
			return err
		}
	}
	if line.hasRespSize() {
		if err := addRespSizeCharts(charts); err != nil {
			return err
		}
	}
	if line.hasReqMethod() {
		if err := addMethodCharts(charts); err != nil {
			return err
		}
	}
	if line.hasHierCode() {
		if err := addHierCodeCharts(charts); err != nil {
			return err
		}
	}
	if line.hasServerAddress() {
		if err := addServerAddressCharts(charts); err != nil {
			return err
		}
	}
	if line.hasMimeType() {
		if err := addMimeTypeCharts(charts); err != nil {
			return err
		}
	}
	s.charts = charts
	return nil
}

func addRespTimeCharts(charts *Charts) error {
	return charts.Add(respTime.Copy())
}

func addClientAddressCharts(charts *Charts) error {
	return charts.Add(uniqClientsCurPoll.Copy())
}

func addCacheCodeCharts(charts *Charts) error {
	cs := []Chart{
		cacheCode,
		cacheCodeTransport,
		cacheCodeHandling,
		cacheCodeObject,
		cacheCodeLoadSource,
		cacheCodeError,
	}
	for _, chart := range cs {
		if err := charts.Add(chart.Copy()); err != nil {
			return err
		}
	}
	return nil
}
func addHTTPCodeCharts(charts *Charts) error {
	cs := []Chart{
		reqTypes,
		respCodeClass,
		respCodes,
	}
	for _, chart := range cs {
		if err := charts.Add(chart.Copy()); err != nil {
			return err
		}
	}
	return nil
}

func addRespSizeCharts(charts *Charts) error {
	return charts.Add(bandwidth.Copy())
}

func addMethodCharts(charts *Charts) error {
	return charts.Add(reqByMethod.Copy())
}

func addHierCodeCharts(charts *Charts) error {
	return charts.Add(reqByHierCode.Copy())
}
func addServerAddressCharts(charts *Charts) error {
	return charts.Add(reqByServer.Copy())
}

func addMimeTypeCharts(charts *Charts) error {
	return charts.Add(reqByMimeType.Copy())
}
