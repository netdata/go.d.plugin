package godplugin

import (
	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func (gd *GoDPlugin) jobsStart(jobs chan *job.Job) {
	if jobs == nil {
		return
	}
	for j := range jobs {
		if v := j.Mod.UpdateEvery(); v != j.UpdateEvery {
			j.UpdateEvery = v
		}

		if v := j.Mod.ModuleName(); v != j.RealModuleName {
			logger.SetModName(j.Logger, v)
			j.RealModuleName = v
		}

		if gd.cmd.Debug || j.UpdateEvery < gd.cmd.UpdateEvery {
			j.UpdateEvery = gd.cmd.UpdateEvery
		}

		c := j.Mod.GetCharts()
		if c == nil {
			j.Error("GetCharts() [FAILED]")
			continue
		}
		j.Obs.Set(c)

		gd.wg.Add(1)
		go j.Start(&gd.wg)
	}
	gd.wg.Wait()
}
