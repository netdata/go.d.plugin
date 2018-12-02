package modules

import (
	"github.com/netdata/go.d.plugin/logger"
)

// Module is an interface that represents a module.
type Module interface {
	// Init does initialization.
	// If it return false, the job will be disabled.
	Init() bool

	// Check is called after Init.
	// If it return false, the job will be disabled.
	Check() bool

	// Charts returns the chart definition.
	// Make sure not to share returned instance.
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
