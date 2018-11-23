package tcpcheck

import (
	"fmt"

	"github.com/netdata/go.d.plugin/modules"
)

type (
	Charts = modules.Charts
	Dims   = modules.Dims
)

func chartsTemplate(port int) Charts {
	family := sprintf("port %d", port)
	return Charts{
		{
			ID:    sprintf("status_%d", port),
			Title: "Port Check Status", Units: "boolean", Fam: family, Ctx: "tcpcheck.status",
			Dims: Dims{
				{ID: sprintf("success_%d", port), Name: "success"},
				{ID: sprintf("failed_%d", port), Name: "failed"},
				{ID: sprintf("timeout_%d", port), Name: "timeout"},
			},
		},
		{
			ID:    sprintf("instate_%d", port),
			Title: "Current State Duration", Units: "seconds", Fam: family, Ctx: "tcpcheck.instate",
			Dims: Dims{
				{ID: sprintf("instate_%d", port), Name: "time"},
			},
		},
		{
			ID:    sprintf("latency_%d", port),
			Title: "TCP Connect Latency", Units: "ms", Fam: family, Ctx: "tcpcheck.latency",
			Dims: Dims{
				{ID: sprintf("latency_%d", port), Name: "time", Div: 1000000},
			},
		},
	}
}

func sprintf(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}
