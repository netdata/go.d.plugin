package modules

import (
	"github.com/l2isbad/go.d.plugin/charts/cooked"
	"github.com/l2isbad/go.d.plugin/charts/raw"
)

type Module interface {
	CheckDataGetter // has to be implemented
	Charts          // has to be added
}

// CheckDataGetter must be implemented by any module
// mandatory
type CheckDataGetter interface {
	Check() bool
	GetData() map[string]int64
}

// Charts must be added by any module
// mandatory
type Charts interface {
	AddOne(*raw.Chart) error
	AddMany(*raw.Charts) int
	GetCharts() cooked.ChartsMap
	GetChartByID(string) *cooked.Chart
}

// BaseConfHook should be added by modules that need to get/set values from base conf
// optional
type BaseConfHook interface {
	UpdateEvery() int
	// more methods can be added if needed
}

// Logger should be added by modules that need to log messages
// optional
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

// NoConfiger should be added/implemented by modules that can work without configuration file
// optional
type NoConfiger interface {
	NoConfig()
}

// Unsafer should be added/implemented if module GetData has a chance to panic
// optional
type Unsafer interface {
	Unsafe()
}
