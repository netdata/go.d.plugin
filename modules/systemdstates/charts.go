package systemdstates

import (
	"github.com/netdata/go-orchestrator/module"
)

type (
	// Charts is an alias for module.Charts
	Charts = module.Charts

	// Dim is an alias for module.Dim
	Dim = module.Dim
)

var charts = Charts{
	{
		ID:    "systemd_service_active_state",
		Title: "Systemd Service Active units",
		Units: "state",
		Fam:   "service",
		Ctx:   "systemd.systemd_service_active_state",
	},
	{
		ID:    "systemd_socket_active_state",
		Title: "Systemd Socket Active units",
		Units: "state",
		Fam:   "socket",
		Ctx:   "systemd.systemd_socket_active_state",
	},
	{
		ID:    "systemd_target_active_state",
		Title: "Systemd Target Active units",
		Units: "state",
		Fam:   "target",
		Ctx:   "systemd.systemd_target_active_state",
	},
	{
		ID:    "systemd_path_active_state",
		Title: "Systemd path Active units",
		Units: "state",
		Fam:   "path",
		Ctx:   "systemd.systemd_path_active_state",
	},
	{
		ID:    "systemd_device_active_state",
		Title: "Systemd device Active units",
		Units: "state",
		Fam:   "device",
		Ctx:   "systemd.systemd_device_active_state",
	},
	{
		ID:    "systemd_mount_active_state",
		Title: "Systemd mount Active units",
		Units: "state",
		Fam:   "mount",
		Ctx:   "systemd.systemd_mount_active_state",
	},
	{
		ID:    "systemd_automount_active_state",
		Title: "Systemd automount Active units",
		Units: "state",
		Fam:   "automount",
		Ctx:   "systemd.systemd_automount_active_state",
	},
	{
		ID:    "systemd_swap_active_state",
		Title: "Systemd swap Active units",
		Units: "state",
		Fam:   "swap",
		Ctx:   "systemd.systemd_swap_active_state",
	},
	{
		ID:    "systemd_timer_active_state",
		Title: "Systemd timer Active units",
		Units: "state",
		Fam:   "timer",
		Ctx:   "systemd.systemd_timer_active_state",
	},
	{
		ID:    "systemd_scope_active_state",
		Title: "Systemd scope Active units",
		Units: "state",
		Fam:   "scope",
		Ctx:   "systemd.systemd_scope_active_state",
	},
}
