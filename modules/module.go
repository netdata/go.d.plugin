package modules

import (
	"github.com/l2isbad/go.d.plugin/logger"
)

// Module Module
type Module interface {
	// Init is called after UpdateEvery, ModuleName are set.
	Init() error

	// Check is called after Init or AutoDetectionRetry.
	// If it return false, this job will be disabled.
	Check() bool

	// GetCharts returns the chart definition.
	// Make sure not to share the return instance.
	GetCharts() *Charts

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

	// Cleanup Cleanup
	Cleanup()
}

// Base is a helper struct. All modules should embed this struct.
type Base struct {
	*logger.Logger
	updateEvery int
	moduleName  string
}

// Init Init
func (b *Base) Init() error { return nil }

// SetLogger SetLogger
func (b *Base) SetLogger(l *logger.Logger) { b.Logger = l }

// SetUpdateEvery SetUpdateEvery
func (b *Base) SetUpdateEvery(v int) { b.updateEvery = v }

// SetModuleName SetModuleName
func (b *Base) SetModuleName(v string) { b.moduleName = v }

// UpdateEvery UpdateEvery
func (b Base) UpdateEvery() int { return b.updateEvery }

// moduleName ModuleName
func (b Base) ModuleName() string { return b.moduleName }

// Cleanup Cleanup
func (b Base) Cleanup() {}
