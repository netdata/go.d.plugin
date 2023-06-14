# Chrony collector

## Overview

[Chrony](https://chrony.tuxfamily.org/) is a versatile implementation of the Network Time Protocol (NTP).

This collector monitors the system's clock performance and peers activity status using Chrony communication protocol v6.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                      |                        Dimensions                        |  Unit   |
|-----------------------------|:--------------------------------------------------------:|:-------:|
| chrony.stratum              |                         stratum                          |  level  |
| chrony.current_correction   |                    current_correction                    | seconds |
| chrony.root_delay           |                        root_delay                        | seconds |
| chrony.root_dispersion      |                        root_delay                        | seconds |
| chrony.last_offset          |                          offset                          | seconds |
| chrony.rms_offset           |                          offset                          | seconds |
| chrony.frequency            |                        frequency                         |   ppm   |
| chrony.residual_frequency   |                    residual_frequency                    |   ppm   |
| chrony.skew                 |                           skew                           |   ppm   |
| chrony.update_interval      |                     update_interval                      | seconds |
| chrony.ref_measurement_time |                   ref_measurement_time                   | seconds |
| chrony.leap_status          |   normal, insert_second, delete_second, unsynchronised   | status  |
| chrony.activity             | online, offline, burst_online, burst_offline, unresolved | sources |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/chrony.conf`.

The file format is YAML. Generally the format is:

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
sudo ./edit-config go.d/chrony.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                        |    Default    | Required |
|:-------------------:|--------------------------------------------------------------------|:-------------:|:--------:|
|    update_every     | Data collection frequency.                                         |       5       |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |       0       |          |
|       address       | Server address. The format is IP:PORT.                             | 127.0.0.1:323 |   yes    |
|       timeout       | Connection timeout. Zero means no timeout.                         |       1       |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:323
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
    address: 127.0.0.1:323

  - name: remote
    address: 192.0.2.1:323
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `chrony` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m chrony
  ```
