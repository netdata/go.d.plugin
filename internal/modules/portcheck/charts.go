package portcheck

import (
	"fmt"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/raw"
)

type (
	Charts      = raw.Charts
	Order       = raw.Order
	Definitions = raw.Definitions
	Chart       = raw.Chart
	Options     = raw.Options
	Dimensions  = raw.Dimensions
	Dimension   = raw.Dimension
)

func charts(port int) *Charts {
	family := sprintf("number %d", port)
	return &raw.Charts{
		Order: Order{
			sprintf("status_%d", port),
			sprintf("instate_%d", port),
			sprintf("latency_%d", port)},
		Definitions: Definitions{
			&Chart{
				ID:      sprintf("status_%d", port),
				Options: Options{"Port Check Status", "boolean", family, "portcheck.status"},
				Dimensions: Dimensions{
					Dimension{sprintf("success_%d", port), "success"},
					Dimension{sprintf("failed_%d", port), "failed"},
					Dimension{sprintf("timeout_%d", port), "timeout"},
				},
			},
			&Chart{
				ID:      sprintf("instate_%d", port),
				Options: Options{"Current State Duration", "seconds", family, "portcheck.instate"},
				Dimensions: Dimensions{
					Dimension{sprintf("instate_%d", port), "time"},
				},
			},
			&Chart{
				ID:      sprintf("latency_%d", port),
				Options: Options{"TCP Connect Latency", "ms", family, "portcheck.latency"},
				Dimensions: Dimensions{
					Dimension{sprintf("latency_%d", port), "time", "", 1, 1e6},
				},
			},
		},
	}
}

func sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}
