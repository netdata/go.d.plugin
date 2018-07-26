package springboot2

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
)

var uCharts = charts.Charts{
	{
		ID:   "heap",
		Opts: charts.Opts{Title: "Threads", Units: "threads", Family: "threads", Type: charts.Area},
		Dims: charts.Dims{
			{ID: "threads_daemon", Name: "daemon"},
			{ID: "threads", Name: "total"},
		},
	},
}
