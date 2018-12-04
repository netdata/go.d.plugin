package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockJob_FullName(t *testing.T) {
	m := &mockJob{}

	expected := "name"

	assert.NotEqual(t, expected, m.FullName())
	m.fullName = func() string { return expected }
	assert.Equal(t, expected, m.FullName())
}

func TestMockJob_ModuleName(t *testing.T) {
	m := &mockJob{}

	expected := "name"

	assert.NotEqual(t, expected, m.ModuleName())
	m.moduleName = func() string { return expected }
	assert.Equal(t, expected, m.ModuleName())
}

func TestMockJob_Name(t *testing.T) {
	m := &mockJob{}

	expected := "name"

	assert.NotEqual(t, expected, m.Name())
	m.name = func() string { return expected }
	assert.Equal(t, expected, m.Name())
}

func TestMockJob_AutoDetectionRetry(t *testing.T) {
	m := &mockJob{}

	expected := -1

	assert.NotEqual(t, expected, m.AutoDetectionRetry())
	m.autoDetectionRetry = func() int { return expected }
	assert.Equal(t, expected, m.AutoDetectionRetry())
}

func TestMockJob_Panicked(t *testing.T) {
	m := &mockJob{}

	assert.False(t, m.Panicked())
	m.panicked = func() bool { return true }
	assert.True(t, m.Panicked())
}

func TestMockJob_Init(t *testing.T) {
	m := &mockJob{}

	assert.True(t, m.Init())
	m.init = func() bool { return false }
	assert.False(t, m.Init())
}

func TestMockJob_Check(t *testing.T) {
	m := &mockJob{}

	assert.True(t, m.Check())
	m.check = func() bool { return false }
	assert.False(t, m.Check())
}

func TestMockJob_PostCheck(t *testing.T) {
	m := &mockJob{}

	assert.True(t, m.PostCheck())
	m.postCheck = func() bool { return false }
	assert.False(t, m.PostCheck())
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
