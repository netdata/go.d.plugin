package portcheck

import (
	"fmt"

	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims
)

func chartsTemplate(port int) Charts {
	fam := fmt.Sprintf("port %d", port)

	return Charts{
		{
			ID:    fmt.Sprintf("status_%d", port),
			Title: "Port Check Status", Units: "boolean", Fam: fam, Ctx: "portcheck.status",
			Dims: Dims{
				{
					ID:   fmt.Sprintf("success_%d", port),
					Name: "success",
				},
				{
					ID:   fmt.Sprintf("failed_%d", port),
					Name: "failed",
				},
				{
					ID:   fmt.Sprintf("timeout_%d", port),
					Name: "timeout",
				},
			},
		},
		{
			ID:    fmt.Sprintf("instate_%d", port),
			Title: "Current State Duration", Units: "seconds", Fam: fam, Ctx: "portcheck.instate",
			Dims: Dims{
				{
					ID:   fmt.Sprintf("instate_%d", port),
					Name: "time",
				},
			},
		},
		{
			ID:    fmt.Sprintf("latency_%d", port),
			Title: "TCP Connect Latency", Units: "ms", Fam: fam, Ctx: "portcheck.latency",
			Dims: Dims{
				{
					ID:   fmt.Sprintf("latency_%d", port),
					Name: "time",
					Div:  1000000,
				},
			},
		},
	}
}
