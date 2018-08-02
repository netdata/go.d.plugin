package godplugin

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func (p *Plugin) jobsInit(created jobStack) jobStack {
	var js jobStack

	if created.empty() {
		return nil
	}

	for _, job := range created {
		l := logger.New(job.RealModuleName, job.JobName())
		job.Logger = l

		job.Module.SetUpdateEvery(job.UpdateEvery)
		job.Module.SetModuleName(job.RealModuleName)
		job.Module.SetLogger(l)

		job.Module.Init()

		js.push(job)
	}

	created.destroy()
	return js
}
