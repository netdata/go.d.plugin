package systemdunits

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
		ID:    "service_states",
		Title: "Systemd Service Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "service",
		Ctx:   "systemd.service_states",
	},
	{
		ID:    "socket_states",
		Title: "Systemd Socket Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "socket",
		Ctx:   "systemd.socket_states",
	},
	{
		ID:    "target_states",
		Title: "Systemd Target Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "target",
		Ctx:   "systemd.target_states",
	},
	{
		ID:    "path_states",
		Title: "Systemd Path Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "path",
		Ctx:   "systemd.path_states",
	},
	{
		ID:    "device_states",
		Title: "Systemd Device Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "device",
		Ctx:   "systemd.device_states",
	},
	{
		ID:    "mount_states",
		Title: "Systemd Mount Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "mount",
		Ctx:   "systemd.mount_states",
	},
	{
		ID:    "automount_states",
		Title: "Systemd Automount Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "automount",
		Ctx:   "systemd.automount_states",
	},
	{
		ID:    "swap_states",
		Title: "Systemd Swap Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "swap",
		Ctx:   "systemd.swap_states",
	},
	{
		ID:    "timer_states",
		Title: "Systemd Timer Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "timer",
		Ctx:   "systemd.timer_states",
	},
	{
		ID:    "scope_states",
		Title: "Systemd Scope Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "scope",
		Ctx:   "systemd.scope_states",
	},
	{
		ID:    "slice_states",
		Title: "Systemd Slice Units (1: active, 2: activating, 3: failed, 4: inactive, 5: deactivating)",
		Units: "state",
		Fam:   "slice",
		Ctx:   "systemd.slice_states",
	},
}
