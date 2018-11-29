package springboot2

import (
	"github.com/netdata/go.d.plugin/modules"
)

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
	// Area is an alias for modules.Area
)

var charts = Charts{
	{
		ID:    "heap",
		Title: "Threads", Units: "threads", Fam: "threads", Type: modules.Area,
		Dims: Dims{
			{ID: "threads_daemon", Name: "daemon"},
			{ID: "threads", Name: "total"},
		},
	},
}
