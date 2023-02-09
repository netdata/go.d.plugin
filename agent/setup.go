// SPDX-License-Identifier: GPL-3.0-or-later

package agent

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"
	"github.com/netdata/go.d.plugin/agent/job/discovery"
	"github.com/netdata/go.d.plugin/agent/job/discovery/dummy"
	"github.com/netdata/go.d.plugin/agent/job/discovery/file"
	"github.com/netdata/go.d.plugin/agent/module"

	"gopkg.in/yaml.v2"
)

func defaultConfig() config {
	return config{
		Enabled:    true,
		DefaultRun: true,
		MaxProcs:   0,
		Modules:    nil,
	}
}

type config struct {
	Enabled    bool            `yaml:"enabled"`
	DefaultRun bool            `yaml:"default_run"`
	MaxProcs   int             `yaml:"max_procs"`
	Modules    map[string]bool `yaml:"modules"`
}

func (c *config) String() string {
	return fmt.Sprintf("enabled '%v', default_run '%v', max_procs '%d'",
		c.Enabled, c.DefaultRun, c.MaxProcs)
}

func (a *Agent) loadPluginConfig() config {
	a.Info("loading config file")

	if len(a.ConfDir) == 0 {
		a.Info("config dir not provided, will use defaults")
		return defaultConfig()
	}

	cfgPath := a.Name + ".conf"
	a.Infof("looking for '%s' in %v", cfgPath, a.ConfDir)

	path, err := a.ConfDir.Find(cfgPath)
	if err != nil || path == "" {
		a.Warning("couldn't find config, will use defaults")
		return defaultConfig()
	}
	a.Infof("found '%s", path)

	cfg := defaultConfig()
	if err := loadYAML(&cfg, path); err != nil {
		a.Warningf("couldn't load config '%s': %v, will use defaults", path, err)
		return defaultConfig()
	}
	a.Info("config successfully loaded")
	return cfg
}

func (a *Agent) loadEnabledModules(cfg config) module.Registry {
	a.Info("loading modules")

	all := a.RunModule == "all" || a.RunModule == ""
	enabled := module.Registry{}

	for name, creator := range a.ModuleRegistry {
		if !all && a.RunModule != name {
			continue
		}
		if all && creator.Disabled && !cfg.isExplicitlyEnabled(name) {
			a.Infof("'%s' module disabled by default, should be explicitly enabled in the config", name)
			continue
		}
		if all && !cfg.isImplicitlyEnabled(name) {
			a.Infof("'%s' module disabled in the config file", name)
			continue
		}
		enabled[name] = creator
	}
	a.Infof("enabled/registered modules: %d/%d", len(enabled), len(a.ModuleRegistry))
	return enabled
}

func (a *Agent) buildDiscoveryConf(enabled module.Registry) discovery.Config {
	a.Info("building discovery config")

	reg := confgroup.Registry{}
	for name, creator := range enabled {
		reg.Register(name, confgroup.Default{
			MinUpdateEvery:     a.MinUpdateEvery,
			UpdateEvery:        creator.UpdateEvery,
			AutoDetectionRetry: creator.AutoDetectionRetry,
			Priority:           creator.Priority,
		})
	}

	var readPaths, dummyPaths []string

	if len(a.ModulesConfDir) == 0 {
		if isInsideK8sCluster() {
			return discovery.Config{Registry: reg}
		}
		a.Info("modules conf dir not provided, will use default config for all enabled modules")
		for name := range enabled {
			dummyPaths = append(dummyPaths, name)
		}
		return discovery.Config{
			Registry: reg,
			Dummy:    dummy.Config{Names: dummyPaths},
		}
	}

	for name := range enabled {
		// TODO: properly handle module renaming
		// We need to announce this change in Netdata v1.39.0 release notes and then remove this workaround.
		// This is just a quick fix for wmi=>windows. We need to prefer user wmi.conf over windows.conf
		// 2nd part of this fix is in /agent/job/discovery/file/parse.go parseStaticFormat()
		if name == "windows" {
			cfgName := "wmi.conf"
			a.Infof("looking for '%s' in %v", cfgName, a.ModulesConfDir)

			path, err := a.ModulesConfDir.Find(cfgName)

			if err == nil && strings.Contains(path, "etc/netdata") {
				a.Infof("found '%s", path)
				readPaths = append(readPaths, path)
				continue
			}
		}

		cfgName := name + ".conf"
		a.Infof("looking for '%s' in %v", cfgName, a.ModulesConfDir)

		path, err := a.ModulesConfDir.Find(cfgName)
		if isInsideK8sCluster() {
			if err != nil {
				a.Infof("not found '%s', won't use default (reading stock configs is disabled in k8s)", cfgName)
				continue
			} else if isStockConfig(path) {
				a.Infof("found '%s', but won't load it (reading stock configs is disabled in k8s)", cfgName)
				continue
			}
		}
		if err != nil {
			a.Infof("couldn't find '%s' module config, will use default config", name)
			dummyPaths = append(dummyPaths, name)
		} else {
			a.Infof("found '%s", path)
			readPaths = append(readPaths, path)
		}
	}

	a.Infof("dummy/read/watch paths: %d/%d/%d", len(dummyPaths), len(readPaths), len(a.ModulesSDConfPath))
	return discovery.Config{
		Registry: reg,
		File: file.Config{
			Read:  readPaths,
			Watch: a.ModulesSDConfPath,
		},
		Dummy: dummy.Config{
			Names: dummyPaths,
		},
	}
}

func (c *config) isExplicitlyEnabled(moduleName string) bool {
	return c.isEnabled(moduleName, true)
}

func (c *config) isImplicitlyEnabled(moduleName string) bool {
	return c.isEnabled(moduleName, false)
}

func (c *config) isEnabled(moduleName string, explicit bool) bool {
	if enabled, ok := c.Modules[moduleName]; ok {
		return enabled
	}
	if explicit {
		return false
	}
	return c.DefaultRun
}

func (c *config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	var m map[string]interface{}
	if err := unmarshal(&m); err != nil {
		return err
	}

	for key, value := range m {
		switch key {
		case "enabled", "default_run", "max_procs", "modules":
			continue
		}
		var b bool
		if in, err := yaml.Marshal(value); err != nil || yaml.Unmarshal(in, &b) != nil {
			continue
		}
		if c.Modules == nil {
			c.Modules = make(map[string]bool)
		}
		c.Modules[key] = b
	}
	return nil
}

func loadYAML(conf interface{}, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err = yaml.NewDecoder(f).Decode(conf); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}

var (
	envKubeHost         = os.Getenv("KUBERNETES_SERVICE_HOST")
	envKubePort         = os.Getenv("KUBERNETES_SERVICE_PORT")
	envNDStockConfigDir = os.Getenv("NETDATA_STOCK_CONFIG_DIR")
)

func isInsideK8sCluster() bool { return envKubeHost != "" && envKubePort != "" }

func isStockConfig(path string) bool {
	if envNDStockConfigDir == "" {
		return false
	}
	return strings.HasPrefix(path, envNDStockConfigDir)
}
