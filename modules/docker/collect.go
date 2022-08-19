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

	if err := d.collectUsage(mx); err != nil {
		return nil, err
	}
	if err := d.collectContainersHealth(mx); err != nil {
		return nil, err
	}

	return mx, nil
}

func (d *Docker) collectUsage(mx map[string]int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout.Duration)
	defer cancel()

	usage, err := d.client.DiskUsage(ctx)
	if err != nil {
		return err
	}

	data := struct {
		containersRunning int64
		containersPaused  int64
		containersExited  int64
		imagesTotal       int64
		imagesDangling    int64
		imagesSize        int64
		volumesTotal      int64
		volumesDangling   int64
		volumesSize       int64
	}{}

	for _, v := range usage.Containers {
		switch v.State {
		case "running":
			data.containersRunning++
		case "exited":
			data.containersExited++
		case "paused":
			data.containersPaused++
		}
	}

	for _, v := range usage.Images {
		data.imagesTotal++
		if v.Containers == 0 {
			data.imagesDangling++
		}
		data.imagesSize += v.Size
	}

	for _, v := range usage.Volumes {
		data.volumesTotal++
		if v.UsageData == nil {
			continue
		}
		if v.UsageData.RefCount == 0 {
			data.volumesDangling++
		}
		if v.UsageData.Size != -1 {
			data.volumesSize += v.UsageData.Size
		}
	}

	mx["running_containers"] = data.containersRunning
	mx["exited_containers"] = data.containersExited
	mx["paused_containers"] = data.containersPaused
	mx["images_active"] = data.imagesTotal - data.imagesDangling
	mx["images_dangling"] = data.imagesDangling
	mx["images_size"] = data.imagesSize
	mx["volumes_active"] = data.volumesTotal - data.volumesDangling
	mx["volumes_dangling"] = data.volumesDangling
	mx["volumes_size"] = data.volumesSize

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
