package godplugin

import (
	"io/ioutil"
	"path"

	val "github.com/go-playground/validator"
	"gopkg.in/yaml.v2"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/modules"
)

var validate = val.New()

type jobRawConf struct {
	name string
	conf []byte
}

type jobStack []*job.Job

func (js *jobStack) push(v *job.Job) {
	*js = append(*js, v)
}

func (js jobStack) empty() bool {
	return len(js) == 0
}

func (js *jobStack) destroy() {
	if !js.empty() {
		for i := range *js {
			(*js)[i] = nil
		}
	}
	*js = nil
}

func (gd *goDPlugin) jobsCreate() jobStack {
	var jobs jobStack

	switch gd.cmd.Module {
	default:
		if c, ok := modules.Registry[gd.cmd.Module]; ok {
			create(gd.cmd.Module, c, gd.dir.modulesConf, &jobs)
		} else {
			info()
		}
	case "all":
		for name, creator := range modules.Registry {

			if modules.GetDefault(name).DisabledByDefault() && !gd.conf.Modules[name] {
				log.Infof("module \"%s\" disabled by default", name)
				continue
			}

			if !isModuleEnabled(gd.conf, name) {
				log.Infof("module \"%s\" disabled in configuration file", name)
				continue
			}

			create(name, creator, gd.dir.modulesConf, &jobs)
		}
	}

	return jobs
}

func create(name string, creator modules.Creator, dir string, jobs *jobStack) {
	// Create module and default conf
	conf, mod := job.NewConfig(), creator.MakeModule()

	conf.RealModuleName = name

	setModuleDefaults(name, conf)

	f, err := ioutil.ReadFile(path.Join(dir, name+".conf"))

	// SKIP: config read error and not NoConfiger
	_, ok := mod.(modules.NoConfiger)
	if !ok && err != nil {
		log.Errorf("'%s' skipped: %s", name, err)
		return
	}

	// PUSH: jobs without configuration (only base conf)
	if err != nil {
		log.Debug(err)
		jobs.push(job.New(mod, conf))
		return
	}

	log.Debugf("module '%s' configuration read success", name)

	err = unmarshal(f, conf)

	// SKIP: YAML parse error || validator error
	if err != nil {
		log.Errorf("module '%s': %s", name, err)
		return
	}

	raw := parseModuleConf(f)
	num := len(raw)

	// PUSH: single job config
	if num == 0 {
		jobs.push(job.New(mod, conf))
		return
	}

	for _, r := range raw {
		c, m := *conf, creator.MakeModule()

		err := unmarshal(r.conf, &c)
		// SKIP: validator error
		if err != nil {
			log.Errorf("module '%s', job '%s': %s", name, r.name, err)
			continue
		}

		err = unmarshal(r.conf, m)
		// SKIP: validator error
		if err != nil {
			log.Errorf("module %s, job '%s': %s", name, r.name, err)
			continue
		}

		// do not add job name for multi job with only 1 job
		if num > 1 {
			c.RealJobName = r.name
		}
		// PUSH:
		jobs.push(job.New(m, &c))
	}

}

func isModuleEnabled(c config, n string) bool {
	v, ok := c.Modules[n]

	if c.DefaultRun {
		return !ok || ok && v
	}
	return ok && v
}

func parseModuleConf(f []byte) []jobRawConf {
	var (
		rv []jobRawConf
		m  yaml.MapSlice
	)

	yaml.Unmarshal(f, &m)

	// TODO: All maps of maps are considered as jobs. looks fragile.
	for k := range m {
		if _, ok := m[k].Value.(yaml.MapSlice); !ok {
			continue
		}
		b, _ := yaml.Marshal(m[k].Value)
		rv = append(rv, jobRawConf{m[k].Key.(string), b})
	}

	return rv
}

func setModuleDefaults(n string, c *job.Config) {
	if v, ok := modules.GetDefault(n).UpdateEvery(); ok {
		c.UpdateEvery = v
	}
	if v, ok := modules.GetDefault(n).ChartsCleanup(); ok {
		c.ChartCleanup = v
	}
}

func unmarshal(in []byte, out interface{}) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return err
	}
	if err := validate.Struct(out); err != nil {
		return err
	}
	return nil
}
