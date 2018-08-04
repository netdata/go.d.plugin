package godplugin

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/modules"
)

type checkResult struct {
	ok  bool
	err error
}

func (p *Plugin) checkJobs(jobs jobStack) chan *job.Job {
	if jobs.empty() {
		return nil
	}

	var (
		wg      = new(sync.WaitGroup)
		toStart = make(chan *job.Job)
		started = make(map[string]bool)
	)

	go func() {
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
				started[key] = true

				toStart <- j
				continue
			}

			if j.AutoDetectionRetry != 0 {
				j.Warningf("Check() [RECHECK EVERY %s]", j.AutoDetectionRetry)
				started[key] = true

				go recheck(j, wg, toStart)
				wg.Add(1)
				continue
			}
			j.Error("Check() [FAILED]")
		}

		jobs.destroy()
		modules.Registry.Destroy()

		wg.Wait()
		close(toStart)
	}()

	return toStart
}

//func safeCheck(f func() bool) (res checkResult) {
//	defer func() {
//		if r := recover(); r != nil {
//			res.err = fmt.Errorf("PANIC(%v)", r)
//		}
//	}()
//	res.ok = f()
//	return
//}
//
//func check(j *job.Job) checkResult {
//	resCh := make(chan checkResult)
//
//	go func() {
//		resCh <- safeCheck(j.Module.Check)
//	}()
//
//	select {
//	case res := <-resCh:
//		return res
//	case <-time.After(5 * time.Second):
//		return checkResult{err: errors.New("check timeout")}
//	}
//}

func recheck(j *job.Job, wg *sync.WaitGroup, ch chan *job.Job) {
	c := 0
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
			ch <- j
			wg.Done()
			break
		}
	}
}
