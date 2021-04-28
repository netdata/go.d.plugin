package unity

import "github.com/netdata/go-orchestrator/module"

type (
	Charts = module.Charts
	Chart = module.Chart
	Dims = module.Dims
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "unity_metrics",
		Title: "Unity Dell EMC Metrics",
		Units: "u",
		Fam:   "server",
		Ctx:   "unity.metrics",
		Dims: Dims{
			{ID: "path", Name:"hello world"},
		},
	},
}

func getUnit(path string) string{
	return "u"
}