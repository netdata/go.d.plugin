package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func (d Dimension) validateConfig() error {
	if d.Name == "" || len(d.Name) > 15 {
		return fmt.Errorf("invalid/missing dimension name")
	}
	if d.OID == "" {
		return fmt.Errorf("missing OID value")
	}
	if d.Algorithm == "" &&
		d.Algorithm != string(module.Incremental) &&
		d.Algorithm != string(module.PercentOfIncremental) &&
		d.Algorithm != string(module.PercentOfAbsolute) {
		return fmt.Errorf("missing/invalid algorithm")
	}
	return nil
}

func (u User) validateConfig() error {
	return nil
}

func (c ChartsConfig) validateConfig() error {
	return nil
}
