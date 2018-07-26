package job

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/modules"
)

func New(m modules.Module, c Config) *Job {
	_, u := m.(modules.Unsafer)

	return &Job{
		Module: m,
		C:      c,
		unsafe: u,
	}
}

type Job struct {
	modules.Module
	t       timers
	C       Config
	W       *WrappedCharts
	retries int
	unsafe  bool
}

func (j *Job) Run(wg *sync.WaitGroup) {
Done:
	for {

		sleep := j.nextIn()
		j.Debugf("sleeping for %s to reach frequency of %d sec", sleep, j.C.UpdEvery)
		time.Sleep(sleep)

		j.t.curRun = time.Now()
		if !j.t.lastRun.IsZero() {
			j.t.sinceLast.Duration = j.t.curRun.Sub(j.t.lastRun)
		}

		if ok := j.update(); ok {
			j.retries, j.t.penalty, j.t.lastRun = 0, 0, j.t.curRun
			j.t.spentOnRun.Duration = time.Since(j.t.lastRun)

		} else if !ok && !j.handleRetries() {
			j.Errorf("stopped after %d collection failures in a row", j.C.RetriesMax)
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

	for _, chart := range j.W.items {

		if chart.flags.obsoleted {
			if !chart.canBeUpdated(data) {
				suppressed++
				continue
			}
			chart.refresh()

		} else if j.C.ChartCleanup > 0 && chart.failedUpdates >= j.C.ChartCleanup {
			j.Errorf("item '%s' was suppressed due to non updating", chart.item.ID)
			chart.obsolete()
			suppressed++
			continue
		}

		active++
		if chart.update(data, j.t.sinceLast.ConvertTo(time.Microsecond)) {
			updated++

		}
	}

	j.Debugf("update charts3: updated:%d, active:%d, suppressed:%d", updated, active, suppressed)
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
	next := start.Add(time.Duration(j.C.UpdEvery) * time.Second).Add(j.t.penalty).Truncate(time.Second)
	return time.Duration(next.UnixNano() - start.UnixNano())
}

func (j *Job) handleRetries() bool {
	j.retries++

	if j.retries%5 != 0 {
		return true
	}

	j.t.penalty = time.Duration(j.retries*j.C.UpdEvery/2) * time.Second
	j.Warningf(
		"added %.0f seconds penalty after %d failed updates in a row",
		j.t.penalty.Seconds(),
		j.retries,
	)
	return j.retries < j.C.RetriesMax
}
