package godplugin

import (
	"reflect"

	"github.com/l2isbad/go.d.plugin/internal/pkg/charts"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func (gd *goDPlugin) jobsSet(created jobStack) jobStack {
	var js jobStack

	if created.empty() {
		return nil
	}

	for _, j := range created {
		m := reflect.ValueOf(j.Module)

		if m.Kind() != reflect.Ptr {
			log.Errorf("module '%s' must be a pointer", j.ModuleName())
			continue
		}

		jobCh := charts.New()
		jobLog := logger.New(j.Config)

		j.Logger = jobLog
		j.Obs.Set(jobCh)

		elem := m.Elem()

		if f := elem.FieldByName("Charts"); !valid(f) {
			log.Errorf("module '%s': 'Charts' field must be a 'modules.Charts' interface", j.ModuleName())
			continue
		} else {
			f.Set(reflect.ValueOf(jobCh))
		}

		if f := elem.FieldByName("Logger"); valid(f) {
			f.Set(reflect.ValueOf(jobLog))
		}

		if f := elem.FieldByName("BaseConfHook"); valid(f) {
			f.Set(reflect.ValueOf(j.Config))
		}



		if gd.cmd.Debug || j.UpdEvery < gd.cmd.UpdateEvery {
			j.SetUpdateEvery(gd.cmd.UpdateEvery)
		}

		js.push(j)
	}

	created.destroy()
	return js
}

func valid(v reflect.Value) bool {
	return v.IsValid() && v.Kind() == reflect.Interface && v.CanSet()
}
