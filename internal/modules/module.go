package modules

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

type Module interface {
	Init()
	Check() bool
	GetData() map[string]int64
	GetCharts() *charts.Charts

	SetLogger(l *logger.Logger)
	SetUpdateEvery(v int)
	SetModuleName(v string)
	UpdateEvery() int
	ModuleName() string
}

type ModuleBase struct {
	*logger.Logger
	updateEvery int
	moduleName  string
}

func (m *ModuleBase) SetLogger(l *logger.Logger) {
	m.Logger = l
}

func (m *ModuleBase) SetUpdateEvery(v int) {
	m.updateEvery = v
}

func (m *ModuleBase) SetModuleName(v string) {
	m.moduleName = v
}

func (m ModuleBase) UpdateEvery() int {
	return m.updateEvery
}

func (m ModuleBase) ModuleName() string {
	return m.moduleName
}

// NoConfiger should be added/implemented by modules which don't need configuration file
type NoConfiger interface {
	NoConfig()
}

// Unsafer should be added/implemented if module getData has a chance to panic
type Unsafer interface {
	Unsafe()
}
