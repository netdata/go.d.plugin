<!--
title: "Docker monitoring with Netdata"
description: "Monitor Docker containers health metrics and data usage metrics."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/docker/README.md"
sidebar_label: "Docker"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Virtualized environments/Containers"
-->

# Docker collector

[Docker Engine](https://docs.docker.com/engine/) is an open source containerization technology for building and
containerizing your applications.

This module monitors one or more Docker Engine instances, depending on your configuration.

## Metrics

- All metrics have "docker." prefix.
- container_writeable_layer_size needs `collect_container_size: yes`. Enabling this may result in high CPU usage
  depending on the version of Docker Engine.

Labels per scope:

- global: no labels.
- container: container_name, image.

| Metric                         |   Scope   |                             Dimensions                              |   Units    |
|--------------------------------|:---------:|:-------------------------------------------------------------------:|:----------:|
| containers_state               |  global   |                      running, paused, stopped                       | containers |
| containers_health_status       |  global   | healthy, unhealthy, not_running_unhealthy, starting, no_healthcheck | containers |
| images                         |  global   |                          active, dangling                           |   images   |
| images_size                    |  global   |                                size                                 |   bytes    |
| container_state                | container |    running, paused, exited, created, restarting, removing, dead     |   state    |
| container_health_status        | container | healthy, unhealthy, not_running_unhealthy, starting, no_healthcheck |   status   |
| container_writeable_layer_size | container |                           writeable_layer                           |    size    |

## Configuration

Edit the `go.d/docker.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/docker.conf
```

```yaml
jobs:
  - name: local
    address: 'unix:///var/run/docker.sock'
    collect_container_size: no

  - name: remote
    address: 'tcp://203.0.113.10:2375'
    collect_container_size: no
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/docker.conf).

## Troubleshooting

To troubleshoot issues with the `docker` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m docker
```
