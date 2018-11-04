package modules

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/l2isbad/go.d.plugin/logger"
)

type JobConfig struct {
	JobName            string `yaml:"job_name"`
	OverrideName       string `yaml:"name"`
	UpdateEvery        int    `yaml:"update_every" validate:"gte=1"`
	AutoDetectionRetry int    `yaml:"autodetection_retry" validate:"gte=0"`
	ChartCleanup       int    `yaml:"chart_cleanup" validate:"gte=0"`
	MaxRetries         int    `yaml:"retries" validate:"gte=0"`
}

// JobNewConfig returns JobConfig with default values
func JobNewConfig() *JobConfig {
	return &JobConfig{
		UpdateEvery:        1,
		AutoDetectionRetry: 0,
		ChartCleanup:       10,
		MaxRetries:         60,
	}
}

type Job struct {
	*JobConfig
	*logger.Logger
	module Module

	Inited     bool
	Panicked   bool
	ModuleName string

	charts       *Charts
	tick         chan int
	shutdownHook chan struct{}
	out          io.Writer
	buf          *bytes.Buffer
	apiWriter    apiWriter

	priority int
	retries  int
	prevRun  time.Time
}

func (j Job) Name() string {
	if j.ModuleName == j.JobName {
		return j.ModuleName
	}
	return fmt.Sprintf("%s_%s", j.ModuleName, j.JobName)
}

func NewJob(modName string, module Module, config *JobConfig, out io.Writer) *Job {
	buf := &bytes.Buffer{}
	return &Job{
		ModuleName:   modName,
		module:       module,
		JobConfig:    config,
		out:          out,
		tick:         make(chan int),
		shutdownHook: make(chan struct{}),
		buf:          buf,
		priority:     70000,
		apiWriter:    apiWriter{Writer: buf},
	}
}

func (j *Job) Init() error {
	j.Logger = logger.New(j.ModuleName, j.JobName)
	j.module.SetUpdateEvery(j.UpdateEvery)
	j.module.SetModuleName(j.ModuleName)
	j.module.SetLogger(j.Logger)

	return j.module.Init()
}

func (j *Job) Check() bool {
	defer func() {
		if r := recover(); r != nil {
			j.Panicked = true
			j.Errorf("PANIC %v", r)
		}

	}()
	return j.module.Check()
}

func (j *Job) PostCheck() bool {
	j.UpdateEvery = j.module.UpdateEvery()
	j.ModuleName = j.module.ModuleName()
	logger.SetModName(j.Logger, j.ModuleName)

	charts := j.module.GetCharts()
	if charts == nil {
		j.Error("GetCharts() [FAILED]")
		return false
	}

	j.charts = charts
	return true
}

func (j *Job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Errorf("Skip the tick due to previous run hasn't been finished.")
	}
}

func (j *Job) MainLoop() {
LOOP:
	for {
		select {
		case <-j.shutdownHook:
			break LOOP
		case t := <-j.tick:
			if t%j.UpdateEvery != 0 {
				continue LOOP
			}
			j.Info(11111111111)
		}

		//curTime := time.Now()
		//if j.prevRun.IsZero() {
		//	sinceLast := 0
		//} else {
		//	sinceLast := convertTo(curTime.Sub(j.prevRun), time.Microsecond)
		//}
		//
		//data := j.getData()
		//
		//if data == nil {
		//	j.retries++
		//	continue
		//}
		//j.buf.Reset()
		//// TODO write data
		//io.Copy(j.out, j.buf)
	}
}

func (j *Job) Shutdown() {
	select {
	case j.shutdownHook <- struct{}{}:
		j.module.Cleanup()
	default:
	}
}

func (j *Job) getData() (result map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			j.Panicked = true
		}
	}()
	return j.module.GetData()
}

func (j *Job) AutoDetectionRetry() int {
	return j.JobConfig.AutoDetectionRetry
}

func (j *Job) PopulateMetrics(data map[string]int64, sinceLast int) bool {
	var totalUpdated int

	for _, chart := range *j.charts {
		if !chart.pushed {
			j.apiWriter.chart(
				j.Name(),
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
				j.ModuleName,
			)
			chart.pushed = true
		}

		if data == nil {
			continue
		}

		if chart.Obsolete {
			if !canChartBeUpdated(chart, data) {
				continue
			}
			chart.Obsolete = false
			chart.pushed = false
		}

		if !chart.updated {
			sinceLast = 0
		}

		j.apiWriter.begin("typeName", chart.ID, sinceLast)

		chart.updated = false
		for _, dim := range chart.Dims {
			if v, ok := data[dim.ID]; ok {
				j.apiWriter.set(dim.ID, v)
				chart.updated = true
			}
		}

		for _, variable := range chart.Vars {
			if v, ok := data[variable.ID]; ok {
				j.apiWriter.set(variable.ID, v)
			}
		}

		j.apiWriter.end()

		if !chart.updated {
			chart.retries++
		}

		if j.ChartCleanup > 0 && chart.retries >= j.ChartCleanup {
			chart.Obsolete = true
			chart.pushed = false
		}
	}

	return totalUpdated > 0
}

func convertTo(from time.Duration, to time.Duration) int {
	return int(int64(from) / (int64(to) / int64(time.Nanosecond)))
}

func canChartBeUpdated(chart *Chart, data map[string]int64) bool {
	for _, dim := range chart.Dims {
		if _, ok := data[dim.ID]; ok {
			return true
		}
	}
	return false
}
