# Dnsmasq DHCP collector

## Overview

[Dnsmasq](https://www.thekelleys.org.uk/dnsmasq/doc.html) is a lightweight, easy to configure DNS forwarder, designed to
provide DNS (and optionally DHCP and TFTP) services to a small-scale network.

This collector monitors one or more Dnsmasq DHCP leases databases, depending on your configuration.

By default, it uses:

- `/var/lib/misc/dnsmasq.leases` to read leases.
- `/etc/dnsmasq.conf` to detect dhcp-ranges.
- `/etc/dnsmasq.d` to find additional configurations.

All configured dhcp-ranges are detected automatically.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                   | Dimensions |  Unit  |
|--------------------------|:----------:|:------:|
| dnsmasq_dhcp.dhcp_ranges | ipv4, ipv6 | ranges |
| dnsmasq_dhcp.dhcp_hosts  | ipv4, ipv6 | hosts  |

### dhcp range

These metrics refer to the DHCP range.

Labels:

| Label      | Description                            |
|------------|----------------------------------------|
| dhcp_range | DHCP range in `START_IP:END_IP` format |

Metrics:

| Metric                                   | Dimensions |    Unit    |
|------------------------------------------|:----------:|:----------:|
| dnsmasq_dhcp.dhcp_range_utilization      |    used    | percentage |
| dnsmasq_dhcp.dhcp_range_allocated_leases | allocated  |   leases   |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/dnsmasq_dhcp.conf`.

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
sudo ./edit-config go.d/dnsmasq_dhcp.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                        |                    Default                    | Required |
|:-------------------:|--------------------------------------------------------------------|:---------------------------------------------:|:--------:|
|    update_every     | Data collection frequency.                                         |                       1                       |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check. |                       0                       |          |
|     leases_path     | Path to dnsmasq DHCP leases file.                                  |         /var/lib/misc/dnsmasq.leases          |          |
|      conf_path      | Path to dnsmasq configuration file.                                |               /etc/dnsmasq.conf               |          |
|      conf_dir       | Path to dnsmasq configuration directory.                           | /etc/dnsmasq.d,.dpkg-dist,.dpkg-old,.dpkg-new |          |

</details>

#### Examples

##### Basic

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: dnsmasq_dhcp
    leases_path: /var/lib/misc/dnsmasq.leases
    conf_path: /etc/dnsmasq.conf
    conf_dir: /etc/dnsmasq.d
```

</details>

##### Pi-hole

Dnsmasq DHCP on Pi-hole.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: dnsmasq_dhcp
    leases_path: /etc/pihole/dhcp.leases
    conf_path: /etc/dnsmasq.conf
    conf_dir: /etc/dnsmasq.d
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `dnsmasq_dhcp` collector, run the `go.d.plugin` with the debug option enabled. The
output
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
  ./go.d.plugin -d -m dnsmasq_dhcp
  ```
