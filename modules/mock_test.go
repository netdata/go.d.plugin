package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockModule_Init(t *testing.T) {
	m := &MockModule{}

	// no default Init
	require.Panics(t, func() { m.Init() })

	m.InitFunc = func() bool { return true }
	assert.True(t, m.Init())
}

func TestMockModule_Check(t *testing.T) {
	m := &MockModule{}

	// no default Check
	require.Panics(t, func() { m.Check() })

	m.CheckFunc = func() bool { return true }
	assert.True(t, m.Check())
}

func TestMockModule_Charts(t *testing.T) {
	m := &MockModule{}

	// no default Charts
	require.Panics(t, func() { m.Charts() })

	m.ChartsFunc = func() *Charts { return nil }
	assert.Nil(t, m.Charts())
}

func TestMockModule_GetData(t *testing.T) {
	m := &MockModule{}

	// no default GatherMetrics
	require.Panics(t, func() { m.GatherMetrics() })

	m.GatherMetricsFunc = func() map[string]int64 { return nil }
	assert.Nil(t, m.GatherMetrics())
}

func TestMockModule_Cleanup(t *testing.T) {
	m := &MockModule{}
	assert.False(t, m.CleanupDone)

	m.Cleanup()
	assert.True(t, m.CleanupDone)
}
