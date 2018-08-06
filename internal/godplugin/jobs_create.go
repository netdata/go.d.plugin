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

func (p *Plugin) createJobs() []job.Job {
	var jobs []job.Job

	if p.Option.Module == "all" {
		for moduleName, creator := range p.registry {
			if !p.Config.IsModuleEnabled(moduleName, false) {
				log.Infof("module '%s' is disabled in configuration file", moduleName)
				continue
			}
			if creator.DisabledByDefault && !p.Config.IsModuleEnabled(moduleName, true) {
				log.Infof("module '%s' is disabled by default", moduleName)
				continue
			}
			createdJob := p.createJob(moduleName, creator, p.ModuleConfDir)
			log.Debugf("created job: %v", createdJob)
			jobs = append(jobs, createdJob...)
		}
	} else {
		if creator, ok := p.registry[p.Option.Module]; ok {
			createdJob := p.createJob(p.Option.Module, creator, p.ModuleConfDir)
			log.Debugf("created job: %v", createdJob)
			jobs = append(jobs, createdJob...)
		} else {
			showAvailableModulesInfo()
		}
	}

	return jobs
}

func (p *Plugin) createJob(moduleName string, creator modules.Creator, moduleConfDir string) []job.Job {
	log.Debugf("create jobs for module '%s'", moduleName)
	var jobs []job.Job

	conf := job.NewConfig()
	conf.RealModuleName = moduleName
	conf.UpdateEvery = p.Option.UpdateEvery
	if creator.UpdateEvery != nil {
		conf.UpdateEvery = *creator.UpdateEvery
	}
	if creator.ChartCleanup != nil {
		conf.ChartCleanup = *creator.ChartCleanup
	}

	module := creator.Create()

	if creator.NoConfig {
		jobs = append(jobs, p.newJobFunc(module, conf, p.Out))
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
		jobs = append(jobs, p.newJobFunc(module, conf, p.Out))
		return jobs
	}

	for _, r := range raw {
		c, m := *conf, creator.Create()

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
		jobs = append(jobs, p.newJobFunc(m, &c, p.Out))
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
	for v := range modules.DefaultRegistry {
		s = append(s, v)
	}
	sort.Strings(s)
	for idx, n := range s {
		fmt.Printf("  %d. %s\n", idx+1, n)
	}
}
