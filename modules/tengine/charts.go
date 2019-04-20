package tengine

import "github.com/netdata/go-orchestrator/module"

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

var charts = Charts{
	{
		ID:    "bandwidth",
		Title: "Bandwidth",
		Units: "B/s",
		Fam:   "bandwidth",
		Type:  module.Area,
		Dims: Dims{
			{ID: "bytes_in", Name: "in", Algo: module.Incremental},
			{ID: "bytes_out", Name: "out", Algo: module.Incremental, Mul: -1},
		},
	},
	{
		ID:    "connections",
		Title: "Connections",
		Units: "connections/s",
		Fam:   "connections",
		Dims: Dims{
			{ID: "conn_total", Name: "accepted", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_total",
		Title: "Requests",
		Units: "requests/s",
		Fam:   "requests",
		Dims: Dims{
			{ID: "req_total", Name: "processed", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_per_response_code_family",
		Title: "Requests Per Response Code Family",
		Units: "requests/s",
		Fam:   "requests",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "http_2xx", Name: "2xx", Algo: module.Incremental},
			{ID: "http_5xx", Name: "5xx", Algo: module.Incremental},
			{ID: "http_3xx", Name: "3xx", Algo: module.Incremental},
			{ID: "http_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "http_other_status", Name: "other", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_per_response_code_detailed",
		Title: "Requests Per Response Code Detailed",
		Units: "requests/s",
		Fam:   "requests",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "http_200", Name: "200", Algo: module.Incremental},
			{ID: "http_206", Name: "206", Algo: module.Incremental},
			{ID: "http_302", Name: "302", Algo: module.Incremental},
			{ID: "http_304", Name: "304", Algo: module.Incremental},
			{ID: "http_403", Name: "403", Algo: module.Incremental},
			{ID: "http_404", Name: "404", Algo: module.Incremental},
			{ID: "http_416", Name: "419", Algo: module.Incremental},
			{ID: "http_499", Name: "499", Algo: module.Incremental},
			{ID: "http_500", Name: "500", Algo: module.Incremental},
			{ID: "http_502", Name: "502", Algo: module.Incremental},
			{ID: "http_503", Name: "503", Algo: module.Incremental},
			{ID: "http_504", Name: "504", Algo: module.Incremental},
			{ID: "http_508", Name: "508", Algo: module.Incremental},
			{ID: "http_other_detail_status", Name: "other", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_upstream",
		Title: "Number Of Requests Calling For Upstream",
		Units: "requests/s",
		Fam:   "upstream",
		Dims: Dims{
			{ID: "ups_req", Name: "requests", Algo: module.Incremental},
		},
	},
	{
		ID:    "tries_upstream",
		Title: "Number Of Times Calling For Upstream",
		Units: "calls/s",
		Fam:   "upstream",
		Dims: Dims{
			{ID: "ups_tries", Name: "calls", Algo: module.Incremental},
		},
	},
	{
		ID:    "requests_upstream_per_response_code_family",
		Title: "Upstream Requests Per Response Code Family",
		Units: "requests/s",
		Fam:   "upstream",
		Type:  module.Stacked,
		Dims: Dims{
			{ID: "http_ups_4xx", Name: "4xx", Algo: module.Incremental},
			{ID: "http_ups_5xx", Name: "5xx", Algo: module.Incremental},
		},
	},
}
