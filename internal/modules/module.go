package modules

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

// TODO:
type Module interface {
	Check() bool
	GetData() map[string]int64

	Charts
	Logger
}

type Charts interface {
	AddChart(...*charts.Chart)
	GetChart(string) *charts.Chart
	LookupChart(string) (*charts.Chart, bool)
	DeleteChart(string) bool
	CopyCharts() charts.Charts
}

// Logger should be added by modules that need to log messages
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

// BaseConfHook should be added by modules that need to get/set values from base conf
// optional
type BaseConfHook interface {
	UpdateEvery() int
	// more methods can be added if needed
}

// NoConfiger should be added/implemented by modules that can work without configuration file
// optional
type NoConfiger interface {
	NoConfig()
}

// Unsafer should be added/implemented if module getData has a chance to panic
// optional
type Unsafer interface {
	Unsafe()
}
