package state

import (
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
		prepare      func() *Store
		input        confgroup.Config
		wantItemsNum int
	}{
		"add a cfg to the empty store": {
			prepare: func() *Store {
				return &Store{}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
			),
			wantItemsNum: 1,
		},
		"add a cfg that already in the store": {
			prepare: func() *Store {
				return &Store{
					items: map[string]map[string]string{
						"modName": {"jobName:18299273693089411682": "state"},
					},
				}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
			),
			wantItemsNum: 1,
		},
		"add a cfg with same module, same name, but specific options": {
			prepare: func() *Store {
				return &Store{
					items: map[string]map[string]string{
						"modName": {"jobName:18299273693089411682": "state"},
					},
				}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
				"opt", "val",
			),
			wantItemsNum: 2,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := test.prepare()
			s.add(test.input, "state")
			assert.Equal(t, test.wantItemsNum, calcItemsNum(s))
		})
	}
}

func TestStore_remove(t *testing.T) {
	tests := map[string]struct {
		prepare      func() *Store
		input        confgroup.Config
		wantItemsNum int
	}{
		"remove a cfg from the empty store": {
			prepare: func() *Store {
				return &Store{}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
			),
			wantItemsNum: 0,
		},
		"remove a cfg from the store": {
			prepare: func() *Store {
				return &Store{
					items: map[string]map[string]string{
						"modName": {
							"jobName:18299273693089411682": "state",
							"jobName:18299273693089411683": "state",
						},
					},
				}
			},
			input: prepareConfig(
				"module", "modName",
				"name", "jobName",
			),
			wantItemsNum: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := test.prepare()
			s.remove(test.input)
			assert.Equal(t, test.wantItemsNum, calcItemsNum(s))
		})
	}
}

func calcItemsNum(s *Store) (num int) {
	for _, v := range s.items {
		for range v {
			num += 1
		}
	}
	return num
}

func prepareConfig(values ...string) confgroup.Config {
	cfg := confgroup.Config{}
	for i := 1; i < len(values); i += 2 {
		cfg[values[i-1]] = values[i]
	}
	return cfg
}
