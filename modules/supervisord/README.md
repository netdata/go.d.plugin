# Supervisord collector

## Overview

[Supervisor](http://supervisord.org/) is a client/server system that allows its users to monitor and control a number of
processes on UNIX-like operating systems.

This collector monitors one or more Supervisor instances, depending on your configuration.

It can collect metrics from
both [unix socket](http://supervisord.org/configuration.html?highlight=unix_http_server#unix-http-server-section-values)
and [internal http server](http://supervisord.org/configuration.html?highlight=unix_http_server#inet-http-server-section-settings)

Used methods:

- [`supervisor.getAllProcessInfo`](http://supervisord.org/api.html#supervisor.rpcinterface.SupervisorNamespaceRPCInterface.getAllProcessInfo)

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                        |      Dimensions      |   Unit    |
|-------------------------------|:--------------------:|:---------:|
| supervisord.summary_processes | running, non-running | processes |

### process group

These metrics refer to the process group.

This scope has no labels.

Metrics:

| Metric                          |       Dimensions        |    Unit     |
|---------------------------------|:-----------------------:|:-----------:|
| supervisord.processes           |  running, non-running   |  processes  |
| supervisord.process_state_code  | a dimension per process |    code     |
| supervisord.process_exit_status | a dimension per process | exit status |
| supervisord.process_uptime      | a dimension per process |   seconds   |
| supervisord.process_downtime    | a dimension per process |   seconds   |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/supervisord.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/supervisord.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                        |          Default           | Required |
|:-------------------:|--------------------------------------------------------------------|:--------------------------:|:--------:|
|    update_every     | Data collection frequency.                                         |             1              |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |             0              |          |
|         url         | Server URL.                                                        | http://127.0.0.1:9001/RPC2 |   yes    |
|       timeout       | System bus requests timeout.                                       |             1              |          |

</details>

#### Examples

##### HTTP

Collect metrics via HTTP.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: 'http://127.0.0.1:9001/RPC2'
```

</details>

##### Socket

Collect metrics via Unix socket.
<details>
<summary>Config</summary>

```yaml
- name: local
  url: 'unix:///run/supervisor.sock'
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collect metrics from local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    url: 'http://127.0.0.1:9001/RPC2'

  - name: remote
    url: 'http://192.0.2.1:9001/RPC2'
```

</details>

## Troubleshooting

### Debug mode

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
