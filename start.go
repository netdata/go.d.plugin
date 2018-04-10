package main

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/l2isbad/toml"

	"github.com/l2isbad/go.d.plugin/cmd"
	"github.com/l2isbad/go.d.plugin/goplugin"
	"github.com/l2isbad/go.d.plugin/logger"
	"github.com/l2isbad/go.d.plugin/modules"

	_ "github.com/l2isbad/go.d.plugin/modules/all"
)

var dir string

type namer struct{}

func (*namer) ModuleName() string {
	return "plugin"
}

func (*namer) JobName() string {
	return "main"
}

var pLogger = logger.New(&namer{})

func main() {
	// Plugins dir "usr/libexec/netdata/plugins.d"
	// Plugin config dir "/etc/netdata/"
	// Modules config dir "/etc/netdata/go.d/"
	if dir = os.Getenv("NETDATA_CONFIG_DIR"); dir == "" {
		dir, _ = os.Getwd()
		dir = path.Join(dir, "/../../../../etc/netdata")
	}
	parsedCmd := cmd.Parse()

	// Initial level is INFO
	// All loggers share the severity level. Here we set DEBUG level for ALL jobs.
	if parsedCmd.Debug {
		pLogger.SetLevel(logger.DEBUG)
	}

	conf := goplugin.NewConf()
	// TODO Should we start if the configuration file exists, but do not readable?
	f, err := ioutil.ReadFile(path.Join(dir, "go.d.conf"))
	if err != nil {
		pLogger.Error(err)
	}

	// invalid TOML format = no go
	// TODO Should we continue with default options if conf file has invalid TOML format?
	if err == nil {
		if err := toml.Unmarshal(f, conf); err != nil {
			pLogger.Critical(err)
		}
	}

	if !conf.Enabled.Bool {
		pLogger.Info("disabled in configuration file")
		return
	}
	if conf.MaxProcs != 0 {
		pLogger.Warningf("Setting GOMAXPROCS to %d", conf.MaxProcs)
		runtime.GOMAXPROCS(conf.MaxProcs)
	}

	goplugin.New(parsedCmd.UpdEvery, enabledModules(conf, parsedCmd.ModRun), pLogger, dir).Run()
}

// enabledModules returns map of enabled modules (creators actually)
func enabledModules(conf *goplugin.Conf, modRun string) modules.CreatorsMap {
	rv := make(modules.CreatorsMap)
	switch modRun {
	case "all":
		for k, v := range modules.Registry {
			if isEnabled(conf, k) {
				rv[k] = v
			}
		}
		return rv
	default:
		if v, ok := modules.Registry[modRun]; !ok {
			cmd.Info()
		} else {
			rv[modRun] = v
		}
		return rv
	}
}

// isEnabled returns whether the module is enabled in plugin configuration file
func isEnabled(conf *goplugin.Conf, modName string) bool {
	if conf.DefaultRun.Bool {
		return conf.Modules[modName] == nil || conf.Modules[modName].Bool
	}
	return conf.Modules[modName] != nil && conf.Modules[modName].Bool
}
