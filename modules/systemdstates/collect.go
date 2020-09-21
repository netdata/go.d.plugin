package systemdstates

import (
	"fmt"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
)

func convertUnitState(state string) int64 {

	switch state {
	case "active":
		return 1
	case "activating":
		return 2
	case "failed":
		return 3
	case "inactive":
		return 4
	case "deactivating":
		return 5
	default:
		return -1
	}

}

func (s *SystemdStates) collect() (map[string]int64, error) {

	var err error
	conn, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	allUnits, err := conn.ListUnits()
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)

	units := s.filterUnits(allUnits)
	for _, unit := range units {

		ut, err := extractUnitType(unit.Name)
		if err != nil {
			return nil, err
		}

		chartID := fmt.Sprintf("systemd_%s_active_state", ut)
		chart := s.charts.Get(chartID)
		if !chart.HasDim(unit.Name) {
			_ = chart.AddDim(&Dim{ID: unit.Name})
		}

		mx[unit.Name] = convertUnitState(unit.ActiveState)
	}

	return mx, nil
}

func (s SystemdStates) filterUnits(units []dbus.UnitStatus) []dbus.UnitStatus {

	var i int
	for _, unit := range units {

		if unit.LoadState == "loaded" && s.unitsMatcher.MatchString(unit.Name) {
			units[i] = unit
			i++
		}
	}

	return units[:i]

}

func extractUnitType(unit string) (string, error) {
	ut := ""
	idx := strings.LastIndexByte(unit, '.')

	switch suffix := unit[idx:]; suffix {
	case ".service":
		ut = "service"
	case ".socket":
		ut = "socket"
	case ".device":
		ut = "device"
	case ".mount":
		ut = "mount"
	case ".automount":
		ut = "automount"
	case ".swap":
		ut = "swap"
	case ".target":
		ut = "target"
	case ".path":
		ut = "path"
	case ".timer":
		ut = "timer"
	case ".scope":
		ut = "scope"
	}

	if ut == "" || idx == -1 {
		return "", fmt.Errorf("Could not find a type for : %v", unit)
	}
	return ut, nil
}
