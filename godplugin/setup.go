package godplugin

import (
	"fmt"
	"runtime"
)

func (p *Plugin) populateActiveModules() {
	if p.Option.Module != "all" {
		if creator, exist := p.registry[p.Option.Module]; exist {
			p.modules[p.Option.Module] = creator
		}
		return
	}

	for name, creator := range p.registry {
		if creator.DisabledByDefault && !p.config.isModuleEnabled(name, true) {
			log.Infof("'%s' disabled by default", name)
			continue
		}
		if !p.config.isModuleEnabled(name, false) {
			log.Infof("'%s' disabled in configuration file", name)
			continue
		}
		p.modules[name] = creator
	}
}

func (p *Plugin) Setup() bool {
	configFile, err := p.ConfigPath.Find(p.confName)

	log.Debug("plugin config file: ", configFile)

	if err != nil {
		log.Critical("find config file error: ", err)
		return false
	}

<<<<<<< HEAD
	if err := loadYAML(p.config, configFile); err != nil {
		log.Critical("loadYAML config error: ", err)
=======
	if err = load(p.config, name); err != nil {
		log.Critical("load config error: ", err)
>>>>>>> master
		return false
	}

	if !p.config.Enabled {
		_, _ = fmt.Fprintln(p.Out, "DISABLE")
		log.Info("disabled in configuration file")
		return false
	}

	p.populateActiveModules()

	if len(p.modules) == 0 {
		log.Critical("no modules to run")
		return false
	}

	if p.config.MaxProcs > 0 {
		log.Infof("setting maximum number of used CPUs to %d", p.config.MaxProcs)
		runtime.GOMAXPROCS(p.config.MaxProcs)
	} else {
		log.Infof("maximum number of used CPUs %d", runtime.NumCPU())
	}

	log.Infof("minimum update every %d", p.Option.UpdateEvery)

	return true
}
