package systemdstates

import (
	"fmt"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
)

const (
	inactive = 0
	active   = 1
	failed   = 2
)

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

		chartID := fmt.Sprintf("systemd_%s_active_state", s.unitType(unit.Name))
		chart := s.charts.Get(chartID)
		if !chart.HasDim(unit.Name) {
			chart.AddDim(&Dim{ID: unit.Name})
		}

		state := -1
		if unit.ActiveState == "active" {
			state = active
		}
		if unit.ActiveState == "inactive" {
			state = inactive
		}
		if unit.ActiveState == "failed" {
			state = failed
		}
		mx[unit.Name] = int64(state)
	}

	return mx, nil
}

func (s SystemdStates) filterUnits(units []dbus.UnitStatus) []dbus.UnitStatus {

	filtered := make([]dbus.UnitStatus, 0, len(units))
	for _, unit := range units {

		if s.unitsMatcher.MatchString(unit.Name) && unit.LoadState == "loaded" {
			filtered = append(filtered, unit)
		}
	}

	return filtered
}

func (s SystemdStates) unitType(unit string) string {
	validTypes := []string{"service", "socket", "device", "mount", "automount", "swap", "target", "path", "timer", "scope"}
	var ut string
	for _, t := range validTypes {
		if strings.HasSuffix(unit, fmt.Sprintf(".%s", t)) {
			ut = t
			break
		}
	}
	return ut

}
