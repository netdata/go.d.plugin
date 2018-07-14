package godplugin

import (
	"time"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"sync"
)

func (gd *goDPlugin) jobsRun(jobs jobStack) {
	if jobs.Empty() {
		return
	}

	for _, j := range jobs {
		ok := check(j)

		if ok {
			j.Info("Check() [OK]")
			go j.Run(&gd.wg)
			gd.wg.Add(1)
			continue
		}

		if !ok && j.AutoDetectionRetry != 0 {
			j.Warningf("Check() [RECHECK EVERY %s]", j.AutoDetectionRetry)
			recheck(j, &gd.wg)
			gd.wg.Add(1)
			continue
		}

		j.Error("Check() [FAILED]")
	}

	jobs.Destroy()
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
