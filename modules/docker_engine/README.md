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

All metrics have "docker_engine." prefix.

| Metric                                    | Scope  |                                                                                                         Dimensions                                                                                                          |   Units    |
|-------------------------------------------|:------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:----------:|
| engine_daemon_container_actions           | global |                                                                                           changes, commit, create, delete, start                                                                                            | actions/s  |
| engine_daemon_container_states_containers | global |                                                                                                  running, paused, stopped                                                                                                   | containers |
| builder_builds_failed_total               | global | build_canceled, build_target_not_reachable_error, command_not_supported_error, dockerfile_empty_error, dockerfile_syntax_error, error_processing_commands_error, missing_onbuild_arguments_error, unknown_instruction_error |  fails/s   |
| engine_daemon_health_checks_failed_total  | global |                                                                                                            fails                                                                                                            |  events/s  |
| swarm_manager_leader                      | global |                                                                                                          is_leader                                                                                                          |    bool    |
| swarm_manager_object_store                | global |                                                                                     nodes, services, tasks, networks, secrets, configs                                                                                      |  objects   |
| swarm_manager_nodes_per_state             | global |                                                                                             ready, down, unknown, disconnected                                                                                              |   nodes    |
| swarm_manager_tasks_per_state             | global |                                                running, failed, ready, rejected, starting, shutdown, new, orphaned, preparing, pending, complete, remove, accepted, assigned                                                |   tasks    |

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
