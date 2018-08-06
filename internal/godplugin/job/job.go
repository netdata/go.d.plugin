package job

import (
	"time"

	"bytes"

	"io"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/apiwriter"
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

// go:generate mockgen -source=job.go
type (
	Factory func(module modules.Module, config *Config, out io.Writer) Job

	Job interface {
		Init() error
		Check() bool
		PostCheck() bool
		MainLoop()
		Shutdown()
		Tick(clock int)

		ModuleName() string
		FullName() string
		JobName() string
		AutoDetectionRetry() int
	}
	job struct {
		*Config
		*logger.Logger
		Module       modules.Module
		tick         chan int
		shutdownHook chan int
		observer     *observer
		out          io.Writer
		buf          *bytes.Buffer
		apiWriter    apiwriter.APIWriter
		retries      int
		sinceLast    time.Time
	}
)

func (j *job) Tick(clock int) {
	select {
	case j.tick <- clock:
	default:
		j.Errorf("Skip the tick due to previous run has not been finished.")
	}
}

func New(module modules.Module, config *Config, out io.Writer) Job {
	buf := &bytes.Buffer{}
	return &job{
		Module:       module,
		Config:       config,
		tick:         make(chan int),
		shutdownHook: make(chan int),
		observer:     newObserver(config),
		out:          out,
		buf:          buf,
		apiWriter:    apiwriter.APIWriter{Writer: buf},
	}
}

func (j *job) Init() error {
	l := logger.New(j.RealModuleName, j.JobName())
	j.Logger = l

	j.Module.SetUpdateEvery(j.UpdateEvery)
	j.Module.SetModuleName(j.RealModuleName)
	j.Module.SetLogger(l)

	return j.Module.Init()
}

func (j *job) Check() bool {
	okCh := make(chan bool)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				j.Errorf("PANIC: %v", r)
				okCh <- false
			}
		}()
		okCh <- j.Module.Check()
	}()

	var ok bool
	select {
	case ok = <-okCh:
	case <-time.After(5 * time.Second):
		j.Error("check timeout")
	}
	return ok
}

func (j *job) PostCheck() bool {
	j.UpdateEvery = j.Module.UpdateEvery()

	modName := j.Module.ModuleName()
	logger.SetModName(j.Logger, modName)
	j.RealModuleName = modName

	c := j.Module.GetCharts()
	if c == nil {
		j.Error("GetCharts() [FAILED]")
		return false
	}
	j.observer.Set(c)
	return true
}

func (j *job) MainLoop() {
LOOP:
	for {
		select {
		case <-j.shutdownHook:
			break LOOP
		case t := <-j.tick:
			if t%j.UpdateEvery != 0 {
				continue LOOP
			}
		}
		data := j.getData()
		if data == nil {
			j.retries++
			continue
		}
		j.buf.Reset()
		// TODO write data
		io.Copy(j.out, j.buf)
	}
}

func (j *job) Shutdown() {
	select {
	case j.shutdownHook <- 1:
	default:
	}
}

func (j *job) getData() (result map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC: %v", r)
			result = nil
		}
	}()
	result = j.Module.GetData()
	return
}

func (j *job) AutoDetectionRetry() int {
	return j.Config.AutoDetectionRetry
}

//func (j *job) Start(wg *sync.WaitGroup) {
//Done:
//	for {
//
//		sleep := j.nextIn()
//		j.Debugf("sleeping for %s to reach frequency of %d sec", sleep, j.UpdateEvery)
//		time.Sleep(sleep)
//
//		j.curRun = time.Now()
//		if !j.lastRun.IsZero() {
//			j.sinceLast.Duration = j.curRun.Sub(j.lastRun)
//		}
//
//		if ok := j.update(); ok {
//			j.retries, j.penalty, j.lastRun = 0, 0, j.curRun
//			j.spentOnRun.Duration = time.Since(j.lastRun)
//
//		} else if !ok && !j.handleRetries() {
//			j.Errorf("stopped after %d collection failures in a row", j.MaxRetries)
//			break Done
//		}
//
//	}
//	wg.Done()
//}

//func (j *job) update() bool {
//
//	data := j.getData()
//
//	if data == nil {
//		j.Debug("getData() failed")
//		return false
//	}
//
//	var (
//		updated    int
//		active     int
//		suppressed int
//	)
//
//	for _, v := range *j.observer.charts {
//		if _, ok := j.observer.items[v.ID]; !ok {
//			j.observer.add(v)
//		}
//		chart := j.observer.items[v.ID]
//
//		if chart.obsoleted {
//			if canBeUpdated(*chart, data) {
//				chart.refresh()
//			} else {
//				suppressed++
//				continue
//			}
//		}
//
//		if j.ChartCleanup > 0 && chart.retries >= j.ChartCleanup {
//			j.Errorf("item '%s' was suppressed due to non updating", chart.item.ID)
//			chart.obsolete()
//			suppressed++
//			continue
//		}
//
//		active++
//		if chart.update(data, j.sinceLast.Round(time.Microsecond)) {
//			updated++
//		}
//	}
//
//	j.Debugf("update items: updated:%d, active:%d, suppressed:%d", updated, active, suppressed)
//	return updated > 0
//}
//

//
//func (j *job) nextIn() time.Duration {
//	start := time.Now()
//	next := start.Add(time.Duration(j.UpdateEvery) * time.Second).Add(j.penalty).Truncate(time.Second)
//	return time.Duration(next.UnixNano() - start.UnixNano())
//}
//
//func (j *job) handleRetries() bool {
//	j.retries++
//
//	if j.retries%5 != 0 {
//		return true
//	}
//
//	j.penalty = time.Duration(j.retries*j.UpdateEvery/2) * time.Second
//	j.Warningf(
//		"added %.0f seconds penalty after %d failed updates in a row",
//		j.penalty.Seconds(),
//		j.retries,
//	)
//	return j.retries < j.MaxRetries
//}
