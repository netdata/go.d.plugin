package snmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// We want to ensure that module is a reference type, nothing more.

	assert.IsType(t, (*SNMP)(nil), New())
}

func TestSNMP_Init(t *testing.T) {
	// 'Init() bool' initializes the module with an appropriate config, so to test it we need:
	// - provide the config.
	// - set module.Config field with the config.
	// - call Init() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success on default config": {
			config: New().Config,
		},
		"success when only 'charts' set": {
			config: Config{
				Charts: ConfigCharts{
					Dims: 2,
				},
			},
		},
		"success when 'charts' and 'hidden_charts' set": {
			config: Config{
				Charts: ConfigCharts{
					Dims: 2,
				},
			},
		},
		"fails when 'charts' and 'hidden_charts' set, but 'num' == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Dims: 2,
				},
			},
		},
		"fails when only 'charts' set, 'num' > 0, but 'dimensions' == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Dims: 0,
				},
			},
		},
		"fails when only 'hidden_charts' set, 'num' > 0, but 'dimensions' == 0": {
			wantFail: true,
			config:   Config{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			example := New()
			example.Config = test.config

			if test.wantFail {
				assert.False(t, example.Init())
			} else {
				assert.True(t, example.Init())
			}
		})
	}
}

func TestSNMP_Check(t *testing.T) {
	// 'Check() bool' reports whether the module is able to collect any data, so to test it we need:
	// - provide the module with a specific config.
	// - initialize the module (call Init()).
	// - call Check() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare  func() *SNMP
		wantFail bool
	}{
		"success on default":                            {prepare: prepareSNMPDefault},
		"success when only 'charts' set":                {prepare: prepareSNMPOnlyCharts},
		"success when only 'hidden_charts' set":         {prepare: prepareSNMPOnlyHiddenCharts},
		"success when 'charts' and 'hidden_charts' set": {prepare: prepareSNMPChartsAndHiddenCharts},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			example := test.prepare()
			require.True(t, example.Init())

			if test.wantFail {
				assert.False(t, example.Check())
			} else {
				assert.True(t, example.Check())
			}
		})
	}
}

func TestSNMP_Charts(t *testing.T) {
	// We want to ensure that initialized module does not return 'nil'.
	// If it is not 'nil' we are ok.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare func(t *testing.T) *SNMP
		wantNil bool
	}{
		"not initialized collector": {
			wantNil: true,
			prepare: func(t *testing.T) *SNMP {
				return New()
			},
		},
		"initialized collector": {
			prepare: func(t *testing.T) *SNMP {
				example := New()
				require.True(t, example.Init())
				return example
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			example := test.prepare(t)

			if test.wantNil {
				assert.Nil(t, example.Charts())
			} else {
				assert.NotNil(t, example.Charts())
			}
		})
	}
}

func TestSNMP_Cleanup(t *testing.T) {
	// Since this module has nothing to clean up,
	// we want just to ensure that Cleanup() not panics.

	assert.NotPanics(t, New().Cleanup)
}

func TestSNMP_Collect(t *testing.T) {
	// 'Collect() map[string]int64' returns collected data, so to test it we need:
	// - provide the module with a specific config.
	// - initialize the module (call Init()).
	// - call Collect() and compare its return value with the expected value.

	// 'test' map contains different test cases.
	tests := map[string]struct {
		prepare       func() *SNMP
		wantCollected map[string]int64
	}{
		"default config": {
			prepare: prepareSNMPDefault,
			wantCollected: map[string]int64{
				"random_0_random0": 1,
				"random_0_random1": -1,
				"random_0_random2": 1,
				"random_0_random3": -1,
			},
		},
		"only 'charts' set": {
			prepare: prepareSNMPOnlyCharts,
			wantCollected: map[string]int64{
				"random_0_random0": 1,
				"random_0_random1": -1,
				"random_0_random2": 1,
				"random_0_random3": -1,
				"random_0_random4": 1,
				"random_1_random0": 1,
				"random_1_random1": -1,
				"random_1_random2": 1,
				"random_1_random3": -1,
				"random_1_random4": 1,
			},
		},
		"only 'hidden_charts' set": {
			prepare: prepareSNMPOnlyHiddenCharts,
			wantCollected: map[string]int64{
				"hidden_random_0_random0": 1,
				"hidden_random_0_random1": -1,
				"hidden_random_0_random2": 1,
				"hidden_random_0_random3": -1,
				"hidden_random_0_random4": 1,
				"hidden_random_1_random0": 1,
				"hidden_random_1_random1": -1,
				"hidden_random_1_random2": 1,
				"hidden_random_1_random3": -1,
				"hidden_random_1_random4": 1,
			},
		},
		"'charts' and 'hidden_charts' set": {
			prepare: prepareSNMPChartsAndHiddenCharts,
			wantCollected: map[string]int64{
				"hidden_random_0_random0": 1,
				"hidden_random_0_random1": -1,
				"hidden_random_0_random2": 1,
				"hidden_random_0_random3": -1,
				"hidden_random_0_random4": 1,
				"hidden_random_1_random0": 1,
				"hidden_random_1_random1": -1,
				"hidden_random_1_random2": 1,
				"hidden_random_1_random3": -1,
				"hidden_random_1_random4": 1,
				"random_0_random0":        1,
				"random_0_random1":        -1,
				"random_0_random2":        1,
				"random_0_random3":        -1,
				"random_0_random4":        1,
				"random_1_random0":        1,
				"random_1_random1":        -1,
				"random_1_random2":        1,
				"random_1_random3":        -1,
				"random_1_random4":        1,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			example := test.prepare()
			require.True(t, example.Init())

			collected := example.Collect()

			assert.Equal(t, test.wantCollected, collected)
			ensureCollectedHasAllChartsDimsVarsIDs(t, example, collected)
		})
	}
}

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, e *SNMP, collected map[string]int64) {
	for _, chart := range *e.Charts() {
		if chart.Obsolete {
			continue
		}
		for _, dim := range chart.Dims {
			_, ok := collected[dim.ID]
			assert.Truef(t, ok,
				"collected metrics has no data for dim '%s' chart '%s'", dim.ID, chart.ID)
		}
		for _, v := range chart.Vars {
			_, ok := collected[v.ID]
			assert.Truef(t, ok,
				"collected metrics has no data for var '%s' chart '%s'", v.ID, chart.ID)
		}
	}
}

func prepareSNMPDefault() *SNMP {
	return prepareSNMP(New().Config)
}

func prepareSNMPOnlyCharts() *SNMP {
	return prepareSNMP(Config{
		Charts: ConfigCharts{
			Dims: 5,
		},
	})
}

func prepareSNMPOnlyHiddenCharts() *SNMP {
	return prepareSNMP(Config{})
}

func prepareSNMPChartsAndHiddenCharts() *SNMP {
	return prepareSNMP(Config{
		Charts: ConfigCharts{
			Dims: 5,
		},
	})
}

func prepareSNMP(cfg Config) *SNMP {
	example := New()
	example.Config = cfg
	example.randInt = func() int64 { return 1 }
	return example
}
