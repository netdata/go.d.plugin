// SPDX-License-Identifier: GPL-3.0-or-later

package dyncfg

import (
	"io"

	"github.com/netdata/go.d.plugin/agent/confgroup"
	"github.com/netdata/go.d.plugin/agent/functions"
	"github.com/netdata/go.d.plugin/agent/module"
)

type Config struct {
	PluginName       string
	Out              io.Writer
	FunctionRegistry FunctionRegistry
	Modules          module.Registry
	ModuleDefaults   confgroup.Registry
}

type FunctionRegistry interface {
	Register(name string, reg func(functions.Function))
}

func validateConfig(cfg Config) error {
	return nil
}
