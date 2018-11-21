package modules

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/l2isbad/go.d.plugin/logger"
)

const (
	penaltyStep = 5
	maxPenalty  = 600
)

func NewJob(modName string, module Module, out io.Writer) *job {
	buf := &bytes.Buffer{}
	return &job{
		moduleName: modName,
		module:     module,
		out:        out,

		runtimeChart: &Chart{
			typeName: "netdata",
			Title:    "Execution Time for",
			Units:    "ms",
			Fam:      "go.d",
			Ctx:      "netdata.go_d_execution_time", Priority: 145000,
			Dims: Dims{
				{ID: "time"},
			},
		},
		tick:         make(chan int),
		shutdownHook: make(chan struct{}),
		buf:          buf,
		apiWriter:    apiWriter{Writer: buf},
	}
}

type job struct {
	*logger.Logger
	module     Module
	moduleName string

	Nam             string `yaml:"name"`
	UpdateEvery     int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectRetry int    `yaml:"autodetection_retry" validate:"gte=0"`

	initialized bool
	panicked    bool

	runtimeChart *Chart
	charts       *Charts
	tick         chan int
	shutdownHook chan struct{}
	out          io.Writer
	buf          *bytes.Buffer
	apiWriter    apiWriter

	retries int
	prevRun time.Time
}

func (j job) FullName() string {
	if j.Name() == "" {
		return j.ModuleName()
	}
	return fmt.Sprintf("%s_%s", j.ModuleName(), j.Name())
}

func (j job) ModuleName() string {
	return j.moduleName
}

func (j job) Name() string {
	return j.Nam
}

func (j job) Initialized() bool {
	return j.initialized
}

func (j job) Panicked() bool {
	return j.panicked
}

func (j *job) Init() bool {
	defer func() {
		if r := recover(); r != nil {
			j.panicked = true
			j.module.Cleanup()
			j.Errorf("PANIC %v", r)
		}

	}()

	j.Logger = logger.New(j.ModuleName(), j.Name())
	j.module.SetLogger(j.Logger)

	return j.module.Init()
}

func (j *job) Check() bool {
	defer func() {
		if r := recover(); r != nil {
			j.panicked = true
			j.module.Cleanup()
			j.Errorf("PANIC %v", r)
		}
	}()

	return j.module.Check()
}

func (j *job) PostCheck() bool {
	j.charts = j.module.GetCharts()
	return j.charts != nil
}

func (j *job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Errorf("Skip the tick due to previous run hasn't been finished.")
	}
}

func (j *job) MainLoop() {
LOOP:
	for {
		select {
		case <-j.shutdownHook:
			break LOOP
		case t := <-j.tick:
			if t%(j.UpdateEvery+j.penalty()) != 0 {
				continue LOOP
			}
			j.runOnce()
		}
	}
}

func (j *job) runOnce() {
	curTime := time.Now()
	sinceLastRun := calcSinceLastRun(curTime, j.prevRun)

	data := j.getData()

	//if j.Panicked {
	//
	//}

	if j.populateMetrics(data, curTime, sinceLastRun) {
		j.prevRun = curTime
	} else {
		j.retries++
	}

	io.Copy(j.out, j.buf)
	j.buf.Reset()
}

func (j *job) Shutdown() {
	select {
	case j.shutdownHook <- struct{}{}:
		j.module.Cleanup()
	default:
	}
}

func (j job) AutoDetectionRetry() int {
	return j.AutoDetectRetry
}

func (j *job) getData() (result map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			j.panicked = true
		}
	}()
	return j.module.GetData()
}

func (j *job) populateMetrics(data map[string]int64, startTime time.Time, sinceLastRun int) bool {
	if !j.runtimeChart.pushed {
		j.runtimeChart.ID = fmt.Sprintf("execution_time_of_%s", j.FullName())
		j.createChart(j.runtimeChart)
	}

	var totalUpdated int
	elapsed := int64(durationTo(time.Now().Sub(startTime), time.Microsecond))

	for _, chart := range *j.charts {

		if !chart.pushed {
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

func (j *job) createChart(chart *Chart) {
	j.apiWriter.chart(
		firstNotEmpty(chart.typeName, j.FullName()),
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

	chart.pushed = true
}

func (j *job) updateChart(chart *Chart, data map[string]int64, sinceLastRun int) bool {
	if !chart.updated {
		sinceLastRun = 0
	}

	j.apiWriter.begin(
		firstNotEmpty(chart.typeName, j.FullName()),
		chart.ID,
		sinceLastRun,
	)

	var updated int

	for _, dim := range chart.Dims {
		if v, ok := data[dim.ID]; ok {
			j.apiWriter.set(dim.ID, v)
			updated++
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

func (j job) penalty() int {
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
