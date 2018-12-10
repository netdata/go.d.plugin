package nginxvts

import "github.com/netdata/go.d.plugin/modules"

type (
	// Charts is an alias for modules.Charts
	Charts = modules.Charts
	// Dims is an alias for modules.Dims
	Dims = modules.Dims
)

var charts = Charts{
	{
		ID:    "shared_mem",
		Title: "Shared Memory Usage", Units: "B", Fam: "shared_zones",
		Dims: Dims{
			{ID: "shared_used_size", Name: "used"},
			{ID: "shared_max_size", Name: "max"},
		},
	},
}
