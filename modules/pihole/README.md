<!--
title: "Pi-hole monitoring with Netdata"
description: "Monitor the health and performance of Pi-hole instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/pihole/README.md"
sidebar_label: "Pi-hole"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Apps"
-->

# Pi-hole collector

[`Pi-hole`](https://pi-hole.net) is a Linux network-level advertisement and Internet tracker blocking application which
acts as a DNS sinkhole, intended for use on a private network.

This module will monitor one or more `Pi-hole` instances using [PHP API](https://github.com/pi-hole/AdminLTE).

The data provided by the API is for the last 24 hours. All collected values refer to this time period and not to the
module's collection interval.

## Metrics

All metrics have "pihole." prefix.

| Metric                            | Scope  |            Dimensions            |   Units    |
|-----------------------------------|:------:|:--------------------------------:|:----------:|
| dns_queries_total                 | global |             queries              |  queries   |
| dns_queries                       | global |    cached, blocked, forwarded    |  queries   |
| dns_queries_percentage            | global |    cached, blocked, forwarded    | percentage |
| unique_clients                    | global |              unique              |  clients   |
| domains_on_blocklist              | global |            blocklist             |  domains   |
| blocklist_last_update             | global |               ago                |  seconds   |
| unwanted_domains_blocking_status  | global |        enabled, disabled         |   status   |
| dns_queries_types                 | global | a, aaaa, any, ptr, soa, srv, txt | percentage |
| dns_queries_forwarded_destination | global |      cached, blocked, other      | percentage |

## Configuration

Edit the `go.d/pihole.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/pihole.conf
```

Module automatically detects `Pihole` web password reading `setupVars.conf` file. It expects to find the file in
the `/etc/pihole/` directory.

If you want to monitor remote instance you need to set the password in the module configuration file.

Here is an example for local and remote instances:

```yaml
jobs:
  - name: local
    top_clients_entries: 10
    top_items_entries: 10  # top permitted and top blocked domains charts

  - name: remote
    url: http://203.0.113.10
    password: 1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b

  - name: remote_https
    url: https://203.0.113.11
    password: 1ebd33f882f9aa5fac26a7cb74704742f91100228eb322e41b7bd6e6aeb8f74b
    tls_skip_verify: yes  # self signed certificate verification skip

```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pihole.conf).

## Troubleshooting

To troubleshoot issues with the `pihole` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m pihole
```
