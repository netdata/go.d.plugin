// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import (
	"github.com/netdata/go.d.plugin/agent/module"

	dockerClient "github.com/docker/docker/client"
)

func init() {
	module.Register("docker", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *Docker {
	return &Docker{}
}

type Config struct {
}

type Docker struct {
	module.Base
	Config `yaml:",inline"`

	client *dockerClient.Client
}

func (d *Docker) Init() bool {
	return true
}

func (d *Docker) Check() bool {
	return len(d.Collect()) > 0
}

func (d *Docker) Charts() *module.Charts {
	return nil
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
