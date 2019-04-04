package coredns

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Chart is an alias for module.Chart
	Chart = module.Chart
	// Dims is an alias for module.Dims
	Dims = module.Dims
	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var summaryCharts = Charts{
	{
		ID:    "requests_rate",
		Title: "Requests Rate",
		Units: "events/s",
		Fam:   "summary",
		Ctx:   "coredns.requests_rate",
		Dims: Dims{
			{ID: "request_total", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "responses_rate",
		Title: "Responses Rate",
		Units: "events/s",
		Fam:   "summary",
		Ctx:   "coredns.responses_rate",
		Dims: Dims{
			{ID: "response_total", Name: "responses", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_rate_by_status",
		Title: "Requests Rate By Status",
		Units: "requests/ss",
		Fam:   "summary",
		Ctx:   "coredns.requests_rate_by_status",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "request_by_status_processed", Name: "processed", Algo: module.Incremental},
			{ID: "request_by_status_dropped", Name: "dropped", Algo: module.Incremental},
		},
	},
	{
		ID:    "panics_rate",
		Title: "Panics Rate",
		Units: "events/s",
		Fam:   "summary",
		Ctx:   "coredns.panics_rate",
		Dims: Dims{
			{ID: "panic_total", Name: "panics", Algo: module.Incremental},
		},
	},
}

var serverCharts = Charts{
	{
		ID:    "%s_requests_rate",
		Title: "Requests Rate",
		Units: "events/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate",
		Dims: Dims{
			{ID: "%s_request_total", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_responses_rate",
		Title: "Responses Rate",
		Units: "events/s",
		Fam:   "server %s",
		Ctx:   "coredns.responses_rate",
		Dims: Dims{
			{ID: "%s_response_total", Name: "responses", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_status",
		Title: "Requests Rate By Status",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_status",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_status_processed", Name: "processed", Algo: module.Incremental},
			{ID: "%s_request_by_status_dropped", Name: "dropped", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_proto",
		Title: "Requests Rate By Proto",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_proto",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_proto_udp", Name: "udp", Algo: module.Incremental},
			{ID: "%s_request_by_proto_tcp", Name: "tcp", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_ip_family",
		Title: "Requests Rate By IP Family",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_ip_family",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_ip_family_v4", Name: "v4", Algo: module.Incremental},
			{ID: "%s_request_by_ip_family_v6", Name: "v6", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_requests_rate_by_type",
		Title: "Requests Rate By Type",
		Units: "requests/s",
		Fam:   "server %s",
		Ctx:   "coredns.requests_rate_by_type",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_request_by_type_A", Name: "A", Algo: module.Incremental},
			{ID: "%s_request_by_type_AAAA", Name: "AAAA", Algo: module.Incremental},
			{ID: "%s_request_by_type_MX", Name: "MX", Algo: module.Incremental},
			{ID: "%s_request_by_type_SOA", Name: "SOA", Algo: module.Incremental},
			{ID: "%s_request_by_type_CNAME", Name: "CNAME", Algo: module.Incremental},
			{ID: "%s_request_by_type_PTR", Name: "PTR", Algo: module.Incremental},
			{ID: "%s_request_by_type_TXT", Name: "TXT", Algo: module.Incremental},
			{ID: "%s_request_by_type_NS", Name: "NS", Algo: module.Incremental},
			{ID: "%s_request_by_type_DS", Name: "DS", Algo: module.Incremental},
			{ID: "%s_request_by_type_DNSKEY", Name: "DNSKEY", Algo: module.Incremental},
			{ID: "%s_request_by_type_RRSIG", Name: "RRSIG", Algo: module.Incremental},
			{ID: "%s_request_by_type_NSEC", Name: "NSEC", Algo: module.Incremental},
			{ID: "%s_request_by_type_NSEC3", Name: "NSEC3", Algo: module.Incremental},
			{ID: "%s_request_by_type_IXFR", Name: "IXFR", Algo: module.Incremental},
			{ID: "%s_request_by_type_ANY", Name: "ANY", Algo: module.Incremental},
			{ID: "%s_request_by_type_other", Name: "other", Algo: module.Incremental},
		},
	},
	{
		ID:    "%s_responses_rate_by_type",
		Title: "Responses Rate By Type",
		Units: "responses/s",
		Fam:   "server %s",
		Ctx:   "coredns.responses_rate_by_type",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "%s_response_by_rcode_NOERROR", Name: "NOERROR", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_FORMERR", Name: "FORMERR", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_SERVFAIL", Name: "SERVFAIL", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_NXDOMAIN", Name: "NXDOMAIN", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_NOTIMP", Name: "NOTIMP", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_REFUSED", Name: "REFUSED", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_YXDOMAIN", Name: "YXDOMAIN", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_YXRRSET", Name: "YXRRSET", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_NXRRSET", Name: "NXRRSET", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_NOTAUTH", Name: "NOTAUTH", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_NOTZONE", Name: "NOTZONE", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADSIG", Name: "BADSIG", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADKEY", Name: "BADKEY", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADTIME", Name: "BADTIME", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADMODE", Name: "BADMODE", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADNAME", Name: "BADNAME", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADALG", Name: "BADALG", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADTRUNC", Name: "BADTRUNC", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_BADCOOKIE", Name: "BADCOOKIE", Algo: module.Incremental},
			{ID: "%s_response_by_rcode_other", Name: "other", Algo: module.Incremental},
		},
	},
}
