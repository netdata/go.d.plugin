package godplugin

import (
	"errors"
	"reflect"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

func (gd *goDPlugin) jobsSet(created jobStack) jobStack {
	var js jobStack

	if created.empty() {
		return nil
	}

	for _, j := range created {
		charts := job.NewCharts(&j.Config)
		j.Charts = charts

		err := setModuleInterfaces(j.Module, &j.Config, charts)

		if err != nil {
			log.Errorf("'%s' %s: %s", j.Config.ModuleName(), j.Config.JobName(), err)
			continue
		}

		if gd.cmd.Debug || j.Config.UpdEvery < gd.cmd.UpdateEvery {
			j.Config.SetUpdateEvery(gd.cmd.UpdateEvery)
		}

		js.push(j)
	}

	created.destroy()
	return js
}

func setModuleInterfaces(mod interface{}, conf *job.Config, charts *job.Charts) error {
	m := reflect.ValueOf(mod)
	if m.Kind() != reflect.Ptr {
		return errors.New("module must be a pointer")
	}
	elem := m.Elem()

	// MANDATORY
	f := elem.FieldByName("Charts")
	if !valid(f) {
		return errors.New("'Charts' field must be a 'modules.Charts' interface")
	}
	f.Set(reflect.ValueOf(charts))

	// MANDATORY
	f = elem.FieldByName("Logger")
	if !valid(f) {
		return errors.New("'Logger' field must be a 'modules.Logger' interface")
	}
	f.Set(reflect.ValueOf(logger.New(conf)))

	// OPTIONAL
	f = elem.FieldByName("BaseConfHook")
	if valid(f) {
		f.Set(reflect.ValueOf(conf))
	}

	return nil
}

func valid(v reflect.Value) bool {
	return v.IsValid() && v.Kind() == reflect.Interface && v.CanSet()
}
