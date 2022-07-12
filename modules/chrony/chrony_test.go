// SPDX-License-Identifier: GPL-3.0-or-later

package chrony

import (
	"errors"
	"testing"

	"github.com/netdata/go.d.plugin/modules/chrony/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChrony_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default config": {
			config: New().Config,
		},
		"unset 'address'": {
			wantFail: true,
			config: Config{
				Address: "",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			chrony := New()
			chrony.Config = test.config

			if test.wantFail {
				assert.False(t, chrony.Init())
			} else {
				assert.True(t, chrony.Init())
			}
		})
	}
}

func TestChrony_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *Chrony
		wantFail bool
	}{
		"tracking: success, activity: success": {
			wantFail: false,
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{}) },
		},
		"tracking: success, activity: fail": {
			wantFail: false,
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{errOnActivity: true}) },
		},
		"tracking: fail, activity: success": {
			wantFail: true,
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{errOnTracking: true}) },
		},
		"tracking: fail, activity: fail": {
			wantFail: true,
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{errOnTracking: true}) },
		},
		"fail on creating client": {
			wantFail: true,
			prepare:  func() *Chrony { return prepareChronyWithMock(nil) },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			chrony := test.prepare()

			require.True(t, chrony.Init())

			if test.wantFail {
				assert.False(t, chrony.Check())
			} else {
				assert.True(t, chrony.Check())
			}
		})
	}
}

func TestChrony_Charts(t *testing.T) {
	assert.Equal(t, len(charts), len(*New().Charts()))
}

func TestChrony_Cleanup(t *testing.T) {
	tests := map[string]struct {
		prepare   func(c *Chrony)
		wantClose bool
	}{
		"after New": {
			wantClose: false,
			prepare:   func(c *Chrony) { return },
		},
		"after Init": {
			wantClose: false,
			prepare:   func(c *Chrony) { c.Init() },
		},
		"after Check": {
			wantClose: true,
			prepare:   func(c *Chrony) { c.Init(); c.Check() },
		},
		"after Collect": {
			wantClose: true,
			prepare:   func(c *Chrony) { c.Init(); c.Collect() },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &mockClient{}
			chrony := prepareChronyWithMock(m)
			test.prepare(chrony)

			require.NotPanics(t, chrony.Cleanup)

			if test.wantClose {
				assert.True(t, m.closeCalled)
			} else {
				assert.False(t, m.closeCalled)
			}
		})
	}
}

func TestChrony_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *Chrony
		expected map[string]int64
	}{
		"tracking: success, activity: success": {
			prepare: func() *Chrony { return prepareChronyWithMock(&mockClient{}) },
			expected: map[string]int64{
				"burst_offline_sources":      3,
				"burst_online_sources":       4,
				"current_correction":         111249,
				"frequency":                  51036781311,
				"last_offset":                -88888,
				"leap_status_delete_second":  0,
				"leap_status_insert_second":  1,
				"leap_status_normal":         0,
				"leap_status_unsynchronised": 0,
				"offline_sources":            2,
				"online_sources":             8,
				"ref_timestamp":              1667,
				"rms_offset":                 359872,
				"root_delay":                 51769230,
				"root_dispersion":            1243559,
				"skew":                       67318372,
				"stratum":                    3,
				"unresolved_sources":         1,
				"update_interval":            1038400390625,
			},
		},
		"tracking: success, activity: fail": {
			prepare: func() *Chrony { return prepareChronyWithMock(&mockClient{errOnActivity: true}) },
			expected: map[string]int64{
				"current_correction":         111249,
				"frequency":                  51036781311,
				"last_offset":                -88888,
				"leap_status_delete_second":  0,
				"leap_status_insert_second":  1,
				"leap_status_normal":         0,
				"leap_status_unsynchronised": 0,
				"ref_timestamp":              1667,
				"rms_offset":                 359872,
				"root_delay":                 51769230,
				"root_dispersion":            1243559,
				"skew":                       67318372,
				"stratum":                    3,
				"update_interval":            1038400390625,
			},
		},
		"tracking: fail, activity: success": {
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{errOnTracking: true}) },
			expected: nil,
		},
		"tracking: fail, activity: fail": {
			prepare:  func() *Chrony { return prepareChronyWithMock(&mockClient{errOnTracking: true}) },
			expected: nil,
		},
		"fail on creating client": {
			prepare:  func() *Chrony { return prepareChronyWithMock(nil) },
			expected: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			chrony := test.prepare()

			require.True(t, chrony.Init())
			_ = chrony.Check()

			collected := chrony.Collect()
			copyRefTimestamp(collected, test.expected)

			assert.Equal(t, test.expected, collected)
		})
	}
}

func prepareChronyWithMock(m *mockClient) *Chrony {
	c := New()
	if m == nil {
		c.newClient = func(c *Chrony) (chronyClient, error) { return nil, errors.New("mock.newClient error") }
	} else {
		c.newClient = func(c *Chrony) (chronyClient, error) { return m, nil }
	}
	return c
}

type mockClient struct {
	errOnTracking bool
	errOnActivity bool
	closeCalled   bool
}

func (m mockClient) Tracking() (*client.TrackingPayload, error) {
	if m.errOnTracking {
		return nil, errors.New("mockClient.Tracking call error")
	}
	tp := client.TrackingPayload{
		RefID: 1540987708,
		Ip: client.IPAddr{
			IPAddrHigh: 6618491809397997568,
			IPAddrLow:  0,
			Family:     1,
			Pad:        0,
		},
		Stratum:    3,
		LeapStatus: 1,
		RefTime: client.ChronyTimespec{
			TvSecHigh: 0,
			TvSecLow:  1657633575,
			TvNSec:    895532067,
		},
		CurrentCorrection:  -387363189,
		LastOffset:         -381315542,
		RmsOffset:          -323179191,
		FreqPpm:            255056470,
		ResidFreqPpm:       -215937554,
		SkewPpm:            -58073545,
		RootDelay:          -86766599,
		RootDispersion:     -257753360,
		LastUpdateInterval: 411159760,
	}
	return &tp, nil
}

func (m mockClient) Activity() (*client.ActivityPayload, error) {
	if m.errOnActivity {
		return nil, errors.New("mockClient.Activity call error")
	}
	ap := client.ActivityPayload{
		Online:       8,
		Offline:      2,
		BurstOnline:  4,
		BurstOffline: 3,
		Unresolved:   1,
	}
	return &ap, nil
}

func (m *mockClient) Close() {
	m.closeCalled = true
}

func copyRefTimestamp(dst, src map[string]int64) {
	if _, ok := dst["ref_timestamp"]; !ok {
		return
	}
	if _, ok := src["ref_timestamp"]; !ok {
		return
	}
	dst["ref_timestamp"] = src["ref_timestamp"]
}
