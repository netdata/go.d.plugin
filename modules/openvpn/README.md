# OpenVPN collector

## Overview

[OpenVPN](https://openvpn.net/) is an open-source commercial software that implements virtual private network
techniques to create secure point-to-point or site-to-site connections in routed or bridged configurations and remote
access facilities.

This collector monitors one or more OpenVPN servers, depending on your configuration.

It uses OpenVPN [Management Interface](https://openvpn.net/community-resources/management-interface/) to collect
metrics.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                 | Dimensions |    Unit    |
|------------------------|:----------:|:----------:|
| openvpn.active_clients |  clients   |  clients   |
| openvpn.total_traffic  |  in, out   | kilobits/s |

### user

These metrics refer to the VPN user.

Labels:

| Label    | Description  |
|----------|--------------|
| username | VPN username |

Metrics:

| Metric                       | Dimensions |    Unit    |
|------------------------------|:----------:|:----------:|
| openvpn.user_traffic         |  in, out   | kilobits/s |
| openvpn.user_connection_time |    time    |  seconds   |

## Setup

### Prerequisites

#### Enable in go.d.conf.

This collector is disabled by default. You need to explicitly enable it
in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

From the documentation for the OpenVPN Management Interface:
> Currently, the OpenVPN daemon can at most support a single management client any one time.

It is disabled to not break other tools which use `Management Interface`.

### Configuration

#### File

The configuration file name is `go.d/openvpn.conf`.

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
sudo ./edit-config go.d/openvpn.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                       |    Default     | Required |
|:-------------------:|-----------------------------------------------------------------------------------|:--------------:|:--------:|
|    update_every     | Data collection frequency.                                                        |       1        |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                |       0        |          |
|       address       | Server address in IP:PORT format.                                                 | 127.0.0.1:7505 |   yes    |
|   per_user_stats    | User selector. Determines which user metrics will be collected.                   |                |          |
|   connect_timeout   | Connection timeout in seconds. The timeout includes name resolution, if required. |       2        |          |
|    read_timeout     | Read timeout in seconds. Sets deadline for read calls.                            |       2        |          |
|    write_timeout    | Write timeout in seconds. Sets deadline for write calls.                          |       2        |          |

</details>

##### per_user_stats

Metrics of users matching the selector will be collected.

- Logic: (pattern1 OR pattern2) AND !(pattern3 or pattern4)
- Pattern syntax: [matcher](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#supported-format).
- Syntax:
  ```yaml
  per_user_stats:
    includes:
      - pattern1
      - pattern2
    excludes:
      - pattern3
      - pattern4
  ```

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:7505
```

</details>

##### With user metrics

Collect metrics of all users.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:7505
    per_user_stats:
      includes:
        - "* *"
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Collecting metrics from local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:7505

  - name: remote
    address: 203.0.113.0:7505
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `openvpn` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m openvpn
  ```
