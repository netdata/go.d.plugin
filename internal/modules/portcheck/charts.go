package portcheck

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	Charts     = charts.Charts
	Chart      = charts.Chart
	Options    = charts.Opts
	Dimensions = charts.Dims
)

func uCharts(port int) Charts {
	family := sprintf("number %d", port)
	return Charts{
		{
			ID: sprintf("status_%d", port),
			Opts: Options{
				Title: "Port Check Status", Units: "boolean", Fam: family, Ctx: "portcheck.status"},
			Dims: Dimensions{
				{ID: sprintf("success_%d", port), Name: "success"},
				{ID: sprintf("failed_%d", port), Name: "failed"},
				{ID: sprintf("timeout_%d", port), Name: "timeout"},
			},
		},
		{
			ID: sprintf("instate_%d", port),
			Opts: Options{
				Title: "Current State Duration", Units: "seconds", Fam: family, Ctx: "portcheck.instate"},
			Dims: Dimensions{
				{ID: sprintf("instate_%d", port), Name: "time"},
			},
		},
		{
			ID:   sprintf("latency_%d", port),
			Opts: Options{Title: "TCP Connect Latency", Units: "ms", Fam: family, Ctx: "portcheck.latency"},
			Dims: Dimensions{
				{ID: sprintf("latency_%d", port), Name: "time", Div: 1000000},
			},
		},
	}
}

func sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}
