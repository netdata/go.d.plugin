<!--
title: "Docker Engine monitoring with Netdata"
description: "Monitor the health and performance of the Docker Engine runtime with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/docker_engine/README.md
sidebar_label: "Docker Engine"
-->

# Docker Engine monitoring with Netdata

[`Docker Engine`](https://docs.docker.com/engine/) is the industry’s de facto container runtime that runs on various
Linux (CentOS, Debian, Fedora, Oracle Linux, RHEL, SUSE, and Ubuntu) and Windows Server operating systems.

This module will monitor one or more `Docker Engines` applications, depending on your configuration.

## Requirements

- Docker with enabled [`metric-address`](https://docs.docker.com/config/thirdparty/prometheus/)

## Charts

It produces the following charts:

- Container Actions in `actions/s`
- Container States in `containers`
- Builder Builds Fails By Reason in `fails/s`
- Health Checks in `events/s`

If Docker is running in in [Swarm mode](https://docs.docker.com/engine/swarm/) and the instance is a Swarm manager:

- Swarm Manager Leader in `bool`
- Swarm Manager Object Store in `count`
- Swarm Manager Nodes Per State in `count`
- Swarm Manager Tasks Per State in `count`

## Configuration

Edit the `go.d/docker_engine.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/docker_engine.conf
```

Needs only `url` to docker `metric-address`. Here is an example for 2 docker instances:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9323/metrics

  - name: remote
    url: http://203.0.113.10:9323/metrics
```

For all available options, please see the
module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/docker_engine.conf).

## Troubleshooting

To troubleshoot issues with the `docker_engine` collector, run the `go.d.plugin` with the debug option enabled. The
output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m docker_engine
```
