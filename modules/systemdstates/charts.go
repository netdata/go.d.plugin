package systemdstates

import (
	"strings"

	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts
	// Dims is an alias for module.Dims
	Dims = module.Dims

	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "systemd_service_active_state",
		Title: "Systemd Service Active units",
		Units: "bool",
		Fam:   "service",
		Ctx:   "systemd.systemd_service_active_state",
		Dims:  Dims{},
	},
	{
		ID:    "systemd_socket_active_state",
		Title: "Systemd Socket Active units",
		Units: "bool",
		Fam:   "socket",
		Ctx:   "systemd.systemd_socket_active_state",
		Dims:  Dims{},
	},
	{
		ID:    "systemd_target_active_state",
		Title: "Systemd Target Active units",
		Units: "bool",
		Fam:   "target",
		Ctx:   "systemd.systemd_target_active_state",
		Dims:  Dims{},
	},
	{
		ID:    "systemd_path_active_state",
		Title: "Systemd path Active units",
		Units: "bool",
		Fam:   "path",
		Ctx:   "systemd.systemd_path_active_state",
		Dims:  Dims{},
	},
}

func (s SystemdStates) charts() *Charts {
	charts := charts.Copy()

	for _, chart := range *charts {

		dims, err := s.chartDims()
		if err != nil {
			s.Error(err)
		}
		for _, unit := range dims {
			if strings.Contains(unit.Name, "."+chart.Fam) {
				chart.Dims = append(chart.Dims, &Dim{ID: unit.Name})
			}

		}
	}
	return charts
}
