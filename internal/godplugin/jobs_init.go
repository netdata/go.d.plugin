package godplugin

import (
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func (gd *goDPlugin) jobsInit(created jobStack) jobStack {
	var js jobStack

	if created.empty() {
		return nil
	}

	for _, j := range created {
		l := logger.New(j.RealModuleName, j.JobName())
		j.Logger = l

		j.Mod.SetUpdateEvery(j.UpdateEvery)
		j.Mod.SetModuleName(j.RealModuleName)
		j.Mod.SetLogger(l)

		j.Mod.Init()

		js.push(j)
	}

	created.destroy()
	return js
}
