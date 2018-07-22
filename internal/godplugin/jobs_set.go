package godplugin

import (
	"errors"
	"reflect"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/pkg/charts/cooked"
	"github.com/l2isbad/go.d.plugin/internal/pkg/logger"
)

const (
	fieldCharts   = "Charts"
	fieldLogger   = "Logger"
	fieldBaseConf = "BaseConfHook"
)

func (gd *goDPlugin) jobsSet(created jobStack) jobStack {
	var js jobStack

	if created.empty() {
		return nil
	}

	for _, j := range created {
		err := setJobFields(j.Module, j.Config)

		if err != nil {
			log.Errorf("\"%s\" %s: %s", j.ModuleName(), j.JobName(), err)
			continue
		}

		if gd.cmd.Debug || j.UpdEvery < gd.cmd.UpdateEvery {
			j.SetUpdateEvery(gd.cmd.UpdateEvery)
		}

		js.push(j)
	}

	created.destroy()
	return js
}

func setJobFields(mod interface{}, conf *job.Config) error {
	v := reflect.ValueOf(mod)
	if v.Kind() != reflect.Ptr {
		return errors.New("module must be a pointer")
	}
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("module must be a struct")
	}
	// Mandatory field
	if !setField(&elem, conf, fieldCharts) {
		return errors.New("'Charts' field must be a 'modules.Charts' interface")
	}
	// Mandatory field
	if !setField(&elem, conf, fieldLogger) {
		return errors.New("'Logger' field must be a 'modules.Logger' interface")
	}
	// Optional field
	setField(&elem, conf, fieldBaseConf)
	return nil
}

func setField(v *reflect.Value, conf *job.Config, fieldName string) bool {
	f := v.FieldByName(fieldName)
	if f.IsValid() && f.Kind() == reflect.Interface && f.Type().Name() == fieldName && f.CanSet() {
		switch fieldName {
		case fieldCharts:
			f.Set(reflect.ValueOf(cooked.NewCharts(conf)))
		case fieldLogger:
			f.Set(reflect.ValueOf(logger.New(conf)))
		case fieldBaseConf:
			f.Set(reflect.ValueOf(conf))
		}
		return true
	}
	return false
}
