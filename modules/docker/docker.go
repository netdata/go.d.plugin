// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import (
	"context"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func init() {
	module.Register("docker", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Docker {
	return &Docker{
		Config: Config{
			Address: docker.DefaultDockerHost,
			Timeout: web.Duration{Duration: time.Second * 1},
		},
		charts: charts.Copy(),
		newClient: func(cfg Config) (dockerClient, error) {
			return docker.NewClientWithOpts(docker.WithHost(cfg.Address))
		},
	}
}

type Config struct {
	Timeout web.Duration `yaml:"timeout"`
	Address string       `yaml:"address"`
}

type (
	Docker struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		newClient func(Config) (dockerClient, error)
		client    dockerClient
	}
	dockerClient interface {
		DiskUsage(ctx context.Context) (types.DiskUsage, error)
		ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error)
		Close() error
	}
)

func (d *Docker) Init() bool {
	return true
}

func (d *Docker) Check() bool {
	return len(d.Collect()) > 0
}

func (d *Docker) Charts() *module.Charts {
	return d.charts
}

func (d *Docker) Collect() map[string]int64 {
	mx, err := d.collect()
	if err != nil {
		d.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (d *Docker) Cleanup() {
	if d.client == nil {
		return
	}
	if err := d.client.Close(); err != nil {
		d.Warningf("error on closing docker client: %v", err)
	}
	d.client = nil
}
