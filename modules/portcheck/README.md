# TCP endpoint collector

## Overview

This collector monitors one or more TCP services availability and response time.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### tcp endpoint

These metrics refer to the TCP endpoint.

Labels:

| Label | Description |
|-------|-------------|
| host  | host        |
| port  | port        |

Metrics:

| Metric                   |        Dimensions        |  Unit   |
|--------------------------|:------------------------:|:-------:|
| portcheck.status         | success, failed, timeout | boolean |
| portcheck.state_duration |           time           | seconds |
| portcheck.latency        |           time           |   ms    |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/portcheck.conf`.

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
sudo ./edit-config go.d/portcheck.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                        | Default | Required |
|:-------------------:|--------------------------------------------------------------------|:-------:|:--------:|
|    update_every     | Data collection frequency.                                         |    5    |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |    0    |          |
|        host         | Remote host address in IPv4, IPv6 format, or DNS name.             |         |   yes    |
|        ports        | Remote host ports. Must be specified in numeric format.            |         |   yes    |
|       timeout       | HTTP request timeout.                                              |    2    |          |

</details>

#### Examples

##### Check SSH and telnet

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: server1
    host: 127.0.0.1
    ports:
      - 22
      - 23
```

</details>

##### Check webserver with IPv6 address

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: server2
    host: "[2001:DB8::1]"
    ports:
      - 80
      - 8080
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Multiple instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: server1
    host: 127.0.0.1
    ports:
      - 22
      - 23

  - name: server2
    host: 203.0.113.10
    ports:
      - 22
      - 23
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `portcheck` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m portcheck
  ```
