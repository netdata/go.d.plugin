package godplugin

import "time"

func (p *Plugin) jobStartLoop() {
	started := make(map[string]bool)
LOOP:
	for {
		select {
		case <-p.jobStartShutdown:
			break LOOP
		case job := <-p.jobCh:
			if started[job.FullName()] {
				log.Infof("skipping %s[%s]: already served by another job", job.ModuleName(), job.Name())
				continue
			}
			if p.initJob(job) && p.checkJob(job) && p.postCheckJob(job) {
				started[job.FullName()] = true
				go job.Start()
				p.loopQueue.add(job)
			}
		}
	}
}

func (p *Plugin) initJob(job Job) bool {
	if !job.Init() {
		log.Errorf("%s[%s] Init failed", job.ModuleName(), job.Name())
		return false
	}
	return true
}

func (p *Plugin) checkJob(job Job) bool {
	ok := job.Check()

	if job.Panicked() {
		return false
	}

	if !ok {
		log.Errorf("%s[%s] Check failed", job.ModuleName(), job.Name())
		if job.AutoDetectionRetry() > 0 {
			go recheckTask(p.jobCh, job)
		}
		return false
	}
	return true
}

func (p *Plugin) postCheckJob(job Job) bool {
	if !job.PostCheck() {
		log.Errorf("%s[%s] PostCheck failed", job.ModuleName(), job.Name())
		return false
	}
	return true
}

func recheckTask(ch chan Job, job Job) {
	log.Infof("%s[%s] scheduling next check in %d seconds",
		job.ModuleName(),
		job.Name(),
		job.AutoDetectionRetry(),
	)
	time.Sleep(time.Second * time.Duration(job.AutoDetectionRetry()))

	t := time.NewTimer(time.Second * 10)
	defer t.Stop()

	select {
	case <-t.C:
	case ch <- job:
	}
}
