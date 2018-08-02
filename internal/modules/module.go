package modules

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

// Module Module
type Module interface {
	// Init is called after UpdateEvery, ModuleName are set.
	Init()

	// Check is called after Init or AutoDetectionRetry.
	// If it return false, this job will be disabled.
	Check() bool
	// GetCharts returns the chart definition.
	// Make sure not to share the return instance.
	GetCharts() *charts.Charts

	// GetData returns the collected metrics.
	GetData() map[string]int64

	// SetLogger SetLogger
	SetLogger(l *logger.Logger)

	// UpdateEvery UpdateEvery
	UpdateEvery() int
	// SetUpdateEvery SetUpdateEvery
	SetUpdateEvery(v int)

	// ModuleName ModuleName
	ModuleName() string
	// SetModuleName SetModuleName
	SetModuleName(v string)
}

// ModuleBase is a helper struct. All modules should embed this struct.
type ModuleBase struct {
	*logger.Logger
	updateEvery int
	moduleName  string
}

// SetLogger SetLogger
func (m *ModuleBase) SetLogger(l *logger.Logger) {
	m.Logger = l
}

// SetUpdateEvery SetUpdateEvery
func (m *ModuleBase) SetUpdateEvery(v int) {
	m.updateEvery = v
}

// SetModuleName SetModuleName
func (m *ModuleBase) SetModuleName(v string) {
	m.moduleName = v
}

// UpdateEvery UpdateEvery
func (m ModuleBase) UpdateEvery() int {
	return m.updateEvery
}

// ModuleName ModuleName
func (m ModuleBase) ModuleName() string {
	return m.moduleName
}
