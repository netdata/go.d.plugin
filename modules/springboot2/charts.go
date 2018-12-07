package springboot2

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "response",
		Title: "Response Codes", Units: "requests/s", Fam: "response_code", Type: modules.Stacked,
		Dims: Dims{
			{ID: "resp_1xx", Name: "1xx", Algo: modules.Incremental},
			{ID: "resp_2xx", Name: "2xx", Algo: modules.Incremental},
			{ID: "resp_3xx", Name: "3xx", Algo: modules.Incremental},
			{ID: "resp_4xx", Name: "4xx", Algo: modules.Incremental},
			{ID: "resp_5xx", Name: "5xx", Algo: modules.Incremental},
		},
	},
	{
		ID:    "thread",
		Title: "Threads", Units: "threads", Fam: "threads", Type: modules.Area,
		Dims: Dims{
			{ID: "threads_daemon", Name: "daemon"},
			{ID: "threads", Name: "total"},
		},
	},
}
