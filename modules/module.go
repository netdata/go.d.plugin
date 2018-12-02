package modules

import (
	"github.com/netdata/go.d.plugin/logger"
)

// Module Module
type Module interface {
	// Init is called after UpdateEvery, ModuleName are set.
	Init() bool

	// Check is called after Init or AutoDetectionRetry.
	// If it return false, this job will be disabled.
	Check() bool

	// Charts returns the chart definition.
	// Make sure not to share the return instance.
	Charts() *Charts

	// GatherMetrics returns metrics.
	GatherMetrics() map[string]int64

	// SetLogger SetLogger
	SetLogger(l *logger.Logger)

	// Cleanup Cleanup
	Cleanup()
}

// Base is a helper struct. All modules should embed this struct.
type Base struct {
	*logger.Logger
}

// SetLogger SetLogger
func (b *Base) SetLogger(l *logger.Logger) { b.Logger = l }
