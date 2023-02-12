<!--
title: "Supervisord monitoring with Netdata"
description: "Monitor the processes running by Supervisor with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/supervisord/README.md"
sidebar_label: "Supervisord"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Virtualized environments/Virtualize hosts"
-->

# Supervisord monitoring with Netdata

[Supervisor](http://supervisord.org/) is a client/server system that allows its users to monitor and control a number of
processes on UNIX-like operating systems.

This module monitors one or more Supervisor instances, depending on your configuration.

It can collect metrics from
both [unix socket](http://supervisord.org/configuration.html?highlight=unix_http_server#unix-http-server-section-values)
and [internal http server](http://supervisord.org/configuration.html?highlight=unix_http_server#inet-http-server-section-settings)

Used methods:

- [`supervisor.getAllProcessInfo`](http://supervisord.org/api.html#supervisor.rpcinterface.SupervisorNamespaceRPCInterface.getAllProcessInfo)

## Metrics

All metrics have "supervisord." prefix.

| Metric              |     Scope     |           Dimensions           |    Units    |
|---------------------|:-------------:|:------------------------------:|:-----------:|
| summary_processes   |    global     |      running, non-running      |  processes  |
| processes           | process group |      running, non-running      |  processes  |
| process_state_code  | process group | <i>a dimension per process</i> |    code     |
| process_exit_status | process group | <i>a dimension per process</i> | exit status |
| process_uptime      | process group | <i>a dimension per process</i> |   seconds   |
| process_downtime    | process group | <i>a dimension per process</i> |   seconds   |

## Configuration

Edit the `go.d/supervisord.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/supervisord.conf
```

Endpoints can be both local or remote as long as they expose their metrics on the provided URL.

Here is an example with two endpoints:

```yaml
jobs:
  # via [unix_http_server]
  - name: local
    url: 'unix:///run/supervisor.sock'

  # via [inet_http_server]
  - name: local
    url: 'http://127.0.0.1:9001/RPC2'
```

For all available options, see the `supervisord`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/supervisord.conf).

## Troubleshooting

To troubleshoot issues with the `supervisord` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m supervisord
  ```
