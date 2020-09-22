package systemdunits

import (
	"testing"
	"regexp"


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

}


func TestSystemdUnits_Init(t *testing.T) {
	job := New()

	assert.True(t, job.Init())
	job.Init()

}

func TestSystemdUnits_FilterUnits(t *testing.T) {

	job := New()

	fixtures := getMockUnit()
	job.selector = regexp.MustCompile("^foo$")

	filtered := job.filterUnits(fixtures[0])
	for _, unit := range filtered {
		if !job.selector.MatchString(unit.Name) {
			t.Error(unit.Name, "should not be in the filtered list")
		}
	}

	if len(filtered) != len(fixtures[0])-1 {
		t.Error("Default filters removed units")
	}
}

func TestSystemdUnits_extractUnitType(t *testing.T) {

	rightUnit, err := extractUnitType("nginx.service")
	assert.Nil(t, err)
	assert.NotNil(t, rightUnit)

	nginxUnit, err := extractUnitType("nginx.service")
	assert.Nil(t, err)
	assert.Equal(t, nginxUnit,"service")


	wrongUnit, err := extractUnitType("nginx.wrong")
	assert.NotNil(t, err)
	assert.Equal(t, wrongUnit,"")

}

func TestSystemdUnits_isUnitTypeValid(t *testing.T) {

	serviceUnit := isUnitTypeValid("service")
	assert.Equal(t, serviceUnit, true)

	mountUnit := isUnitTypeValid("mount")
	assert.Equal(t, mountUnit, true)

	wrongUnit := isUnitTypeValid("wrong")
	assert.Equal(t, wrongUnit,false)

}

func TestSystemdUnits_convertUnitState(t *testing.T) {

	var state int64
	state = convertUnitState("active")
	assert.Equal(t, state, int64(1))

	state = convertUnitState("failed")
	assert.Equal(t, state, int64(3))

	state = convertUnitState("inactive")
	assert.Equal(t, state, int64(4))

	state = convertUnitState("wrong")
	assert.Equal(t, state, int64(-1))

}
