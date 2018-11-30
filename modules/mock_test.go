package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockModule_Init(t *testing.T) {
	mock := &mockModule{}
	f := func() {
		mock.Init()
	}

	require.Panics(t, f)

	mock.initFunc = func() bool { return true }
	assert.True(t, mock.Init())
}

func TestMockModule_Check(t *testing.T) {
	mock := &mockModule{}
	f := func() {
		mock.Check()
	}

	require.Panics(t, f)

	mock.checkFunc = func() bool { return true }
	assert.True(t, mock.Check())
}

func TestMockModule_Charts(t *testing.T) {
	mock := &mockModule{}
	f := func() {
		mock.Charts()
	}

	require.Panics(t, f)

	mock.chartsFunc = func() *Charts { return nil }
	assert.Nil(t, mock.Charts())
}

func TestMockModule_GetData(t *testing.T) {
	mock := &mockModule{}
	f := func() {
		mock.GatherMetrics()
	}

	require.Panics(t, f)

	mock.gatherMetricsFunc = func() map[string]int64 { return nil }
	assert.Nil(t, mock.GatherMetrics())
}

func TestMockModule_Cleanup(t *testing.T) {
	mock := &mockModule{}
	assert.False(t, mock.cleanupDone)

	mock.Cleanup()
	assert.True(t, mock.cleanupDone)
}
