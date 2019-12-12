package squidlog

import (
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
	// netdata specific grouping
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
	reqByMethod = Chart{
		ID:       "requests_by_http_method",
		Title:    "Requests By HTTP Method",
		Units:    "requests/s",
		Fam:      "http method",
		Ctx:      "squid.http_method_requests",
		Type:     module.Stacked,
		Priority: prioReqMethod,
	}
)

func (s *SquidLog) createCharts(line *logLine) error {
	return nil
}
