package plugin

import (
	"testing"

	"github.com/netdata/go-orchestrator/module"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestConfig_UnmarshalYAML(t *testing.T) {
	tests := map[string]struct {
		input   string
		wantCfg config
	}{
		"valid configuration": {
			input: "enabled: yes\ndefault_run: yes\nmodules:\n  module1: yes\n  module2: yes",
			wantCfg: config{
				Enabled:    true,
				DefaultRun: true,
				Modules: map[string]bool{
					"module1": true,
					"module2": true,
				},
			},
		},
		"valid configuration with broken modules section": {
			input: "enabled: yes\ndefault_run: yes\nmodules:\nmodule1: yes\nmodule2: yes",
			wantCfg: config{
				Enabled:    true,
				DefaultRun: true,
				Modules: map[string]bool{
					"module1": true,
					"module2": true,
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var cfg config
			err := yaml.Unmarshal([]byte(test.input), &cfg)
			require.NoError(t, err)
			assert.Equal(t, test.wantCfg, cfg)
		})
	}
}

func TestPlugin_loadConfig(t *testing.T) {
	tests := map[string]struct {
		plugin  Plugin
		wantCfg config
	}{
		"valid config file": {
			plugin: Plugin{
				Name:    "plugin-valid",
				ConfDir: []string{"testdata"},
			},
			wantCfg: config{
				Enabled:    true,
				DefaultRun: true,
				MaxProcs:   1,
				Modules: map[string]bool{
					"module1": true,
					"module2": true,
				},
			},
		},
		"no config path provided": {
			plugin:  Plugin{},
			wantCfg: defaultConfig(),
		},
		"config file not found": {
			plugin: Plugin{
				Name:    "plugin",
				ConfDir: []string{"testdata/not-exist"},
			},
			wantCfg: defaultConfig(),
		},
		"empty config file": {
			plugin: Plugin{
				Name:    "plugin-empty",
				ConfDir: []string{"testdata"},
			},
			wantCfg: defaultConfig(),
		},
		"invalid syntax config file": {
			plugin: Plugin{
				Name:    "plugin-invalid-syntax",
				ConfDir: []string{"testdata"},
			},
			wantCfg: defaultConfig(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.wantCfg, test.plugin.loadPluginConfig())
		})
	}
}

func TestPlugin_loadEnabledModules(t *testing.T) {
	tests := map[string]struct {
		plugin      Plugin
		cfg         config
		wantModules module.Registry
	}{
		"load all, module disabled by default but explicitly enabled": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{
					"module1": module.Creator{Defaults: module.Defaults{Disabled: true}},
				},
			},
			cfg: config{
				Modules: map[string]bool{"module1": true},
			},
			wantModules: module.Registry{
				"module1": module.Creator{Defaults: module.Defaults{Disabled: true}},
			},
		},
		"load all, module disabled by default and not explicitly enabled": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{
					"module1": module.Creator{Defaults: module.Defaults{Disabled: true}},
				},
			},
			wantModules: module.Registry{},
		},
		"load all, module in config modules (default_run=true)": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{
					"module1": module.Creator{},
				},
			},
			cfg: config{
				Modules:    map[string]bool{"module1": true},
				DefaultRun: true,
			},
			wantModules: module.Registry{
				"module1": module.Creator{},
			},
		},
		"load all, module not in config modules (default_run=true)": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{"module1": module.Creator{}},
			},
			cfg: config{
				DefaultRun: true,
			},
			wantModules: module.Registry{"module1": module.Creator{}},
		},
		"load all, module in config modules (default_run=false)": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{
					"module1": module.Creator{},
				},
			},
			cfg: config{
				Modules: map[string]bool{"module1": true},
			},
			wantModules: module.Registry{
				"module1": module.Creator{},
			},
		},
		"load all, module not in config modules (default_run=false)": {
			plugin: Plugin{
				ModuleRegistry: module.Registry{
					"module1": module.Creator{},
				},
			},
			wantModules: module.Registry{},
		},
		"load specific, module exist in registry": {
			plugin: Plugin{
				RunModule: "module1",
				ModuleRegistry: module.Registry{
					"module1": module.Creator{},
				},
			},
			wantModules: module.Registry{
				"module1": module.Creator{},
			},
		},
		"load specific, module doesnt exist in registry": {
			plugin: Plugin{
				RunModule:      "module3",
				ModuleRegistry: module.Registry{},
			},
			wantModules: module.Registry{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.wantModules, test.plugin.loadEnabledModules(test.cfg))
		})
	}
}

// TODO: tech debt
func TestPlugin_buildDiscoveryConf(t *testing.T) {

}
