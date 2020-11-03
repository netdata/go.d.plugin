package module

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/agent/netdataapi"
	"github.com/netdata/go.d.plugin/logger"
)

var writeLock = &sync.Mutex{}

func newRuntimeChart(pluginName string) *Chart {
	return &Chart{
		typeID: "netdata",
		Units:  "ms",
		Fam:    pluginName,
		Ctx:    "netdata.go_plugin_execution_time", Priority: 145000,
		Dims: Dims{
			{ID: "time"},
		},
	}
}

type JobConfig struct {
	PluginName      string
	Name            string
	ModuleName      string
	FullName        string
	Module          Module
	Out             io.Writer
	UpdateEvery     int
	AutoDetectEvery int
	Priority        int
}

const (
	penaltyStep = 5
	maxPenalty  = 600
	infTries    = -1
)

func NewJob(cfg JobConfig) *Job {
	var buf bytes.Buffer
	return &Job{
		pluginName:      cfg.PluginName,
		name:            cfg.Name,
		moduleName:      cfg.ModuleName,
		fullName:        cfg.FullName,
		updateEvery:     cfg.UpdateEvery,
		AutoDetectEvery: cfg.AutoDetectEvery,
		priority:        cfg.Priority,
		module:          cfg.Module,
		out:             cfg.Out,
		AutoDetectTries: infTries,
		runChart:        newRuntimeChart(cfg.PluginName),
		stop:            make(chan struct{}),
		tick:            make(chan int),
		buf:             &buf,
		api:             netdataapi.New(&buf),
	}
}

// Job represents a job. It's a module wrapper.
type Job struct {
	pluginName string
	name       string
	moduleName string
	fullName   string

	updateEvery     int
	AutoDetectEvery int
	AutoDetectTries int
	priority        int

	*logger.Logger

	module Module

	initialized bool
	panicked    bool

	runChart *Chart
	charts   *Charts
	tick     chan int
	out      io.Writer
	buf      *bytes.Buffer
	api      *netdataapi.API

	retries int
	prevRun time.Time

	stop chan struct{}
}

// https://github.com/netdata/netdata/blob/ab0ffcebf802803d1e88f6a5e47a314c292b45e3/database/rrd.h#L59
// Chart type.id (job.FullName() + '.' + chart.ID)
const RRD_ID_LENGTH_MAX = 200

// FullName returns job full name.
func (j Job) FullName() string {
	return j.fullName
}

// ModuleName returns job module name.
func (j Job) ModuleName() string {
	return j.moduleName
}

// Name returns job name.
func (j Job) Name() string {
	return j.name
}

// Panicked returns 'panicked' flag value.
func (j Job) Panicked() bool {
	return j.panicked
}

// AutoDetectionEvery returns value of AutoDetectEvery.
func (j Job) AutoDetectionEvery() int {
	return j.AutoDetectEvery
}

// RetryAutoDetection returns whether it is needed to retry autodetection.
func (j Job) RetryAutoDetection() bool {
	return j.AutoDetectEvery > 0 && (j.AutoDetectTries == infTries || j.AutoDetectTries > 0)
}

// AutoDetection invokes init, check and postCheck. It handles panic.
func (j *Job) AutoDetection() (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			j.Errorf("PANIC %v", r)
			j.panicked = true
			j.disableAutoDetection()
		}
		if !ok {
			j.module.Cleanup()
		}
	}()

	if ok = j.init(); !ok {
		j.Error("init failed")
		j.disableAutoDetection()
		return
	}
	if ok = j.check(); !ok {
		j.Error("check failed")
		return
	}
	j.Info("check success")
	if ok = j.postCheck(); !ok {
		j.Error("postCheck failed")
		j.disableAutoDetection()
		return
	}
	return true
}

// Tick Tick.
func (j *Job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Debug("skip the tick due to previous run hasn't been finished")
	}
}

// Start starts job main loop.
func (j *Job) Start() {
	j.Infof("started, data collection interval %ds", j.updateEvery)
	defer func() { j.Info("stopped") }()

LOOP:
	for {
		select {
		case <-j.stop:
			break LOOP
		case t := <-j.tick:
			if t%(j.updateEvery+j.penalty()) == 0 {
				j.runOnce()
			}
		}
	}
	j.module.Cleanup()
	j.cleanup()
	j.stop <- struct{}{}
}

// Stop stops job main loop. It blocks until the job is stopped.
func (j *Job) Stop() {
	// TODO: should have blocking and non blocking stop
	j.stop <- struct{}{}
	<-j.stop
}

func (j *Job) disableAutoDetection() {
	j.AutoDetectEvery = 0
}

func (j *Job) cleanup() {
	if j.Logger != nil {
		logger.GlobalMsgCountWatcher.Unregister(j.Logger)
	}
	j.buf.Reset()

	if j.runChart.created {
		j.runChart.MarkRemove()
		j.createChart(j.runChart)
	}
	if j.charts != nil {
		for _, chart := range *j.charts {
			if chart.created {
				chart.MarkRemove()
				j.createChart(chart)
			}
		}
	}
	writeLock.Lock()
	_, _ = io.Copy(j.out, j.buf)
	writeLock.Unlock()
}

func (j *Job) init() bool {
	if j.initialized {
		return true
	}

	log := logger.NewLimited(j.ModuleName(), j.Name())
	j.Logger = log
	j.module.GetBase().Logger = log

	j.initialized = j.module.Init()
	return j.initialized
}

func (j *Job) check() bool {
	ok := j.module.Check()
	if !ok && j.AutoDetectTries != infTries {
		j.AutoDetectTries--
	}
	return ok
}

func (j *Job) postCheck() bool {
	if j.charts = j.module.Charts(); j.charts == nil {
		j.Error("nil charts")
		return false
	}
	if err := checkCharts(*j.charts...); err != nil {
		j.Errorf("charts check: %v", err)
		return false
	}
	return true
}

func (j *Job) runOnce() {
	curTime := time.Now()
	sinceLastRun := calcSinceLastRun(curTime, j.prevRun)
	j.prevRun = curTime

	metrics := j.collect()

	if j.panicked {
		return
	}

	if j.processMetrics(metrics, curTime, sinceLastRun) {
		j.retries = 0
	} else {
		j.retries++
	}

	writeLock.Lock()
	_, _ = io.Copy(j.out, j.buf)
	writeLock.Unlock()
	j.buf.Reset()
}

func (j *Job) collect() (result map[string]int64) {
	j.panicked = false
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			j.panicked = true
		}
	}()
	return j.module.Collect()
}

func (j *Job) processMetrics(metrics map[string]int64, startTime time.Time, sinceLastRun int) bool {
	if !j.runChart.created {
		j.runChart.ID = fmt.Sprintf("execution_time_of_%s", j.FullName())
		j.runChart.Title = fmt.Sprintf("Execution Time for %s", j.FullName())
		j.createChart(j.runChart)
	}

	elapsed := int64(durationTo(time.Since(startTime), time.Millisecond))

	var i, updated int
	for _, chart := range *j.charts {
		if !chart.created {
			typeID := fmt.Sprintf("%s.%s", j.FullName(), chart.ID)
			if len(typeID) >= RRD_ID_LENGTH_MAX {
				j.Warningf("chart 'type.id' length (%d) >= max allowed (%d), the chart is ignored (%s)",
					len(typeID), RRD_ID_LENGTH_MAX, typeID)
				chart.ignore = true
			}
			j.createChart(chart)
		}
		if chart.remove {
			continue
		}
		(*j.charts)[i] = chart
		i++
		if len(metrics) == 0 || chart.Obsolete {
			continue
		}
		if j.updateChart(chart, metrics, sinceLastRun) {
			updated++
		}
	}
	*j.charts = (*j.charts)[:i]

	if updated == 0 {
		return false
	}
	j.updateChart(j.runChart, map[string]int64{"time": elapsed}, sinceLastRun)
	return true
}

func (j *Job) createChart(chart *Chart) {
	defer func() { chart.created = true }()
	if chart.ignore {
		return
	}

	if chart.Priority == 0 {
		chart.Priority = j.priority
		j.priority++
	}
	_ = j.api.CHART(
		firstNotEmpty(chart.typeID, j.FullName()),
		chart.ID,
		chart.OverID,
		chart.Title,
		chart.Units,
		chart.Fam,
		chart.Ctx,
		chart.Type.String(),
		chart.Priority,
		j.updateEvery,
		chart.Opts.String(),
		j.pluginName,
		j.moduleName,
	)
	for _, dim := range chart.Dims {
		_ = j.api.DIMENSION(
			dim.ID,
			dim.Name,
			dim.Algo.String(),
			handleZero(dim.Mul),
			handleZero(dim.Div),
			dim.DimOpts.String(),
		)
	}
	for _, v := range chart.Vars {
		_ = j.api.VARIABLE(
			v.ID,
			v.Value,
		)
	}
	_ = j.api.EMPTYLINE()
}

func (j *Job) updateChart(chart *Chart, collected map[string]int64, sinceLastRun int) bool {
	if chart.ignore {
		dims := chart.Dims[:0]
		for _, dim := range chart.Dims {
			if !dim.remove {
				dims = append(dims, dim)
			}
		}
		chart.Dims = dims
		return false
	}

	if !chart.updated {
		sinceLastRun = 0
	}

	_ = j.api.BEGIN(
		firstNotEmpty(chart.typeID, j.FullName()),
		chart.ID,
		sinceLastRun,
	)
	var i, updated int
	for _, dim := range chart.Dims {
		if dim.remove {
			continue
		}
		chart.Dims[i] = dim
		i++
		if v, ok := collected[dim.ID]; !ok {
			_ = j.api.SETEMPTY(dim.ID)
		} else {
			_ = j.api.SET(dim.ID, v)
			updated++
		}
	}
	chart.Dims = chart.Dims[:i]

	for _, vr := range chart.Vars {
		if v, ok := collected[vr.ID]; ok {
			_ = j.api.VARIABLE(vr.ID, v)
		}

	}
	_ = j.api.END()

	if chart.updated = updated > 0; chart.updated {
		chart.Retries = 0
	} else {
		chart.Retries++
	}
	return chart.updated
}

func (j Job) penalty() int {
	v := j.retries / penaltyStep * penaltyStep * j.updateEvery / 2
	if v > maxPenalty {
		return maxPenalty
	}
	return v
}

func calcSinceLastRun(curTime, prevRun time.Time) int {
	if prevRun.IsZero() {
		return 0
	}
	return int((curTime.UnixNano() - prevRun.UnixNano()) / 1000)
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

func handleZero(v int) int {
	if v == 0 {
		return 1
	}
	return v
}
