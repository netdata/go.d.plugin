<!--
title: "Docker Engine monitoring with Netdata"
description: "Monitor the health and performance of the Docker Engine runtime with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/docker_engine/README.md"
sidebar_label: "Docker Engine"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Virtualized environments/Containers"
-->

# Docker Engine collector

[`Docker Engine`](https://docs.docker.com/engine/) is the industry’s de facto container runtime that runs on various
Linux (CentOS, Debian, Fedora, Oracle Linux, RHEL, SUSE, and Ubuntu) and Windows Server operating systems.

This module will monitor one or more `Docker Engines` applications, depending on your configuration.

## Requirements

- Docker with enabled [`metric-address`](https://docs.docker.com/config/thirdparty/prometheus/)

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/docker_engine/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/docker_engine.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

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

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m docker_engine
  ```
