// SPDX-License-Identifier: GPL-3.0-or-later

package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (d *Docker) collect() (map[string]int64, error) {
	if d.client == nil {
		client, err := d.newClient(d.Config)
		if err != nil {
			return nil, err
		}
		d.client = client
	}

	defer func() { _ = d.client.Close() }()

	mx := make(map[string]int64)

	if err := d.collectInfo(mx); err != nil {
		return nil, err
	}
	if err := d.collectContainersHealth(mx); err != nil {
		return nil, err
	}

	return mx, nil
}

func (d *Docker) collectInfo(mx map[string]int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout.Duration)
	defer cancel()

	info, err := d.client.Info(ctx)
	if err != nil {
		return err
	}

	mx["running_containers"] = int64(info.ContainersRunning)
	mx["paused_containers"] = int64(info.ContainersPaused)
	mx["stopped_containers"] = int64(info.ContainersStopped)

	return nil
}

func (d *Docker) collectContainersHealth(mx map[string]int64) error {
	ctx1, cancel1 := context.WithTimeout(context.Background(), d.Timeout.Duration)
	defer cancel1()

	args := filters.NewArgs(filters.KeyValuePair{Key: "health", Value: "healthy"})
	healthy, err := d.client.ContainerList(ctx1, types.ContainerListOptions{Filters: args})
	if err != nil {
		return err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), d.Timeout.Duration)
	defer cancel2()

	args = filters.NewArgs(filters.KeyValuePair{Key: "health", Value: "unhealthy"})
	unhealthy, err := d.client.ContainerList(ctx2, types.ContainerListOptions{Filters: args})
	if err != nil {
		return err
	}

	mx["healthy_containers"] = int64(len(healthy))
	mx["unhealthy_containers"] = int64(len(unhealthy))

	return nil
}
