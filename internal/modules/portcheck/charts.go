package portcheck

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

type (
	Charts     = charts.Charts
	Chart      = charts.Chart
	Options    = charts.Options
	Dimensions = charts.Dimensions
)

func uCharts(port int) Charts {
	family := sprintf("number %d", port)
	return Charts{
		{
			ID: sprintf("status_%d", port),
			Options: Options{
				Title: "Port Check Status", Units: "boolean", Family: family, Context: "portcheck.status"},
			Dimensions: Dimensions{
				{ID: sprintf("success_%d", port), Name: "success"},
				{ID: sprintf("failed_%d", port), Name: "failed"},
				{ID: sprintf("timeout_%d", port), Name: "timeout"},
			},
		},
		{
			ID: sprintf("instate_%d", port),
			Options: Options{
				Title: "Current State Duration", Units: "seconds", Family: family, Context: "portcheck.instate"},
			Dimensions: Dimensions{
				{ID: sprintf("instate_%d", port), Name: "time"},
			},
		},
		{
			ID:      sprintf("latency_%d", port),
			Options: Options{Title: "TCP Connect Latency", Units: "ms", Family: family, Context: "portcheck.latency"},
			Dimensions: Dimensions{
				{ID: sprintf("latency_%d", port), Name: "time", Divisor: 1000000},
			},
		},
	}
}

func sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}
