<!--
title: "Dnsmasq DHCP monitoring with Netdata"
description: "Monitor the health and performance of Dnsmasq DHCP servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/dnsmasq_dhcp/README.md"
sidebar_label: "Dnsmasq DHCP"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "References/Collectors references/Networking"
-->

# Dnsmasq DHCP monitoring with Netdata

[`Dnsmasq`](http://www.thekelleys.org.uk/dnsmasq/doc.html) is a lightweight, easy to configure, DNS forwarder and DHCP
server.

This module monitors `Dnsmasq DHCP` leases database.

## Metrics

All metrics have "dnsmasq_dhcp." prefix.

| Metric                      |   Scope    | Dimensions |   Units    |
|-----------------------------|:----------:|:----------:|:----------:|
| dhcp_ranges                 |   global   | ipv4, ipv6 |   ranges   |
| dhcp_hosts                  |   global   | ipv4, ipv6 |   hosts    |
| dhcp_range_utilization      | dhcp_range |    used    | percentage |
| dhcp_range_allocated_leases | dhcp_range | allocated  |   leases   |

## Auto-detection

Module automatically detects all configured dhcp-ranges reading `dnsmasq` configuration files.

By default it uses:

- `/var/lib/misc/dnsmasq.leases` to read leases.
- `/etc/dnsmasq.conf` to detect dhcp-ranges.
- `/etc/dnsmasq.d` to find additional configurations.

## Configuration

Edit the `go.d/dnsmasq_dhcp.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dnsmasq_dhcp.conf
```

Here is an example:

```yaml
jobs:
  - name: dnsmasq_dhcp
    leases_path: /var/lib/misc/dnsmasq.leases
    conf_path: /etc/dnsmasq.conf
    conf_dir: /etc/dnsmasq.d
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsmasq_dhcp.conf).

## Troubleshooting

To troubleshoot issues with the `dnsmasq_dhcp` collector, run the `go.d.plugin` with the debug option enabled. The
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
  ./go.d.plugin -d -m dnsmasq_dhcp
  ```
