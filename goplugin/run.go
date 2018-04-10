package goplugin

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/charts/cooked"
	"github.com/l2isbad/go.d.plugin/goplugin/job"
	"github.com/l2isbad/go.d.plugin/logger"
	"github.com/l2isbad/go.d.plugin/modules"
)

const (
	fieldCharts   = "Charts"
	fieldLogger   = "Logger"
	fieldBaseConf = "BaseConfHook"
)

func New(updEvery int, modRun modules.CreatorsMap, l Logger, modConf string) *goPlugin {
	return &goPlugin{
		overrideUpdEvery: updEvery,
		modRun:           modRun,
		modConf:          modConf,
		Log:              l,
	}
}

type goPlugin struct {
	overrideUpdEvery int
	modRun           modules.CreatorsMap
	modConf          string
	Log              Logger

	wg sync.WaitGroup
}

// Run loads all jobs, sets Charts, Logger and BaseConf fields and runs jobCheck for every job in sequential order.
// Then it waits until all jobs are completed.
func (p *goPlugin) Run() {
	if len(p.modRun) == 0 {
		return
	}
	loaded := p.loadJobs()
	started := make(map[string]bool)

	for !loaded.empty() {
		j := loaded.pop()
		if err := setJobFields(j.Module, j.Conf); err != nil {
			p.Log.Errorf("'%s %s' %s", j.ModuleName(), j.JobName(), err)
			continue
		}
		j.Logger = logger.CacheGet(j)

		// Module name can be changed in jobCheck(), so we need to copy it before.
		key := j.FullName()

		if started[key] {
			j.Error("[DROPPED] already served by another job")
			continue
		}
		if p.jobCheck(j) {
			started[key] = true
		}

	}
	p.wg.Wait()
	fmt.Println("DISABLE")
}

// jobCheck executes job.Check() in a separate goroutine and starts the 5 second timer.
// Timer is needed because plugin runs job checks sequentially.
// If the job.Check() does not return a value before the timer expires, it returns false.
// All jobs passed the Check() will be started.
// If job Check() fails but it's a AutoDetectionRetry job, jobRecheck(job) will be started.
func (p *goPlugin) jobCheck(j *job.Job) bool {
	p.jobUpdEveryOverride(j)
	check := make(chan bool)
	go func() {
		check <- j.Check()
	}()

	select {
	case ok := <-check:
		if !ok && j.AutoDetectionRetry == 0 {
			j.Error("Check() [FAILED]")
			return false
		}

		if !ok {
			j.Warningf("Check() [RECHECK EVERY %d]", j.AutoDetectionRetry)
			p.wg.Add(1)
			go p.jobRecheck(j) // GO
			return true
		}

		p.wg.Add(1)
		j.Info("Check() [OK]")
		go j.Start(&p.wg) // GO
		return true

	case <-time.After(5 * time.Second):
		j.Error("Check() [TIMEOUT]")
	}
	return false
}

// jobRecheck executes job Check() every job AutoDetectionRetry seconds in a for loop.
// If Check() returns true, jobRecheck exit for loop.
func (p *goPlugin) jobRecheck(j *job.Job) {
	var c int
	for {
		c++
		time.Sleep(time.Duration(j.AutoDetectionRetry) * time.Second)
		if j.Check() {
			j.Infof("Check() [OK] after %d rechecks", c)
			j.Start(&p.wg)
			return
		}
	}
}

// jobUpdEveryOverride overrides job UpdateEvery to UpdateEvery from cmd line if
// 1. plugin runs in DEBUG mode
// 2. job UpdateEvery lower then UpdateEvery from cmd
func (p *goPlugin) jobUpdEveryOverride(j *job.Job) {
	if (p.Log.Level() == logger.DEBUG && j.UpdEvery != p.overrideUpdEvery) || j.UpdEvery < p.overrideUpdEvery {
		j.Infof("update every is changed %d => %d", j.UpdEvery, p.overrideUpdEvery)
		j.SetUpdateEvery(p.overrideUpdEvery)
	}
}

// Set *cooked.charts as modules.Charts interface       (mandatory)
// Set *logger.logger as modules.Logger interface       (optional)
// Set *jobs.Conf     as modules.BaseConfHook interface (optional)
func setJobFields(mod interface{}, conf *job.Conf) error {
	v := reflect.ValueOf(mod)
	if v.Kind() != reflect.Ptr {
		return errors.New("module must be a pointer")
	}
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("module must be a struct")
	}
	// Mandatory field
	if !setField(&elem, conf, fieldCharts) {
		return errors.New("'Charts' field must be a 'modules.Charts' interface")
	}
	// Optional fields
	// logger.New adds new logger instance to the cache
	logger.New(conf)
	setField(&elem, conf, fieldLogger)
	setField(&elem, conf, fieldBaseConf)
	return nil
}

func setField(v *reflect.Value, conf *job.Conf, fieldName string) bool {
	f := v.FieldByName(fieldName)
	if f.IsValid() && f.Kind() == reflect.Interface && f.Type().Name() == fieldName && f.CanSet() {
		switch fieldName {
		case fieldCharts:
			f.Set(reflect.ValueOf(cooked.NewCharts(conf)))
		case fieldLogger:
			f.Set(reflect.ValueOf(logger.CacheGet(conf)))
		case fieldBaseConf:
			f.Set(reflect.ValueOf(conf))
		}
		return true
	}
	return false
}
