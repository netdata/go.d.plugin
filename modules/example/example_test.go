package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.IsType(t, (*Example)(nil), New())
}

func TestExample_Init(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"default": {
			config: New().Config,
		},
		"only charts": {
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"only hidden charts": {
			config: Config{
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"charts and hidden charts": {
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 2,
				},
			},
		},
		"charts->num and hidden_charts->num == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Num:  0,
					Dims: 2,
				},
				HiddenCharts: ConfigCharts{
					Num:  0,
					Dims: 2,
				},
			},
		},
		"charts->num > 0 and charts->dimensions == 0": {
			wantFail: true,
			config: Config{
				Charts: ConfigCharts{
					Num:  1,
					Dims: 0,
				},
			},
		},
		"hidden_charts->num > 0 and hidden_charts->dimensions == 0": {
			wantFail: true,
			config: Config{
				HiddenCharts: ConfigCharts{
					Num:  1,
					Dims: 0,
				},
			},
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

func TestExample_Check(t *testing.T) {
	tests := map[string]struct {
		prepare  func() *Example
		wantFail bool
	}{
		"default":                  {prepare: prepareExampleDefault},
		"only charts":              {prepare: prepareExampleOnlyCharts},
		"only hidden charts":       {prepare: prepareExampleOnlyHiddenCharts},
		"charts and hidden charts": {prepare: prepareExampleChartsAndHiddenCharts},
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

func TestExample_Charts(t *testing.T) {
	tests := map[string]struct {
		prepare func(t *testing.T) *Example
		wantNil bool
	}{
		"not initialized collector": {
			wantNil: true,
			prepare: func(t *testing.T) *Example {
				return New()
			},
		},
		"initialized collector": {
			prepare: func(t *testing.T) *Example {
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

func TestExample_Cleanup(t *testing.T) {
	assert.NotPanics(t, New().Cleanup)
}

func TestExample_Collect(t *testing.T) {
	tests := map[string]struct {
		prepare       func() *Example
		wantCollected map[string]int64
	}{
		"default": {
			prepare: prepareExampleDefault,
			wantCollected: map[string]int64{
				"random_0_random0": 1,
				"random_0_random1": -1,
				"random_0_random2": 1,
				"random_0_random3": -1,
			},
		},
		"only charts": {
			prepare: prepareExampleOnlyCharts,
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
		"only hidden charts": {
			prepare: prepareExampleOnlyHiddenCharts,
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
		"chart and hidden charts": {
			prepare: prepareExampleChartsAndHiddenCharts,
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

func ensureCollectedHasAllChartsDimsVarsIDs(t *testing.T, e *Example, collected map[string]int64) {
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

func prepareExampleDefault() *Example {
	return prepareExample(New().Config)
}

func prepareExampleOnlyCharts() *Example {
	return prepareExample(Config{
		Charts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareExampleOnlyHiddenCharts() *Example {
	return prepareExample(Config{
		HiddenCharts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareExampleChartsAndHiddenCharts() *Example {
	return prepareExample(Config{
		Charts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
		HiddenCharts: ConfigCharts{
			Num:  2,
			Dims: 5,
		},
	})
}

func prepareExample(cfg Config) *Example {
	example := New()
	example.Config = cfg
	example.randInt = func() int64 { return 1 }
	return example
}
