package godplugin

import (
	"io/ioutil"
	"path"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/modules"
	"github.com/l2isbad/yaml"
)

type jobRawConf struct {
	name string
	conf []byte
}

type jobStack []*job.Job

func (js *jobStack) Push(v *job.Job) {
	*js = append(*js, v)
}

func (js *jobStack) Empty() bool {
	return len(*js) == 0
}

func (js *jobStack) Destroy() {
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
	conf, mod := job.NewConf(), creator.MakeModule()

	conf.SetModuleName(name)
	setModuleDefaults(name, conf)

	f, err := ioutil.ReadFile(path.Join(dir, name+".conf"))

	// SKIP: config read error and not NoConfiger
	_, ok := mod.(modules.NoConfiger)
	if !ok && err != nil {
		log.Errorf("\"%s\" skipped: %s", name, err)
		return
	}

	// PUSH: job with default config (1 job module)
	if err != nil {
		log.Debug(err)
		jobs.Push(job.New(mod, conf))
		return
	}

	log.Debugf("module \"%s\" configuration read success", name)

	err = yaml.Unmarshal(f, conf)

	// SKIP: YAML parse error = no go
	if err != nil {
		log.Errorf("module \"%s\" config yaml parse: %s", name, err)
		return
	}

	for _, r := range parseModuleConf(f) {
		c, m := *conf, creator.MakeModule()

		err := yaml.Unmarshal(r.conf, &c)
		// SKIP: job
		if err != nil {
			log.Errorf("module %s, job \"%s\" skipped: %s", name, r.name, err)
			continue
		}

		err = yaml.Unmarshal(r.conf, m)
		// SKIP: job
		if err != nil {
			log.Errorf("module %s, job \"%s\" skipped: %s", name, r.name, err)
			continue
		}

		c.SetJobName(r.name)
		// PUSH: job
		jobs.Push(job.New(m, &c))
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
