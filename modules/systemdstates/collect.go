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

		chartID := fmt.Sprintf("%s_states", ut)
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

		if unit.LoadState == "loaded" && s.selector.MatchString(unit.Name) {
			units[i] = unit
			i++
		}
	}

	return units[:i]

}

func extractUnitType(unit string) (string, error) {

	idx := strings.LastIndexByte(unit, '.')

	if idx <= 0 {
		return "", fmt.Errorf("Could not find a type for : %v", unit)
	}
	ut := unit[idx+1:]
	if !isUnitTypeValid(ut) {
		return "", fmt.Errorf("Could not find a valid type for : %v", unit)
	}

	return ut, nil
}

func isUnitTypeValid(unit string) bool {
	switch unit {
	case "service", "socket", "device", "mount", "automount", "swap", "target", "path", "timer", "scope", "slice":
		return true
	}
	return false
}
