package nginxvts

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

const (
	serverReqChart = "server_req"
)

var charts = Charts{
	{
		ID:    "shared_mem",
		Title: "Shared Memory Usage", Units: "B", Ctx: "nginx-vts", Fam: "shared_zones",
		Dims: Dims{
			{ID: "shared_used_size", Name: "used"},
			{ID: "shared_max_size", Name: "max"},
		},
	},
	{
		ID:    serverReqChart,
		Title: "Requests per Server", Units: "requests/s", Ctx: "nginx-vts", Fam: "server_zones", Type: modules.Stacked,
		Dims: Dims{},
	},
}

func createZoneCharts(prefix string, name string) Charts {
	return Charts{
		{
			ID:    prefix + "_bandwidth_" + name,
			Title: "Bandwidth", Units: "kilobits/s", Ctx: "nginx-vts", Fam: prefix + "_zones." + name, Type: modules.Area,
			Dims: Dims{
				{ID: prefix + "_" + name + "_received", Name: "received", Algo: modules.Incremental, Mul: 8, Div: 1000},
				{ID: prefix + "_" + name + "_sent", Name: "sent", Algo: modules.Incremental, Mul: -8, Div: 1000},
			},
		},
		{
			ID:    prefix + "_response_" + name,
			Title: "Response Codes", Units: "requests/s", Ctx: "nginx-vts", Fam: prefix + "_zones." + name, Type: modules.Stacked,
			Dims: Dims{
				{ID: prefix + "_" + name + "_resp_1xx", Name: "1xx", Algo: modules.Incremental},
				{ID: prefix + "_" + name + "_resp_2xx", Name: "2xx", Algo: modules.Incremental},
				{ID: prefix + "_" + name + "_resp_3xx", Name: "3xx", Algo: modules.Incremental},
				{ID: prefix + "_" + name + "_resp_4xx", Name: "4xx", Algo: modules.Incremental},
				{ID: prefix + "_" + name + "_resp_5xx", Name: "5xx", Algo: modules.Incremental},
			},
		},
	}
}

func createZoneReqDim(dimID string, name string) *modules.Dim {
	return &modules.Dim{ID: dimID, Name: name, Algo: modules.Incremental}
}
