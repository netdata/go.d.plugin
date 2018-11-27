package modules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockModule_Init(t *testing.T) {
	mock := &MockModule{}
	f := func() {
		mock.Init()
	}

	require.Panics(t, f)

	mock.InitFunc = func() bool { return true }
	assert.True(t, mock.Init())
}

func TestMockModule_Check(t *testing.T) {
	mock := &MockModule{}
	f := func() {
		mock.Check()
	}

	require.Panics(t, f)

	mock.CheckFunc = func() bool { return true }
	assert.True(t, mock.Check())
}

func TestMockModule_GetCharts(t *testing.T) {
	mock := &MockModule{}
	f := func() {
		mock.GetCharts()
	}

	require.Panics(t, f)

	mock.GetChartsFunc = func() *Charts { return nil }
	assert.Nil(t, mock.GetCharts())
}

func TestMockModule_GetData(t *testing.T) {
	mock := &MockModule{}
	f := func() {
		mock.GetData()
	}

	require.Panics(t, f)

	mock.GetDataDunc = func() map[string]int64 { return nil }
	assert.Nil(t, mock.GetData())
}

func TestMockModule_Cleanup(t *testing.T) {
	mock := &MockModule{}
	assert.False(t, mock.CleanupDone)

	mock.Cleanup()
	assert.True(t, mock.CleanupDone)
}
