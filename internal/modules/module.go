package modules

import "github.com/l2isbad/go.d.plugin/internal/pkg/charts"

type Module interface {
	CheckDataGetter // has to be implemented
	Charts          // has to be added
	Logger          // has to be added
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
	AddChart(...charts.Chart)
	GetChart(string) Chart
	LookupChart(string) (Chart, bool)
}

// BaseConfHook should be added by modules that need to get/set values from base conf
// optional
type BaseConfHook interface {
	UpdateEvery() int
	// more methods can be added if needed
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

// User Chart
type Chart interface {
	AddDim(charts.Dimension)
	AddVar(charts.Variable)
	//GetDimByID(string) Dimension
	//GetVarByID(string) Variable
}

//type Dimension interface {
//}
//
//type Variable interface {
//}
