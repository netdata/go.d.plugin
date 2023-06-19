# NTP daemon collector

## Overview

The NTPd (Network Time Protocol daemon) is an operating system program that maintains the system time in synchronization
with time-servers using the Network Time Protocol.

This collector monitors the system variables of the local `ntpd` daemon (optional incl. variables of the polled peers)
using the NTP Control Message Protocol via UDP socket, similar to `ntpq`,
the [standard NTP query program](https://doc.ntp.org/current-stable/ntpq.html).

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric             |    Dimensions    |     Unit     |
|--------------------|:----------------:|:------------:|
| ntpd.sys_offset    |      offset      | milliseconds |
| ntpd.sys_jitter    |  system, clock   | milliseconds |
| ntpd.sys_frequency |    frequency     |     ppm      |
| ntpd.sys_wander    |      clock       |     ppm      |
| ntpd.sys_rootdelay |      delay       | milliseconds |
| ntpd.sys_rootdisp  |    dispersion    | milliseconds |
| ntpd.sys_stratum   |     stratum      |   stratum    |
| ntpd.sys_tc        | current, minimum |     log2     |
| ntpd.sys_precision |    precision     |     log2     |

### peer

These metrics refer to the NTPd peer.

Labels:

| Label        | Description              |
|--------------|--------------------------|
| peer_address | peer's source IP address |

Metrics:

| Metric               | Dimensions |     Unit     |
|----------------------|:----------:|:------------:|
| ntpd.peer_offset     |   offset   | milliseconds |
| ntpd.peer_delay      |   delay    | milliseconds |
| ntpd.peer_dispersion | dispersion | milliseconds |
| ntpd.peer_jitter     |   jitter   | milliseconds |
| ntpd.peer_xleave     |   xleave   | milliseconds |
| ntpd.peer_rootdelay  | rootdelay  | milliseconds |
| ntpd.peer_rootdisp   | dispersion | milliseconds |
| ntpd.peer_stratum    |  stratum   |   stratum    |
| ntpd.peer_hmode      |   hmode    |    hmode     |
| ntpd.peer_pmode      |   pmode    |    pmode     |
| ntpd.peer_hpoll      |   hpoll    |     log2     |
| ntpd.peer_ppoll      |   ppoll    |     log2     |
| ntpd.peer_precision  | precision  |     log2     |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/ntpd.conf`.

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
sudo ./edit-config go.d/ntpd.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                        |    Default    | Required |
|:-------------------:|--------------------------------------------------------------------|:-------------:|:--------:|
|    update_every     | Data collection frequency.                                         |       1       |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |       0       |          |
|       address       | Server address in IP:PORT format.                                  | 127.0.0.1:123 |   yes    |
|       timeout       | Connection/read/write timeout.                                     |       3       |          |
|    collect_peers    | Determines whether peer metrics will be collected.                 |      no       |          |

</details>

#### Examples

##### Basic

A basic example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:123
```

</details>

##### With peers metrics

Collect peers metrics.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    address: 127.0.0.1:123
    collect_peers: yes
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
    address: 127.0.0.1:123

  - name: remote
    address: 203.0.113.0:123
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `ntpd` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m ntpd
  ```
