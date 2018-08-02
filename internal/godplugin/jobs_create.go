package godplugin

import (
	"io/ioutil"
	"path"

	val "github.com/go-playground/validator"
	"github.com/go-yaml/yaml"

	"fmt"
	"sort"

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

func (p *Plugin) createJobs() jobStack {
	var jobs jobStack

	if p.Option.Module == "all" {
		for moduleName, creator := range modules.Registry {
			if p.Config.IsModuleEnabled(moduleName, false) {
				log.Infof("module \"%s\" is disabled in configuration file", moduleName)
				continue
			}
			module := creator.MakeModule()
			if module.DisabledByDefault() && !p.Config.IsModuleEnabled(moduleName, true) {
				log.Infof("module \"%s\" is disabled by default", moduleName)
				continue
			}
			jobs = append(jobs, p.createJob(moduleName, module, creator, p.ModuleConfDir)...)
		}
	} else {
		if creator, ok := modules.Registry[p.Option.Module]; ok {
			module := creator.MakeModule()
			jobs = append(jobs, p.createJob(p.Option.Module, module, creator, p.ModuleConfDir)...)
		} else {
			showAvailableModulesInfo()
		}
	}

	return jobs
}

func (p *Plugin) createJob(moduleName string, module modules.Module, creator modules.Creator, moduleConfDir string) []*job.Job {
	var jobs []*job.Job
	conf := job.NewConfig()
	conf.RealModuleName = moduleName
	conf.UpdateEvery = module.UpdateEvery()
	conf.ChartCleanup = module.DefaultChartCleanup()

	if !module.RequireConfig() {
		jobs = append(jobs, job.New(module, conf))
		return jobs
	}

	confFile := path.Join(moduleConfDir, moduleName+".conf")
	f, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Errorf("module '%s': read config file '%s' error: %v", moduleName, confFile, err)
		return jobs
	}

	if err = unmarshal(f, conf); err != nil {
		log.Errorf("module '%s', unmarshal config file error: %v", moduleName, err)
	}

	raw := parseModuleConf(f)
	num := len(raw)

	// PUSH: single job config
	if num == 0 {
		jobs = append(jobs, job.New(module, conf))
		return jobs
	}

	for _, r := range raw {
		c, m := *conf, creator.MakeModule()

		err := unmarshal(r.conf, &c)
		// SKIP: validator error
		if err != nil {
			log.Errorf("module '%s', job '%s': %s", moduleName, r.name, err)
			continue
		}

		err = unmarshal(r.conf, m)
		// SKIP: validator error
		if err != nil {
			log.Errorf("module %s, job '%s': %s", moduleName, r.name, err)
			continue
		}

		// do not add job moduleName for multi job with only 1 job
		if num > 1 {
			c.RealJobName = r.name
		}
		jobs = append(jobs, job.New(m, &c))
	}
	return jobs
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

func unmarshal(in []byte, out interface{}) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return err
	}
	if err := validate.Struct(out); err != nil {
		return err
	}
	return nil
}

func showAvailableModulesInfo() {
	fmt.Println("Available modules:")
	var s []string
	for v := range modules.Registry {
		s = append(s, v)
	}
	sort.Strings(s)
	for idx, n := range s {
		fmt.Printf("  %d. %s\n", idx+1, n)
	}
}
