<!--
title: "Docker Hub repository monitoring with Netdata"
description: "Monitor the health and performance of Docker Hub repositories with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/dockerhub/README.md"
sidebar_label: "Docker Hub repositories"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# Docker Hub repository collector

[`Docker Hub`](https://docs.docker.com/docker-hub/) is a service provided by Docker for finding and sharing container
images with your team.

This module will collect `Docker Hub` repositories statistics.

## Metrics

All metrics have "docker_engine." prefix.

| Metric       | Scope  |            Dimensions             |  Units  |
|--------------|:------:|:---------------------------------:|:-------:|
| pulls_sum    | global |                sum                |  pulls  |
| pulls        | global | <i>a dimension per repository</i> |  pulls  |
| pulls_rate   | global | <i>a dimension per repository</i> | pulls/s |
| stars        | global | <i>a dimension per repository</i> |  stars  |
| status       | global | <i>a dimension per repository</i> | status  |
| last_updated | global | <i>a dimension per repository</i> | seconds |

## Configuration

Edit the `go.d/dockerhub.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dockerhub.conf
```

Needs only list of `repositories`. Here is an example:

```yaml
jobs:
  - name: me
    repositories:
      - 'me/repo1'
      - 'me/repo2'
      - 'me/repo3' 
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dockerhub.conf).

## Troubleshooting

To troubleshoot issues with the `dockerhub` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m dockerhub
  ```

