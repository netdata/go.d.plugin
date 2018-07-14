package job

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
)

func New(m modules.Module, c *Config) *Job {
	_, u := m.(modules.Unsafer)

	return &Job{
		Module: m,
		Config: c,
		unsafe: u,
	}
}

type Job struct {
	modules.Module
	*Config
	timers
	retries int
	unsafe  bool
}

func (j *Job) Run(wg *sync.WaitGroup) {
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
			j.Errorf("stopped after %d collection failures in a row", j.RetriesMax)
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

	for _, chart := range j.GetCharts() {

		if chart.IsObsoleted() {
			if !chart.CanBeUpdated(data) {
				suppressed++
				continue
			}
			chart.Refresh()

		} else if j.ChartCleanup > 0 && chart.FailedUpdates >= j.ChartCleanup {
			j.Errorf("chart '%s' was suppressed due to non updating", chart.ID())
			chart.Obsolete()
			suppressed++
			continue
		}

		active++
		if chart.Update(data, j.sinceLast.ConvertTo(time.Microsecond)) {
			updated++

		}
	}

	j.Debugf("update charts: updated:%d, active:%d, suppressed:%d", updated, active, suppressed)
	return updated > 0
}

func (j *Job) safeGetData() (m map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC(%s)", r)
		}
	}()

	m = j.Module.GetData()
	return
}

func (j *Job) getData() map[string]int64 {
	if j.unsafe {
		return j.safeGetData()
	}
	return j.Module.GetData()
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
	return j.retries < j.RetriesMax
}
