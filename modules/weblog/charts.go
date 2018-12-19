package weblog

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Chart is an alias for modules.Chart
	Chart = modules.Chart
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
	// Dim is an alias for modules.Dim
	Dim = modules.Dim
)

// NOTE: inconsistency between contexts with python web_log
var (
	responseStatuses = Chart{
		ID:    "response_statuses",
		Title: "Response Statuses",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_statuses",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "successful_requests", Name: "success", Algo: modules.Incremental},
			{ID: "server_errors", Name: "error", Algo: modules.Incremental},
			{ID: "redirects", Name: "redirect", Algo: modules.Incremental},
			{ID: "bad_requests", Name: "bad", Algo: modules.Incremental},
			{ID: "other_requests", Name: "other", Algo: modules.Incremental},
		},
	}
	responseCodes = Chart{
		ID:    "response_codes",
		Title: "Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "2xx", Algo: modules.Incremental},
			{ID: "5xx", Algo: modules.Incremental},
			{ID: "3xx", Algo: modules.Incremental},
			{ID: "4xx", Algo: modules.Incremental},
			{ID: "1xx", Algo: modules.Incremental},
			{ID: "0xx", Algo: modules.Incremental},
			{ID: "unmatched", Algo: modules.Incremental},
		},
	}
	responseCodesDetailed = Chart{
		ID:    "detailed_response_codes",
		Title: "Detailed Response Codes",
		Units: "requests/s",
		Fam:   "responses",
		Ctx:   "web_log.response_codes_detailed",
		Type:  modules.Stacked,
	}
	bandwidth = Chart{
		ID:    "bandwidth",
		Title: "Bandwidth",
		Units: "kilobits/s",
		Fam:   "bandwidth",
		Ctx:   "web_log.bandwidth",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "resp_length", Name: "received", Algo: modules.Incremental, Mul: 8, Div: 1000},
			{ID: "bytes_sent", Name: "sent", Algo: modules.Incremental, Mul: -8, Div: 1000},
		},
	}
	responseTime = Chart{
		ID:    "response_time",
		Title: "Processing Time",
		Units: "milliseconds",
		Fam:   "timings",
		Ctx:   "web_log.response_time",
		Type:  modules.Area,
		Dims: Dims{
			{ID: "resp_time_min", Name: "min", Algo: modules.Incremental, Div: 1000},
			{ID: "resp_time_max", Name: "max", Algo: modules.Incremental, Div: 1000},
			{ID: "resp_time_avg", Name: "avg", Algo: modules.Incremental, Div: 1000},
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
		Type:  modules.Area,
		Dims: Dims{
			{ID: "resp_time_upstream_min", Name: "min", Algo: modules.Incremental, Div: 1000},
			{ID: "resp_time_upstream_max", Name: "max", Algo: modules.Incremental, Div: 1000},
			{ID: "resp_time_upstream_avg", Name: "avg", Algo: modules.Incremental, Div: 1000},
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
		Type:  modules.Stacked,
	}
	requestsPerUserDefined = Chart{
		ID:    "requests_per_user_defined",
		Title: "Requests Per User Defined Pattern",
		Units: "requests/s",
		Fam:   "user defined",
		Ctx:   "web_log.requests_per_user_defined",
		Type:  modules.Stacked,
	}
	requestsPerHTTPMethod = Chart{
		ID:    "requests_per_http_method",
		Title: "Requests Per HTTP Method",
		Units: "requests/s",
		Fam:   "http methods",
		Ctx:   "web_log.requests_per_http_method",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "GET", Algo: modules.Incremental},
		},
	}
	requestsPerHTTPVersion = Chart{
		ID:    "requests_per_http_version",
		Title: "Requests Per HTTP Version",
		Units: "requests/s",
		Fam:   "http versions",
		Ctx:   "web_log.requests_per_http_version",
		Type:  modules.Stacked,
	}
	requestsPerIPProto = Chart{
		ID:    "requests_per_ip_proto",
		Title: "Requests Per IP Protocol",
		Units: "requests/s",
		Fam:   "ip protocols",
		Ctx:   "web_log.requests_per_ip_proto",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "req_ipv4", Name: "ipv4", Algo: modules.Incremental},
			{ID: "req_ipv6", Name: "ipv6", Algo: modules.Incremental},
		},
	}
	requestsPerVhost = Chart{
		ID:    "requests_per_vhost",
		Title: "Requests Per Vhost",
		Units: "requests/s",
		Fam:   "vhost",
		Ctx:   "web_log.requests_per_vhost",
		Type:  modules.Stacked,
	}
	currentPollIPs = Chart{
		ID:    "clients_current",
		Title: "Current Poll Unique Client IPs",
		Units: "unique ips",
		Fam:   "clients",
		Ctx:   "web_log.current_poll_ips",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "unique_current_poll_ipv4", Name: "ipv4", Algo: modules.Incremental},
			{ID: "unique_current_poll_ipv6", Name: "ipv6", Algo: modules.Incremental},
		},
	}
	allTimeIPs = Chart{
		ID:    "clients_all_time",
		Title: "All Time Unique Client IPs",
		Units: "unique ips",
		Fam:   "clients",
		Ctx:   "web_log.all_time_ips",
		Type:  modules.Stacked,
		Dims: Dims{
			{ID: "unique_all_time_ipv4", Name: "ipv4"},
			{ID: "unique_all_time_ipv6", Name: "ipv6"},
		},
	}
)

func responseCodesDetailedPerFamily() []*Chart {
	return []*Chart{
		{
			ID:    responseCodesDetailed.ID + "_1xx",
			Title: "Detailed Response Codes 1xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_1xx",
			Type:  modules.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_2xx",
			Title: "Detailed Response Codes 2xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_2xx",
			Type:  modules.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_3xx",
			Title: "Detailed Response Codes 3xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_3xx",
			Type:  modules.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_4xx",
			Title: "Detailed Response Codes 4xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_4xx",
			Type:  modules.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_5xx",
			Title: "Detailed Response Codes 5xx",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_5xx",
			Type:  modules.Stacked,
		},
		{
			ID:    responseCodesDetailed.ID + "_other",
			Title: "Detailed Response Codes Other",
			Units: "requests/s",
			Fam:   "responses",
			Ctx:   "web_log.response_codes_detailed_other",
			Type:  modules.Stacked,
		},
	}
}

func perCategoryStats(id string) []*Chart {
	return []*Chart{
		{
			ID:    responseCodesDetailed.ID + "_" + id,
			Title: "Detailed Response Codes",
			Units: "requests/s",
			Fam:   "url " + id,
			Ctx:   "web_log.response_codes_detailed_per_url",
			Type:  modules.Stacked,
		},
		{
			ID:    bandwidth.ID + "_" + id,
			Title: "Bandwidth",
			Units: "kilobits/s",
			Fam:   "url " + id,
			Ctx:   "web_log.bandwidth_per_url",
			Type:  modules.Area,
			Dims: Dims{
				{ID: id + "_resp_length", Name: "received", Algo: modules.Incremental, Mul: 8, Div: 1000},
				{ID: id + "_bytes_sent", Name: "sent", Algo: modules.Incremental, Mul: -8, Div: 1000},
			},
		},
		{
			ID:    responseTime.ID + "_" + id,
			Title: "Processing Time",
			Units: "milliseconds",
			Fam:   "url " + id,
			Ctx:   "web_log.response_time_per_url",
			Type:  modules.Area,
			Dims: Dims{
				{ID: id + "_resp_time_min", Name: "min", Algo: modules.Incremental, Div: 1000},
				{ID: id + "_resp_time_max", Name: "max", Algo: modules.Incremental, Div: 1000},
				{ID: id + "_resp_time_avg", Name: "avg", Algo: modules.Incremental, Div: 1000},
			},
		},
	}
}

func (w *WebLog) createCharts() {
	var charts modules.Charts

	_ = charts.Add(responseStatuses.Copy(), responseCodes.Copy())

	if w.DoCodesDetailed {
		if w.DoCodesAggregate {
			_ = charts.Add(responseCodesDetailed.Copy())
		} else {
			_ = charts.Add(responseCodesDetailedPerFamily()...)
		}
	}

	if w.gm.has(keyBytesSent) || w.gm.has(keyRespLength) {
		_ = charts.Add(bandwidth.Copy())
	}

	if w.gm.has(keyRequest) && len(w.urlCats) > 0 {
		chart := requestsPerURL.Copy()
		_ = charts.Add(chart)
		for _, cat := range w.urlCats {
			_ = chart.AddDim(&Dim{
				ID:   cat.name,
				Algo: modules.Incremental,
			})
		}
	}

	if w.gm.has(keyRequest) && len(w.urlCats) > 0 && w.DoPerURLCharts {
		for _, cat := range w.urlCats {
			for _, chart := range perCategoryStats(cat.name) {
				_ = charts.Add(chart)
				for _, d := range chart.Dims {
					w.metrics[d.ID] = 0
				}
			}
		}
	}

	if w.gm.has(keyRequest) && len(w.userCats) > 0 {
		chart := requestsPerUserDefined.Copy()
		_ = charts.Add(chart)

		for _, cat := range w.userCats {
			_ = chart.AddDim(&Dim{
				ID:   cat.name,
				Algo: modules.Incremental,
			})
		}
	}

	if w.gm.has(keyRespTime) {
		_ = charts.Add(responseTime.Copy())
	}

	if w.gm.has(keyRespTime) && len(w.histograms) != 0 {
		chart := responseTimeHistogram.Copy()
		_ = charts.Add(chart)
		for _, v := range w.histograms[keyRespTimeHistogram] {
			_ = chart.AddDim(&Dim{
				ID:   v.id,
				Name: v.name,
				Algo: modules.Incremental,
			})
		}
	}

	if w.gm.has(keyRespTimeUpstream) {
		_ = charts.Add(responseTimeUpstream.Copy())
	}

	if w.gm.has(keyRespTimeUpstream) && len(w.histograms) != 0 {
		chart := responseTimeUpstreamHistogram.Copy()
		_ = charts.Add(chart)
		for _, v := range w.histograms[keyRespTimeUpstreamHistogram] {
			_ = chart.AddDim(&Dim{
				ID:   v.id,
				Name: v.name,
				Algo: modules.Incremental,
			})
		}
	}

	if w.gm.has(keyRequest) {
		_ = charts.Add(requestsPerHTTPMethod.Copy())
		_ = charts.Add(requestsPerHTTPVersion.Copy())
	}

	if w.gm.has(keyVhost) {
		_ = charts.Add(requestsPerVhost.Copy())
	}

	if w.gm.has(keyAddress) {
		_ = charts.Add(requestsPerIPProto.Copy())
		_ = charts.Add(currentPollIPs.Copy())
		if w.DoAllTimeIPs {
			_ = charts.Add(allTimeIPs.Copy())
		}
	}

	w.charts = &charts
}

func (w *WebLog) Charts() *Charts {
	return w.charts
}
