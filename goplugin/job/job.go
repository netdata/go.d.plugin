package job

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/modules"
)

func New(m modules.Module, c *Conf) *Job {
	_, u := m.(modules.Unsafer)
	return &Job{
		Module: m,
		Conf:   c,
		unsafe: u,
	}
}

type Job struct {
	modules.Module
	modules.Logger
	*Conf
	timers
	retries int
	unsafe  bool
}

// Start runs job.update() in a for loop every UpdateEvery.
func (j *Job) Start(wg *sync.WaitGroup) {
	rc := runtimeChart{}
	rc.create(j.FullName(), j.UpdEvery)
Done:
	for {
		sleep := j.nextIn()
		j.Debugf("sleeping for %s to reach frequency of %d sec", sleep, j.UpdEvery)
		time.Sleep(sleep)
		j.curRun = time.Now()
		if !j.lastRun.IsZero() {
			j.sinceLast.Duration = j.curRun.Sub(j.lastRun)
		}
		switch j.update() {
		case true:
			j.retries, j.penalty, j.lastRun = 0, 0, j.curRun
			j.spentOnRun.Duration = time.Since(j.lastRun)
			rc.update(
				j.FullName(),
				j.sinceLast.ConvertTo(time.Microsecond),
				j.spentOnRun.ConvertTo(time.Millisecond))
		case false:
			rc.updated = false
			if !j.handleRetries() {
				j.Errorf("stopped after %d collection failures in a row", j.RetriesMax)
				break Done
			}
		}
	}
	wg.Done()
}

// SafeGetData is a wrapper around job.GetData() which invokes recover at the end.
func (j *Job) SafeGetData() (m map[string]int64) {
	defer func() {
		if r := recover(); r != nil {
			j.Errorf("PANIC(%s)", r)
		}
	}()
	m = j.Module.GetData()
	return
}

// GetData overrides modules.Module GetData.
// It runs SafeGetData or Module.GetData depends on job safe field.
func (j *Job) GetData() map[string]int64 {
	if j.unsafe {
		return j.SafeGetData()
	}
	return j.Module.GetData()
}

// update is the core function of the job. It invokes every UpdateEvery.
// In general it runs GetData (not module GetData!) and if there is no error it
// iterates over job's charts and updates it.
// Returns true if at least one chart have been updated.
func (j *Job) update() bool {
	data := j.GetData()
	if data == nil || len(data) == 0 {
		j.Debug("GetData failed")
		return false
	}
	var updated, active, suppressed int

	for _, chart := range j.GetCharts() {
		if chart.IsObsoleted() {
			suppressed++
			if !chart.CanBeUpdated(data) {
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

// nextIng finds the time of the next loop
func (j *Job) nextIn() time.Duration {
	start := time.Now()
	next := start.Add(time.Duration(j.UpdEvery) * time.Second).Add(j.penalty).Truncate(time.Second)
	return time.Duration(next.UnixNano() - start.UnixNano())
}

// handleRetries invokes if job.update() fails. It will add time penalty (slowdown) to
// the job's next run calculation every 5 fails in a row.
// TODO 5 is hardcoded. Should we move all hardcoded values to a separate file?
func (j *Job) handleRetries() bool {
	j.retries++
	if j.retries%5 == 0 {
		j.penalty = time.Duration(j.retries*j.UpdEvery/2) * time.Second
		j.Warningf("added %.0f seconds penalty after %d failed updates in a row", j.penalty.Seconds(), j.retries)
	}
	return j.retries < j.RetriesMax
}
