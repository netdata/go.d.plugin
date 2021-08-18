package traefik

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

var chartTmplEntrypointRequests = module.Chart{
	ID:    "entrypoint_requests_%s_%s",
	Title: "Processed HTTP requests for entrypoint <code>%s</code> protocol <code>%s</code>",
	Units: "requests/s",
	Fam:   "entrypoint requests",
	Ctx:   "traefik.entrypoint_requests",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: prefixEntrypointRequests + "1xx_%s_%s", Name: "1xx", Algo: module.Incremental},
		{ID: prefixEntrypointRequests + "2xx_%s_%s", Name: "2xx", Algo: module.Incremental},
		{ID: prefixEntrypointRequests + "3xx_%s_%s", Name: "3xx", Algo: module.Incremental},
		{ID: prefixEntrypointRequests + "4xx_%s_%s", Name: "4xx", Algo: module.Incremental},
		{ID: prefixEntrypointRequests + "5xx_%s_%s", Name: "5xx", Algo: module.Incremental},
	},
}

var chartTmplEntrypointRequestDuration = module.Chart{
	ID:    "entrypoint_request_duration_%s_%s",
	Title: "Processed HTTP Request average duration for entrypoint <code>%s</code> protocol <code>%s</code>",
	Units: "milliseconds",
	Fam:   "entrypoint request duration",
	Ctx:   "traefik.entrypoint_request_duration_average",
	Type:  module.Stacked,
	Dims: module.Dims{
		{ID: prefixEntrypointReqDurAvg + "1xx_%s_%s", Name: "1xx"},
		{ID: prefixEntrypointReqDurAvg + "2xx_%s_%s", Name: "2xx"},
		{ID: prefixEntrypointReqDurAvg + "3xx_%s_%s", Name: "3xx"},
		{ID: prefixEntrypointReqDurAvg + "4xx_%s_%s", Name: "4xx"},
		{ID: prefixEntrypointReqDurAvg + "5xx_%s_%s", Name: "5xx"},
	},
}

var chartTmplEntrypointOpenConnections = module.Chart{
	ID:    "entrypoint_open_connections_%s_%s",
	Title: "Open connections for entrypoint <code>%s</code> protocol <code>%s</code>",
	Units: "connections",
	Fam:   "entrypoint connections",
	Ctx:   "traefik.entrypoint_open_connections",
	Type:  module.Stacked,
}

func newChartEntrypointRequests(entrypoint, proto string) *module.Chart {
	return newEntrypointChart(chartTmplEntrypointRequests, entrypoint, proto)
}

func newChartEntrypointRequestDuration(entrypoint, proto string) *module.Chart {
	return newEntrypointChart(chartTmplEntrypointRequestDuration, entrypoint, proto)
}

func newChartEntrypointOpenConnections(entrypoint, proto string) *module.Chart {
	return newEntrypointChart(chartTmplEntrypointOpenConnections, entrypoint, proto)
}

func newEntrypointChart(tmpl module.Chart, entrypoint, proto string) *module.Chart {
	chart := tmpl.Copy()
	chart.ID = fmt.Sprintf(chart.ID, entrypoint, proto)
	chart.Title = fmt.Sprintf(chart.Title, entrypoint, proto)
	for _, d := range chart.Dims {
		d.ID = fmt.Sprintf(d.ID, entrypoint, proto)
	}
	return chart
}
