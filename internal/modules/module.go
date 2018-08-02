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

	RequireConfig() bool
	DefaultUpdateEvery(globalUpdateEvery int) int
	DefaultChartCleanup() int
	DisabledByDefault() bool

	SetLogger(l *logger.Logger)
	UpdateEvery() int
	SetUpdateEvery(v int)
	ModuleName() string
	SetModuleName(v string)
}

type ModuleBase struct {
	*logger.Logger
	updateEvery int
	moduleName  string
}

func (m *ModuleBase) RequireConfig() bool {
	return true
}

func (m *ModuleBase) DefaultUpdateEvery(globalUpdateEvery int) int {
	return globalUpdateEvery
}
func (m *ModuleBase) DefaultChartCleanup() int {
	return 0
}
func (m *ModuleBase) DisabledByDefault() bool {
	return false
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

// Unsafer should be added/implemented if module getData has a chance to panic
type Unsafer interface {
	Unsafe()
}
