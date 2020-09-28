package systemdunits

import (
	"fmt"
	"strings"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	job := New()
	assert.Implements(t, (*module.Module)(nil), job)
}

func TestSystemdUnits_Init(t *testing.T) {
	job := New()
	job.Selector.Includes = []string{"* *.service"}
	assert.True(t, job.Init())
}

func TestSystemdUnits_Charts(t *testing.T) {
	job := New()
	for _, chart := range *job.Charts() {
		idx := strings.IndexByte(chart.ID, '_')
		unit := chart.ID[:idx]
		chartID := fmt.Sprintf("%s_states", unit)
		assert.Equal(t, chart.ID, chartID)
		assert.Equal(t, chart.Fam, unit)
	}
}

func TestSystemdUnits_Collect(t *testing.T) {
	tests := map[string]struct {
		name  string
		state int64
		units []dbus.UnitStatus
	}{
		"service1": {
			"service1.service",
			1, // 1: active
			[]dbus.UnitStatus{
				{
					Name:        "service1.service",
					Description: "service1 desc",
					LoadState:   "loaded",
					ActiveState: "active",
					SubState:    "running",
					Followed:    "",
					Path:        "/org/freedesktop/systemd1/unit/service1",
					JobId:       0,
					JobType:     "",
					JobPath:     "/",
				},
			},
		},
		"service2": {
			"service2.service",
			2, // 2: activating
			[]dbus.UnitStatus{
				{
					Name:        "service2.service",
					Description: "service2 desc",
					LoadState:   "loaded",
					ActiveState: "activating",
					SubState:    "running",
					Followed:    "",
					Path:        "/org/freedesktop/systemd1/unit/service1",
					JobId:       0,
					JobType:     "",
					JobPath:     "/",
				},
			},
		},
		"service3": {
			"service3.service",
			3, // 3: failed
			[]dbus.UnitStatus{
				{
					Name:        "service3.service",
					Description: "service3 desc",
					LoadState:   "loaded",
					ActiveState: "failed",
					SubState:    "running",
					Followed:    "",
					Path:        "/org/freedesktop/systemd1/unit/service1",
					JobId:       0,
					JobType:     "",
					JobPath:     "/",
				},
			},
		},
		"service4": {
			"service4.service",
			4, // 4: inactive
			[]dbus.UnitStatus{
				{
					Name:        "service4.service",
					Description: "service4 desc",
					LoadState:   "loaded",
					ActiveState: "inactive",
					SubState:    "running",
					Followed:    "",
					Path:        "/org/freedesktop/systemd1/unit/service1",
					JobId:       0,
					JobType:     "",
					JobPath:     "/",
				},
			},
		},
		"service5": {
			"service5.service",
			5, // 5: deactivating
			[]dbus.UnitStatus{
				{
					Name:        "service5.service",
					Description: "service5 desc",
					LoadState:   "loaded",
					ActiveState: "deactivating",
					SubState:    "running",
					Followed:    "",
					Path:        "/org/freedesktop/systemd1/unit/service1",
					JobId:       0,
					JobType:     "",
					JobPath:     "/",
				},
			},
		},
	}
	job := New()
	job.Selector.Includes = []string{"* *.service"}
	require.True(t, job.Init())

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			require.True(t, job.Init())
			job.units = test.units
			collected := job.Collect()
			assert.Equal(t, collected[test.name], test.state)
		})
	}

	job2 := New()
	job2.Selector.Includes = []string{"* *.target"}
	require.True(t, job2.Init())
	require.True(t, job2.Init())
	job2.units = tests["service1"].units
	assert.Nil(t, job2.Collect())
}
