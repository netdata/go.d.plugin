package systemdstates

import (
	"github.com/coreos/go-systemd/dbus"
)

type unit struct {
	dbus.UnitStatus
}

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

	result, err := getUnits(conn)
	if err != nil {
		return nil, err
	}

	mx := make(map[string]int64)

	filtered := s.filterUnits(result)
	for _, unit := range filtered {
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

func (s SystemdStates) chartDims() ([]unit, error) {
	conn, err := dbus.New()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	result, err := getUnits(conn)
	if err != nil {
		return nil, err
	}

	return s.filterUnits(result), nil

}

func (s SystemdStates) filterUnits(units []unit) []unit {

	filtered := make([]unit, 0, len(units))
	for _, unit := range units {

		if s.unitsMatcher.MatchString(unit.Name) && unit.LoadState == "loaded" {
			filtered = append(filtered, unit)
		}
	}

	return filtered
}

func getUnits(conn *dbus.Conn) ([]unit, error) {
	allUnits, err := conn.ListUnits()

	if err != nil {
		return nil, err
	}

	result := make([]unit, 0, len(allUnits))
	for _, status := range allUnits {

		unit := unit{
			UnitStatus: status,
		}
		result = append(result, unit)
	}

	return result, nil
}
