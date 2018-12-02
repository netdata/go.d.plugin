package godplugin

import (
	"fmt"
	"runtime"

	"github.com/netdata/go.d.plugin/modules"
)

func (p *Plugin) populateActiveModules() {
	if p.Option.Module != "all" {
		if creator, exist := modules.DefaultRegistry[p.Option.Module]; exist {
			p.modules[p.Option.Module] = creator
		}
		return
	}

	for name, creator := range modules.DefaultRegistry {
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
	name, err := p.ConfigPath.Find("go.d.conf")

	if err != nil {
		log.Critical(err)
		return false
	}

	if err := load(p.config, name); err != nil {
		log.Critical(err)
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
