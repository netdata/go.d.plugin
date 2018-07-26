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
		charts := job.NewWrappedCharts(&j.C)
		j.W = charts

		err := setModuleInterfaces(j.Module, &j.C, charts)

		if err != nil {
			log.Errorf("'%s' %s: %s", j.C.ModuleName(), j.C.JobName(), err)
			continue
		}

		if gd.cmd.Debug || j.C.UpdEvery < gd.cmd.UpdateEvery {
			j.C.SetUpdateEvery(gd.cmd.UpdateEvery)
		}

		js.push(j)
	}

	created.destroy()
	return js
}

func setModuleInterfaces(mod interface{}, conf *job.Config, charts *job.WrappedCharts) error {
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
