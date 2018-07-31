package job

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func New(m modules.Module, c *Config) *Job {
	_, u := m.(modules.Unsafer)

	return &Job{
		Mod:    m,
		Config: c,
		unsafe: u,
		Obs:    newObserver(c),
	}
}

type Job struct {
	Mod modules.Module

	timers
	*Config
	*logger.Logger
	Obs *observer

	retries int
	unsafe  bool
}

func (j *Job) Start(wg *sync.WaitGroup) {
Done:
	for {

		sleep := j.nextIn()
		j.Debugf("sleeping for %s to reach frequency of %d sec", sleep, j.UpdateEvery)
		time.Sleep(sleep)

		j.curRun = time.Now()
		if !j.lastRun.IsZero() {
			j.sinceLast.Duration = j.curRun.Sub(j.lastRun)
		}

		if ok := j.update(); ok {
			j.retries, j.penalty, j.lastRun = 0, 0, j.curRun
			j.spentOnRun.Duration = time.Since(j.lastRun)

		} else if !ok && !j.handleRetries() {
			j.Errorf("stopped after %d collection failures in a row", j.MaxRetries)
			break Done
		}

	}
	wg.Done()
}

func (j *Job) update() bool {

	data := j.getData()

	if data == nil {
		j.Debug("GetData() failed")
		return false
	}

	var (
		updated    int
		active     int
		suppressed int
	)

	for _, v := range *j.Obs.charts {
		if _, ok := j.Obs.items[v.ID]; !ok {
			j.Obs.add(v)
		}
		chart := j.Obs.items[v.ID]

		if chart.obsoleted {
			if canBeUpdated(*chart, data) {
				chart.refresh()
			} else {
				suppressed++
				continue
			}
		}

		if j.ChartCleanup > 0 && chart.retries >= j.ChartCleanup {
			j.Errorf("item '%s' was suppressed due to non updating", chart.item.ID)
			chart.obsolete()
			suppressed++
			continue
		}

		active++
		if chart.update(data, j.sinceLast.ConvertTo(time.Microsecond)) {
			updated++
		}
	}

	j.Debugf("update items: updated:%d, active:%d, suppressed:%d", updated, active, suppressed)
	return updated > 0
}

func (j *Job) safeGetData() (m map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC(%s)", r)
		}
	}()

	m = j.Mod.GetData()
	return
}

func (j *Job) getData() map[string]int64 {
	if j.unsafe {
		return j.safeGetData()
	}
	return j.Mod.GetData()
}

func (j *Job) nextIn() time.Duration {
	start := time.Now()
	next := start.Add(time.Duration(j.UpdateEvery) * time.Second).Add(j.penalty).Truncate(time.Second)
	return time.Duration(next.UnixNano() - start.UnixNano())
}

func (j *Job) handleRetries() bool {
	j.retries++

	if j.retries%5 != 0 {
		return true
	}

	j.penalty = time.Duration(j.retries*j.UpdateEvery/2) * time.Second
	j.Warningf(
		"added %.0f seconds penalty after %d failed updates in a row",
		j.penalty.Seconds(),
		j.retries,
	)
	return j.retries < j.MaxRetries
}
