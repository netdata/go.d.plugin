package systemdunits

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/coreos/go-systemd/v22/dbus"
)

func (s *SystemdUnits) collect() (map[string]int64, error) {
	units, err := s.getLoadedUnits()
	if err != nil {
		return nil, err
	}

	if len(units) == 0 {
		return nil, nil
	}

	mx := make(map[string]int64)
	for _, unit := range units {
		name := cleanUnitName(unit.Name)

		typ, err := extractUnitType(name)
		if err != nil {
			continue
		}

		if !s.collectedUnits[name] {
			s.collectedUnits[name] = true
			s.addUnitToCharts(name, typ)
		}
		mx[name] = convertUnitState(unit.ActiveState)
	}
	return mx, nil
}

func (s *SystemdUnits) getLoadedUnits() ([]dbus.UnitStatus, error) {
	if s.conn == nil {
		conn, err := s.client.connect()
		if err != nil {
			return nil, err
		}
		s.conn = conn
	}

	units, err := s.conn.ListUnitsByPatterns(
		[]string{"active", "activating", "failed", "inactive", "deactivating"},
		s.Include,
	)
	if err != nil {
		s.conn.Close()
		s.conn = nil
		return nil, err
	}

	loaded := units[:0]
	for _, unit := range units {
		if unit.LoadState == "loaded" {
			loaded = append(loaded, unit)
		}
	}
	s.Debugf("got total/loaded %d/%d units", len(units), len(loaded))
	return loaded, nil
}

func (s *SystemdUnits) addUnitToCharts(name, typ string) {
	id := fmt.Sprintf("%s_unit_state", typ)
	chart := s.Charts().Get(id)
	if chart == nil {
		s.Warningf("add dimension (unit '%s'): can't find '%s' chart", name, id)
		return
	}

	dim := &module.Dim{
		ID:   name,
		Name: name[:len(name)-len(typ)-1], // name.type => name
	}
	if err := chart.AddDim(dim); err != nil {
		s.Warningf("add dimension (unit '%s'): %v", name, err)
	}
	chart.MarkNotCreated()
}

func extractUnitType(name string) (string, error) {
	// name.type => type
	idx := strings.LastIndexByte(name, '.')
	if idx <= 0 {
		return "", fmt.Errorf("could not find a type for: %s", name)
	}

	typ := name[idx+1:]
	if !isUnitTypeValid(typ) {
		return "", fmt.Errorf("could not find a valid type for: %s", name)
	}
	return typ, nil
}

func isUnitTypeValid(typ string) bool {
	// https://www.freedesktop.org/software/systemd/man/systemd.html
	switch typ {
	case "service",
		"socket",
		"device",
		"mount",
		"automount",
		"swap",
		"target",
		"path",
		"timer",
		"scope",
		"slice":
		return true
	}
	return false
}

func convertUnitState(state string) int64 {
	// https://www.freedesktop.org/software/systemd/man/systemd.html
	switch state {
	case "active":
		return 1
	case "inactive":
		return 2
	case "activating":
		return 3
	case "deactivating":
		return 4
	case "failed":
		return 5
	default:
		return -1
	}
}

func cleanUnitName(name string) string {
	// dev-disk-by\x2duuid-DE44\x2dCEE0.device => dev-disk-by-uuid-DE44-CEE0.device
	if strings.IndexByte(name, '\\') == -1 {
		return name
	}
	v, err := strconv.Unquote("\"" + name + "\"")
	if err != nil {
		return name
	}
	return v
}
