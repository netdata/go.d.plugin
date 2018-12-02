package godplugin

import "time"

func (p *Plugin) checkJobs() {
	started := make(map[string]bool)

	for job := range p.checkCh {
		if started[job.FullName()] {
			log.Infof("skipping %s[%s]: already served by another job", job.ModuleName(), job.Name())
			continue
		}

		if !job.Initialized() && !job.Init() {
			log.Errorf("%s[%s] Init failed", job.ModuleName(), job.Name())
			continue
		}

		ok := job.Check()

		if job.Panicked() {
			continue
		}

		if !ok {
			log.Errorf("%s[%s] Check failed", job.ModuleName(), job.Name())
			if job.AutoDetectionRetry() > 0 {
				go recheckTask(p.checkCh, job)
			}
			continue
		}

		if !job.PostCheck() {
			log.Errorf("%s[%s] PostCheck failed", job.ModuleName(), job.Name())
			continue
		}

		started[job.FullName()] = true

		log.Infof("%s[%s]: Check OK", job.ModuleName(), job.Name())

		go job.Start()
		p.loopQueue.add(job)
	}
}

func recheckTask(ch chan Job, job Job) {
	log.Infof("%s[%s] scheduling next check in %d seconds",
		job.ModuleName(),
		job.Name(),
		job.AutoDetectionRetry(),
	)
	time.Sleep(time.Second * time.Duration(job.AutoDetectionRetry()))
	ch <- job
}
