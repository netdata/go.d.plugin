package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockJob_FullName(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockFullName, m.FullName())
}

func TestMockJob_ModuleName(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockModuleName, m.ModuleName())
}

func TestMockJob_Name(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockName, m.Name())
}

func TestMockJob_AutoDetectionRetry(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockAutoDetectionRetry, m.AutoDetectionRetry())
}

func TestMockJob_Panicked(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockPanicked, m.Panicked())
}

func TestMockJob_Init(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockInit, m.Init())
}

func TestMockJob_Check(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockCheck, m.Check())
}

func TestMockJob_PostCheck(t *testing.T) {
	m := &mockJob{}

	assert.Equal(t, mockPostCheck, m.PostCheck())
}

func TestMockJob_Tick(t *testing.T) {
	m := &mockJob{}

	assert.NotPanics(t, func() { m.Tick(1) })
}

func TestMockJob_Start(t *testing.T) {
	m := &mockJob{}

	assert.NotPanics(t, func() { m.Start() })
}

func TestMockJob_Stop(t *testing.T) {
	m := &mockJob{}

	assert.NotPanics(t, func() { m.Stop() })
}
