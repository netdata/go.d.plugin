package systemdunits

import "github.com/coreos/go-systemd/v22/dbus"

type systemdDBusClient struct{}

func (systemdDBusClient) connect() (systemdConnection, error) {
	// TODO: switch when available
	// New => NewWithContext
	// ListUnits => ListUnitsContext
	return dbus.New()
}

func newSystemdDBusClient() *systemdDBusClient {
	return &systemdDBusClient{}
}
