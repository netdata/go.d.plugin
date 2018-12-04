package godplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockJob_FullName(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.FullName() })
	// OK case
	m.fullName = func() string { return "name" }
	assert.Equal(t, "name", m.FullName())
}

func TestMockJob_ModuleName(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.ModuleName() })
	// OK case
	m.moduleName = func() string { return "name" }
	assert.Equal(t, "name", m.ModuleName())
}

func TestMockJob_Name(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Name() })
	// OK case
	m.name = func() string { return "name" }
	assert.Equal(t, "name", m.Name())
}

func TestMockJob_AutoDetectionRetry(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.AutoDetectionRetry() })
	// OK case
	m.autoDetectionRetry = func() int { return 1 }
	assert.Equal(t, 1, m.AutoDetectionRetry())
}

func TestMockJob_Panicked(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Panicked() })
	// OK case
	m.panicked = func() bool { return true }
	assert.True(t, m.Panicked())
}

func TestMockJob_Init(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Init() })
	// OK case
	m.init = func() bool { return true }
	assert.True(t, m.Init())
}

func TestMockJob_Check(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Check() })
	// OK case
	m.check = func() bool { return true }
	assert.True(t, m.Check())
}

func TestMockJob_PostCheck(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.PostCheck() })
	// OK case
	m.postCheck = func() bool { return true }
	assert.True(t, m.PostCheck())
}

func TestMockJob_Tick(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Tick(1) })
	// OK case
	m.tick = func(int) {}
	assert.NotPanics(t, func() { m.Tick(1) })
}

func TestMockJob_Start(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Start() })
	// OK case
	m.start = func() {}
	assert.NotPanics(t, func() { m.Start() })
}

func TestMockJob_Stop(t *testing.T) {
	m := &mockJob{}

	// PANIC case
	assert.Panics(t, func() { m.Stop() })
	// OK case
	m.stop = func() {}
	assert.NotPanics(t, func() { m.Stop() })
}
