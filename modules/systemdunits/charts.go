// +build linux

package systemdunits

import (
	"github.com/netdata/go.d.plugin/agent/module"
)

// systemd unit types: https://www.freedesktop.org/software/systemd/man/systemd.html
var charts = module.Charts{
	{
		ID:    "service_unit_state",
		Title: "Service Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "service",
		Ctx:   "systemd.service_units_state",
	},
	{
		ID:    "socket_unit_state",
		Title: "Socket Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "socket",
		Ctx:   "systemd.socket_unit_state",
	},
	{
		ID:    "target_unit_state",
		Title: "Target Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "target",
		Ctx:   "systemd.target_unit_state",
	},
	{
		ID:    "path_unit_state",
		Title: "Path Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "path",
		Ctx:   "systemd.path_unit_state",
	},
	{
		ID:    "device_unit_state",
		Title: "Device Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "device",
		Ctx:   "systemd.device_unit_state",
	},
	{
		ID:    "mount_unit_state",
		Title: "Mount Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "mount",
		Ctx:   "systemd.mount_unit_state",
	},
	{
		ID:    "automount_unit_state",
		Title: "Automount Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "automount",
		Ctx:   "systemd.automount_unit_state",
	},
	{
		ID:    "swap_unit_state",
		Title: "Swap Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "swap",
		Ctx:   "systemd.swap_unit_state",
	},
	{
		ID:    "timer_unit_state",
		Title: "Timer Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "timer",
		Ctx:   "systemd.timer_unit_state",
	},
	{
		ID:    "scope_unit_state",
		Title: "Scope Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "scope",
		Ctx:   "systemd.scope_unit_state",
	},
	{
		ID:    "slice_unit_state",
		Title: "Slice Unit State (1: active, 2: inactive, 3: activating, 4: deactivating, 5: failed)",
		Units: "state",
		Fam:   "slice",
		Ctx:   "systemd.slice_unit_state",
	},
}
