# Ping collector

## Overview

This module measures round-tripe time and packet loss by sending ping messages to network hosts.

There are two operational modes:

- privileged (send raw ICMP ping, default). Requires
  CAP_NET_RAW [capability](https://man7.org/linux/man-pages/man7/capabilities.7.html) or root privileges:
  > **Note**: set automatically during Netdata installation.

  ```bash
  sudo setcap CAP_NET_RAW=eip <INSTALL_PREFIX>/usr/libexec/netdata/plugins.d/go.d.plugin
  ```

- unprivileged (send UDP ping, Linux only).
  Requires configuring [ping_group_range](https://www.man7.org/linux/man-pages/man7/icmp.7.html):

  ```bash
  sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
  ```
  To persist the change add `net.ipv4.ping_group_range="0 2147483647"` to `/etc/sysctl.conf` and
  execute `sudo sysctl -p`.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### host

These metrics refer to the remote host.

Labels:

| Label | Description |
|-------|-------------|
| host  | remote host |

Metrics:

| Metric                |   Dimensions   |     Unit     |
|-----------------------|:--------------:|:------------:|
| ping.host_rtt         | min, max, avg  | milliseconds |
| ping.host_std_dev_rtt |    std_dev     | milliseconds |
| ping.host_packet_loss |      loss      |  percentage  |
| ping.host_packets     | received, sent |   packets    |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/ping.conf`.

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
sudo ./edit-config go.d/ping.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                            | Default | Required |
|:-------------------:|----------------------------------------------------------------------------------------|:-------:|:--------:|
|    update_every     | Data collection frequency.                                                             |    5    |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                     |    0    |          |
|        hosts        | Network hosts.                                                                         |         |   yes    |
|     privileged      | Ping packets type. "no" means send an "unprivileged" UDP ping,  "yes" - raw ICMP ping. |   yes   |          |
|       packets       | Number of ping packets to send.                                                        |    5    |          |
|      interval       | Timeout between sending ping packets.                                                  |  100ms  |          |

</details>

#### Examples

##### IPv4 hosts

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: example
    hosts:
      - 192.0.2.0
      - 192.0.2.1
```

</details>

##### Unprivileged mode

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: example
    privileged: no
    hosts:
      - 192.0.2.0
      - 192.0.2.1
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Multiple instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: example1
    hosts:
      - 192.0.2.0
      - 192.0.2.1

  - name: example2
    packets: 10
    hosts:
      - 192.0.2.3
      - 192.0.2.4
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `ping` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m ping
  ```
