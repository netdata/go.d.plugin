package springboot2

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

var charts = raw.Charts{
	Order: raw.Order{"heap"},
	Definitions: raw.Definitions{
		{
			ID:      "heap",
			Options: raw.Options{"Threads", "threads", "threads", "", "area"},
			Dimensions: raw.Dimensions{
				{"threads_daemon", "daemon"},
				{"threads", "total"},
			},
		},
	},
}
