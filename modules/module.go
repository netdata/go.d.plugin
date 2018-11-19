package modules

import (
	"github.com/l2isbad/go.d.plugin/logger"
)

// Module Module
type Module interface {
	// Init is called after UpdateEvery, ModuleName are set.
	Init() bool

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

	// Cleanup Cleanup
	Cleanup()
}

// Base is a helper struct. All modules should embed this struct.
type Base struct {
	*logger.Logger
}

// Init Init
func (b *Base) Init() bool { return true }

// SetLogger SetLogger
func (b *Base) SetLogger(l *logger.Logger) { b.Logger = l }

// Cleanup Cleanup
func (b Base) Cleanup() {}
