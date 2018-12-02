package godplugin

import (
	"fmt"
	"gopkg.in/yaml.v2"

	"github.com/netdata/go.d.plugin/modules"
)

func (p *Plugin) createJobs() []Job {
	var jobs []Job

	for name, creator := range p.modules {
		log.Infof("loading '%s' configuration", name)

		configPath, err := p.ConfigPath.Find(fmt.Sprintf("go.d/%s.conf", name))
		if err != nil {
			log.Warningf("skipping '%s': %v", name, err)
			continue
		}

		modConf := newModuleConfig()

		if err = load(modConf, configPath); err != nil {
			log.Warningf("skipping '%s': %v", name, err)
			continue
		}

		if len(modConf.Jobs) == 0 {
			log.Warningf("skipping '%s': config 'jobs' section is empty or not exist", name)
			continue
		}

		modConf.updateJobs(creator.UpdateEvery, p.Option.UpdateEvery)

		for _, conf := range modConf.Jobs {
			mod := creator.Create()

			if err := unmarshalAndValidate(conf, mod); err != nil {
				log.Errorf("skipping %s[%s]: %s", name, jobName(conf), err)
				continue
			}

			job := modules.NewJob(name, mod, p.Out, p)

			if err := unmarshalAndValidate(conf, job); err != nil {
				log.Errorf("skipping %s[%s]: %s", name, jobName(conf), err)
				continue
			}

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

func jobName(conf map[string]interface{}) interface{} {
	if name := conf["name"]; name != nil {
		return name
	}
	return "unnamed"
}
