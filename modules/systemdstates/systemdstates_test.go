package systemdstates

import (
	"regexp"
	"testing"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/netdata/go-orchestrator/module"
	"github.com/stretchr/testify/assert"
)

// Creates mock UnitLists
func getMockUnit() [][]dbus.UnitStatus {
	unit1 := []dbus.UnitStatus{
		{

			Name:        "foo",
			Description: "foo desc",
			LoadState:   "loaded",
			ActiveState: "active",
			SubState:    "running",
			Followed:    "",
			Path:        "/org/freedesktop/systemd1/unit/foo",
			JobId:       0,
			JobType:     "",
			JobPath:     "/",
		},
		{

			Name:        "bar",
			Description: "bar desc",
			LoadState:   "not-found",
			ActiveState: "inactive",
			SubState:    "dead",
			Followed:    "",
			Path:        "/org/freedesktop/systemd1/unit/bar",
			JobId:       0,
			JobType:     "",
			JobPath:     "/",
		},
	}

	return [][]dbus.UnitStatus{unit1}
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Len(t, job.metrics, len(job.metrics))
}

func TestSystemdStates_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
	job.Init()

}

func TestSystemdStates_FilterUnits(t *testing.T) {

	job := New()

	fixtures := getMockUnit()
	job.unitsMatcher = regexp.MustCompile("^foo$")

	filtered := job.filterUnits(fixtures[0])
	for _, unit := range filtered {
		if !job.unitsMatcher.MatchString(unit.Name) {
			t.Error(unit.Name, "should not be in the filtered list")
		}
	}

	if len(filtered) != len(fixtures[0])-1 {
		t.Error("Default filters removed units")
	}
}

func TestSystemdStates_unitType(t *testing.T) {

	job := New()

	rightUnit, err := job.unitType("nginx.service")
	assert.Nil(t, err)
	assert.NotNil(t, rightUnit)

	nginxUnit, err := job.unitType("nginx.service")
	assert.Nil(t, err)
	assert.Equal(t, nginxUnit, "service")

	wrongUnit, err := job.unitType("nginx.wrong")
	assert.NotNil(t, err)
	assert.Equal(t, wrongUnit, "")

}
