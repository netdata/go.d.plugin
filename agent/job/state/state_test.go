package state

import (
	"fmt"
	"testing"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"

	"github.com/stretchr/testify/assert"
)

// TODO: tech debt
func TestNewManager(t *testing.T) {

}

// TODO: tech debt
func TestManager_Run(t *testing.T) {

}

// TODO: tech debt
func TestManager_Save(t *testing.T) {

}

// TODO: tech debt
func TestManager_Remove(t *testing.T) {

}

// TODO: tech debt
func TestState_Contains(t *testing.T) {

}

// TODO: tech debt
func TestLoad(t *testing.T) {

}

func TestStore_add(t *testing.T) {
	tests := map[string]struct {
		prepare   func() *Store
		input     confgroup.Config
		wantItems map[string]map[string]string
	}{
		"add an item to the empty store": {
			prepare: func() *Store {
				return &Store{}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
			),
			wantItems: map[string]map[string]string{
				"modName": {
					"jobName:18299273693089411682": "state",
				},
			},
		},
		"add an item with same module, same name, but specific options": {
			prepare: func() *Store {
				return &Store{
					items: map[string]map[string]string{
						"modName": {
							"jobName:18299273693089411682": "state",
						},
					},
				}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
				"opt", "val",
			),
			wantItems: map[string]map[string]string{
				"modName": {
					"jobName:18299273693089411682": "state",
					"jobName:6762067169527372123":  "state",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := test.prepare()
			s.add(test.input, "state")
			fmt.Println(s.items)
			assert.Equal(t, test.wantItems, s.items)
		})
	}
}

func prepareConfig(values ...string) confgroup.Config {
	cfg := confgroup.Config{}
	for i := 1; i < len(values); i += 2 {
		cfg[values[i-1]] = values[i]
	}
	return cfg
}
