package godplugin

//
//func (gd *GoDPlugin) jobsStart(jobs chan *job.Job) {
//	if jobs == nil {
//		return
//	}
//	for j := range jobs {
//		if v := j.Module.UpdateEvery(); v != j.UpdateEvery {
//			j.UpdateEvery = v
//		}
//
//		if v := j.Module.ModuleName(); v != j.RealModuleName {
//			logger.SetModName(j.Logger, v)
//			j.RealModuleName = v
//		}
//
//		if gd.cmd.Debug || j.UpdateEvery < gd.cmd.UpdateEvery {
//			j.UpdateEvery = gd.cmd.UpdateEvery
//		}
//
//		c := j.Module.GetCharts()
//		if c == nil {
//			j.Error("GetCharts() [FAILED]")
//			continue
//		}
//		j.observer.Set(c)
//
//		gd.wg.Add(1)
//		go j.Start(&gd.wg)
//	}
//	gd.wg.Wait()
//}
