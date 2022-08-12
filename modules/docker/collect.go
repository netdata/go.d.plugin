package docker

import (
	"context"
	"fmt"

	dockerClient "github.com/docker/docker/client"
)

func (d *Docker) collect() (map[string]int64, error) {
	if d.client == nil {
		client, err := dockerClient.NewClientWithOpts()
		if err != nil {
			return nil, err
		}
		d.client = client
	}

	defer d.client.Close()

	info, err := d.client.Info(context.Background())
	if err != nil {
		return nil, err
	}

	fmt.Println(info.Containers, info.ContainersRunning, info.ContainersPaused, info.ContainersStopped)
	return nil, nil
}
