package godplugin

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/netdata/go.d.plugin/modules"
)

func (p *Plugin) loadModuleConfig(name string) *moduleConfig {

	log.Infof("loading '%s' configuration", name)

	configPath, err := p.ConfigPath.Find(fmt.Sprintf("go.d/%s.conf", name))
	if err != nil {
		log.Warningf("skipping '%s': %v", name, err)
		return nil
	}

	modConf := newModuleConfig()
	modConf.name = name

	if err = loadYAML(modConf, configPath); err != nil {
		log.Warningf("skipping '%s': %v", name, err)
		return nil
	}

	if len(modConf.Jobs) == 0 {
		log.Warningf("skipping '%s': config 'jobs' section is empty or not exist", name)
		return nil
	}

	return modConf
}

func (p *Plugin) createModuleJobs(modConf *moduleConfig) []Job {
	var jobs []Job

	creator := p.registry[modConf.name]
	modConf.updateJobs(creator.UpdateEvery, p.Option.UpdateEvery)

	jobName := func(conf map[string]interface{}) interface{} {
		if name, ok := conf["name"]; ok {
			return name
		}
		return "unnamed"
	}

	for _, conf := range modConf.Jobs {
		mod := creator.Create()

		if err := unmarshalAndValidate(conf, mod); err != nil {
			log.Errorf("skipping %s[%s]: %s", modConf.name, jobName(conf), err)
			continue
		}

		job := modules.NewJob(modConf.name, mod, p.Out, p)

		if err := unmarshalAndValidate(conf, job); err != nil {
			log.Errorf("skipping %s[%s]: %s", modConf.name, jobName(conf), err)
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs
}

func (p *Plugin) createJobs() []Job {
	var jobs []Job

	for name := range p.modules {
		conf := p.loadModuleConfig(name)
		if conf == nil {
			continue
		}

		for _, job := range p.createModuleJobs(conf) {
			jobs = append(jobs, job)
		}
	}

	return jobs
}

func unmarshalAndValidate(conf interface{}, module interface{}) error {
	b, _ := yaml.Marshal(conf)
	if err := yaml.Unmarshal(b, module); err != nil {
		return err
	}
	if err := validate.Struct(module); err != nil {
		return err
	}
	return nil
}
