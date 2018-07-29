package modules

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

type Module interface {
	Check() bool
	GetData() map[string]int64
	Charts
}

// Mandatory
type Charts interface {
	AddChart(...*charts.Chart)
	GetChart(string) *charts.Chart
	LookupChart(string) (*charts.Chart, bool)
	DeleteChart(string) bool
}

// Optional
type Logger interface {
	Error(...interface{})
	Warning(...interface{})
	Info(...interface{})
	Debug(...interface{})

	Errorf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

// BaseConfHook should be added by modules that need access to the base conf
// optional
type BaseConfHook interface {
	UpdateEvery() int
	// more methods can be added if needed
}

// NoConfiger should be added/implemented by modules which don't need configuration file
// optional
type NoConfiger interface {
	NoConfig()
}

// Unsafer should be added/implemented if module getData has a chance to panic
// optional
type Unsafer interface {
	Unsafe()
}
