// SPDX-License-Identifier: GPL-3.0-or-later

package dummy

import (
	"context"
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/job/confgroup"
	"github.com/netdata/go.d.plugin/logger"
)

type Config struct {
	Registry confgroup.Registry
	Names    []string
}

func validateConfig(cfg Config) error {
	if len(cfg.Registry) == 0 {
		return errors.New("empty config registry")
	}
	if len(cfg.Names) == 0 {
		return errors.New("names not set")
	}
	return nil
}

type Discovery struct {
	*logger.Logger
	reg   confgroup.Registry
	names []string
}

func NewDiscovery(cfg Config) (*Discovery, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("config validation: %v", err)
	}
	d := &Discovery{
		reg:    cfg.Registry,
		names:  cfg.Names,
		Logger: logger.New("discovery", "dummy"),
	}
	return d, nil
}

func (d Discovery) String() string {
	return "dummy discovery"
}

func (d Discovery) Run(ctx context.Context, in chan<- []*confgroup.Group) {
	d.Info("instance is started")
	defer func() { d.Info("instance is stopped") }()

	select {
	case <-ctx.Done():
	case in <- d.groups():
	}
	close(in)
}

func (d Discovery) groups() (groups []*confgroup.Group) {
	for _, name := range d.names {
		groups = append(groups, d.newCfgGroup(name))
	}
	return groups
}

func (d Discovery) newCfgGroup(name string) *confgroup.Group {
	def, ok := d.reg.Lookup(name)
	if !ok {
		return nil
	}

	cfg := confgroup.Config{}
	cfg.SetModule(name)
	cfg.SetSource(name)
	cfg.SetProvider("dummy")
	cfg.Apply(def)

	group := &confgroup.Group{
		Configs: []confgroup.Config{cfg},
		Source:  name,
	}
	return group
}
