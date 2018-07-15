package godplugin

import (
	"sync"
	"time"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/modules"
)

func (gd *goDPlugin) jobsRun(jobs jobStack) {
	if jobs.Empty() {
		return
	}
	started := make(map[string]bool)

	for _, j := range jobs {
		key := j.GetFullName()

		if started[key] {
			j.Info("[DROPPED] already served by another job")
			continue
		}

		ok := check(j)

		if ok {
			j.Info("Check() [OK]")
			started[key] = true

			go j.Run(&gd.wg)
			gd.wg.Add(1)
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

	jobs.Destroy()
	modules.Registry.Destroy()
}

func check(j *job.Job) bool {
	var (
		check = make(chan bool)
		limit = time.After(5 * time.Second)
		ok    bool
	)
	go func() {
		check <- j.Check()
	}()

	select {
	case ok = <-check:
	case <-limit:
	}

	return ok
}

func recheck(j *job.Job, wg *sync.WaitGroup) {
	var c int
	go func() {
		for {
			c++
			time.Sleep(time.Duration(j.AutoDetectionRetry) * time.Second)
			if j.Check() {
				j.Infof("Check() [OK] after %d rechecks", c)
				j.Run(wg)
			}
		}
	}()
}
