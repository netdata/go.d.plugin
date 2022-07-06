// SPDX-License-Identifier: GPL-3.0-or-later

package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockJob_FullName(t *testing.T) {
	m := &MockJob{}
	expected := "name"

	assert.NotEqual(t, expected, m.FullName())
	m.FullNameFunc = func() string { return expected }
	assert.Equal(t, expected, m.FullName())
}

func TestMockJob_ModuleName(t *testing.T) {
	m := &MockJob{}
	expected := "name"

	assert.NotEqual(t, expected, m.ModuleName())
	m.ModuleNameFunc = func() string { return expected }
	assert.Equal(t, expected, m.ModuleName())
}

func TestMockJob_Name(t *testing.T) {
	m := &MockJob{}
	expected := "name"

	assert.NotEqual(t, expected, m.Name())
	m.NameFunc = func() string { return expected }
	assert.Equal(t, expected, m.Name())
}

func TestMockJob_AutoDetectionEvery(t *testing.T) {
	m := &MockJob{}
	expected := -1

	assert.NotEqual(t, expected, m.AutoDetectionEvery())
	m.AutoDetectionEveryFunc = func() int { return expected }
	assert.Equal(t, expected, m.AutoDetectionEvery())
}

func TestMockJob_RetryAutoDetection(t *testing.T) {
	m := &MockJob{}
	expected := true

	assert.True(t, m.RetryAutoDetection())
	m.RetryAutoDetectionFunc = func() bool { return expected }
	assert.True(t, m.RetryAutoDetection())
}

func TestMockJob_AutoDetection(t *testing.T) {
	m := &MockJob{}
	expected := true

	assert.True(t, m.AutoDetection())
	m.AutoDetectionFunc = func() bool { return expected }
	assert.True(t, m.AutoDetection())
}

func TestMockJob_Tick(t *testing.T) {
	m := &MockJob{}

	assert.NotPanics(t, func() { m.Tick(1) })
}

func TestMockJob_Start(t *testing.T) {
	m := &MockJob{}

	assert.NotPanics(t, func() { m.Start() })
}

func TestMockJob_Stop(t *testing.T) {
	m := &MockJob{}

	assert.NotPanics(t, func() { m.Stop() })
}
