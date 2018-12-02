package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockModule_Init(t *testing.T) {
	m := &mockModule{}

	// no default Init
	require.Panics(t, func() { m.Init() })

	m.init = func() bool { return true }
	assert.True(t, m.Init())
}

func TestMockModule_Check(t *testing.T) {
	m := &mockModule{}

	// no default Check
	require.Panics(t, func() { m.Check() })

	m.check = func() bool { return true }
	assert.True(t, m.Check())
}

func TestMockModule_Charts(t *testing.T) {
	m := &mockModule{}

	// no default Charts
	require.Panics(t, func() { m.Charts() })

	m.charts = func() *Charts { return nil }
	assert.Nil(t, m.Charts())
}

func TestMockModule_GetData(t *testing.T) {
	m := &mockModule{}

	// no default GatherMetrics
	require.Panics(t, func() { m.GatherMetrics() })

	m.gatherMetrics = func() map[string]int64 { return nil }
	assert.Nil(t, m.GatherMetrics())
}

func TestMockModule_Cleanup(t *testing.T) {
	m := &mockModule{}
	assert.False(t, m.cleanupDone)

	m.Cleanup()
	assert.True(t, m.cleanupDone)
}
