package modules

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/netdata/go.d.plugin/logger"
)

const (
	penaltyStep = 5
	maxPenalty  = 600
)

type Observer interface {
	RemoveFromQueue(fullName string)
}

func NewJob(modName string, module Module, out io.Writer, observer Observer) *Job {
	buf := &bytes.Buffer{}
	return &Job{
		moduleName: modName,
		module:     module,
		out:        out,
		observer:   observer,

		runtimeChart: &Chart{
			typeID: "netdata",
			Title:  "Execution Time for",
			Units:  "ms",
			Fam:    "go.d",
			Ctx:    "netdata.go_d_execution_time", Priority: 145000,
			Dims: Dims{
				{ID: "time"},
			},
		},
		stopHook:  make(chan struct{}, 1),
		tick:      make(chan int),
		buf:       buf,
		apiWriter: apiWriter{Writer: buf},
		priority:  70000,
	}
}

type Job struct {
	*logger.Logger
	module     Module
	moduleName string

	Nam             string `yaml:"name"`
	UpdateEvery     int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectRetry int    `yaml:"autodetection_retry" validate:"gte=0"`

	initialized bool
	panicked    bool

	stopHook     chan struct{}
	observer     Observer
	runtimeChart *Chart
	charts       *Charts
	tick         chan int
	out          io.Writer
	buf          *bytes.Buffer
	apiWriter    apiWriter

	priority int
	retries  int
	prevRun  time.Time
}

func (j Job) FullName() string {
	if j.Name() == "" {
		return j.ModuleName()
	}
	return fmt.Sprintf("%s_%s", j.ModuleName(), j.Name())
}

func (j Job) ModuleName() string {
	return j.moduleName
}

func (j Job) Name() string {
	if j.Nam == "" {
		return j.moduleName
	}
	return j.Nam
}

func (j Job) Initialized() bool {
	return j.initialized
}

func (j Job) Panicked() bool {
	return j.panicked
}

func (j *Job) Init() bool {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC %v", r)
			j.panicked = true
			j.module.Cleanup()
		}
	}()

	j.Logger = logger.New(j.ModuleName(), j.Name())
	j.module.SetLogger(j.Logger)

	return j.module.Init()
}

func (j *Job) Check() bool {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC %v", r)
			j.panicked = true
			j.module.Cleanup()
		}
	}()

	return j.module.Check()
}

func (j *Job) PostCheck() bool {
	j.charts = j.module.GetCharts()

	if j.charts == nil {
		j.Error("charts can't be nil")
		j.module.Cleanup()
		return false
	}

	return true
}

func (j *Job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Errorf("Skip the tick due to previous run hasn't been finished.")
	}
}

func (j *Job) Start() {
	j.MainLoop()
}

func (j *Job) Stop() {
	j.stopHook <- struct{}{}
}

func (j *Job) MainLoop() {
LOOP:
	for {
		select {
		case <-j.stopHook:
			j.module.Cleanup()
			break LOOP
		case t := <-j.tick:
			doRun := t%(j.UpdateEvery+j.penalty()) == 0
			if doRun {
				j.runOnce()
			}
		}
	}
}

func (j *Job) runOnce() {
	curTime := time.Now()
	sinceLastRun := calcSinceLastRun(curTime, j.prevRun)

	data := j.getData()

	if j.panicked {
		j.observer.RemoveFromQueue(j.FullName())
		go j.module.Cleanup()
		return
	}

	if j.populateMetrics(data, curTime, sinceLastRun) {
		j.prevRun = curTime
	} else {
		j.retries++
	}

	io.Copy(j.out, j.buf)
	j.buf.Reset()
}

func (j Job) AutoDetectionRetry() int {
	return j.AutoDetectRetry
}

func (j *Job) getData() (result map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			j.panicked = true
		}
	}()
	return j.module.GetData()
}

func (j *Job) populateMetrics(data map[string]int64, startTime time.Time, sinceLastRun int) bool {
	if !j.runtimeChart.created {
		j.runtimeChart.ID = fmt.Sprintf("execution_time_of_%s", j.FullName())
		j.createChart(j.runtimeChart)
	}

	var totalUpdated int
	elapsed := int64(durationTo(time.Now().Sub(startTime), time.Microsecond))

	for _, chart := range *j.charts {

		if !chart.created {
			j.createChart(chart)
		}

		if data == nil || chart.Obsolete {
			continue
		}

		if j.updateChart(chart, data, sinceLastRun) {
			totalUpdated++
		}
	}

	if totalUpdated == 0 {
		return false
	}

	j.updateChart(
		j.runtimeChart,
		map[string]int64{"time": elapsed},
		sinceLastRun,
	)

	return true
}

func (j *Job) createChart(chart *Chart) {
	if chart.Priority == 0 {
		chart.Priority = j.priority
		j.priority++
	}
	j.apiWriter.chart(
		firstNotEmpty(chart.typeID, j.FullName()),
		chart.ID,
		chart.OverID,
		chart.Title,
		chart.Units,
		chart.Fam,
		chart.Ctx,
		chart.Type,
		chart.Priority,
		j.UpdateEvery,
		chart.Opts,
		j.moduleName,
	)
	for _, dim := range chart.Dims {
		j.apiWriter.dimension(
			dim.ID,
			dim.Name,
			dim.Algo,
			dim.Mul,
			dim.Div,
			dim.Hidden,
		)
	}
	for _, v := range chart.Vars {
		if v.Value == 0 {
			continue
		}
		j.apiWriter.set(
			v.ID,
			v.Value,
		)
	}
	j.apiWriter.Write([]byte("\n"))

	chart.created = true
}

func (j *Job) updateChart(chart *Chart, data map[string]int64, sinceLastRun int) bool {
	if !chart.updated {
		sinceLastRun = 0
	}

	j.apiWriter.begin(
		firstNotEmpty(chart.typeID, j.FullName()),
		chart.ID,
		sinceLastRun,
	)

	var updated int

	for _, dim := range chart.Dims {
		if v, ok := data[dim.ID]; ok {
			j.apiWriter.set(dim.ID, v)
			updated++
		} else {
			j.apiWriter.setEmpty(dim.ID)
		}
	}
	for _, variable := range chart.Vars {
		if v, ok := data[variable.ID]; ok {
			j.apiWriter.set(variable.ID, v)
		}
	}

	j.apiWriter.end()

	chart.updated = updated > 0

	if chart.updated {
		chart.Retries = 0
	} else {
		chart.Retries++
	}

	return chart.updated
}

func (j Job) penalty() int {
	v := j.retries / penaltyStep * penaltyStep * j.UpdateEvery / 2
	if v > maxPenalty {
		return maxPenalty
	}
	return v
}

func calcSinceLastRun(curTime, prevRun time.Time) int {
	if prevRun.IsZero() {
		return 0
	}
	return durationTo(curTime.Sub(prevRun), time.Microsecond)
}

func durationTo(duration time.Duration, to time.Duration) int {
	return int(int64(duration) / (int64(to) / int64(time.Nanosecond)))
}

func firstNotEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
