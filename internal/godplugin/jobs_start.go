package godplugin

import (
	"fmt"
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/modules"
)

type result struct {
	ok  bool
	err error
}

func (gd *goDPlugin) jobsRun(jobs jobStack) {
	if jobs.empty() {
		return
	}
	started := make(map[string]bool)

	for _, j := range jobs {
		key := j.FullName()

		if started[key] {
			j.Info("[DROPPED] already served by another job")
			continue
		}

		res := check(j)
		if res.err != nil {
			j.Error(res.err)
			continue
		}

		if res.ok {
			j.Info("Check() [OK]")
			if getCharts(j) {
				started[key] = true
				go j.Start(&gd.wg)
				gd.wg.Add(1)
			}
			continue
		}

		if j.AutoDetectionRetry != 0 {
			j.Warningf("Check() [RECHECK EVERY %s]", j.AutoDetectionRetry)
			started[key] = true

			recheck(j, &gd.wg)
			gd.wg.Add(1)
			continue
		}

		j.Error("Check() [FAILED]")
	}

	jobs.destroy()
	modules.Registry.Destroy()
}

func safeCheck(f func() bool) (res result) {
	defer func() {
		if r := recover(); r != nil {
			res.err = fmt.Errorf("PANIC(%s)", r)
		}
	}()
	res.ok = f()
	return
}

func check(j *job.Job) result {
	var (
		res   result
		resCh = make(chan result)
		limit = time.NewTimer(5 * time.Second)
	)

	go func() {
		resCh <- safeCheck(j.Check)
	}()

	select {
	case res = <-resCh:
	case <-limit.C:
	}

	limit.Stop()
	return res
}

func recheck(j *job.Job, wg *sync.WaitGroup) {
	var c int
	go func() {
		for {
			c++
			time.Sleep(time.Duration(j.AutoDetectionRetry) * time.Second)
			res := check(j)
			if res.err != nil {
				j.Error(res.err)
				break
			}
			if res.ok {
				j.Infof("Check() [OK] after %d rechecks", c)
				if getCharts(j) {
					go j.Start(wg)
				}
				break
			}
		}
	}()
}

func getCharts(j *job.Job) bool {
	c := j.GetCharts()
	if c == nil {
		j.Error("GetCharts() [FAILED]")
		return false
	}

	j.Obs.Set(c)
	return true
}
