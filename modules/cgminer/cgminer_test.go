// SPDX-License-Identifier: GPL-3.0-or-later

// SPDX-License-Identifier: GPL-3.0-or-later

package cgminer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// We want to ensure that module is a reference type, nothing more.

	assert.IsType(t, (*Cgminer)(nil), New("http://localhost:4028"))
}

func TestCgminer_GetPoolInfo(t *testing.T) {
	cgminer := New("http://localhost:4028")

	poolInfo, err := cgminer.GetPoolInfo()
	if err != nil {
		t.Errorf("Error getting pool info: %s", err)
	}

	if len(poolInfo) == 0 {
		t.Errorf("No pool info returned")
	}

	for _, info := range poolInfo {
		if info.URL == "" {
			t.Errorf("Pool URL is empty")
		}
	}
}

func TestCgminer_GetDeviceInfo(t *testing.T) {
	cgminer := New("http://localhost:4028")

	deviceInfo, err := cgminer.GetDeviceInfo()
	if err != nil {
		t.Errorf("Error getting device info: %s", err)
	}

	if len(deviceInfo) == 0 {
		t.Errorf("No device info returned")
	}

	for _, info := range deviceInfo {
		if info.Name == "" {
			t.Errorf("Device name is empty")
		}
	}
}

func TestCgminer_GetMinerConfig(t *testing.T) {
	cgminer := New("http://localhost:4028")

	config, err := cgminer.GetMinerConfig()
	if err != nil {
		t.Errorf("Error getting miner config: %s", err)
	}

	if len(config) == 0 {
		t.Errorf("No config returned")
	}
}
