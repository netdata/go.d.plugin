package godplugin

import (
	"fmt"
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

		ok, err := check(j)
		if err != nil {
			j.Error(err)
			continue
		}

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

func check(j *job.Job) (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("PANIC(%s)", r)
		}
	}()
	ok = j.Check()
	return
}

func recheck(j *job.Job, wg *sync.WaitGroup) {
	var c int
	go func() {
		for {
			c++
			time.Sleep(time.Duration(j.AutoDetectionRetry) * time.Second)
			ok, err := check(j)

			if err != nil {
				j.Error(err)
				break
			}

			if ok {
				j.Infof("Check() [OK] after %d rechecks", c)
				go j.Run(wg)
				break
			}
		}
	}()
}
