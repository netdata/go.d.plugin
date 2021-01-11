<!--
title: "Dnsmasq DNS Forwarder"
description: "Monitor the health and performance of Dnsmasq DNS forwarders with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/dnsmasq/README.md
sidebar_label: "Dnsmasq DNS Forwarder"
-->

# Dnsmasq DNS Forwarder

[`Dnsmasq`](http://www.thekelleys.org.uk/dnsmasq/doc.html) is a lightweight, easy to configure DNS forwarder, designed
to provide DNS (and optionally DHCP and TFTP) services to a small-scale network.

This module monitors one or more `Dnsmasq DNS Forwarder` instances, depending on your configuration.

It collects DNS cache statistics
by [reading the response on the following query](https://manpages.debian.org/stretch/dnsmasq-base/dnsmasq.8.en.html#NOTES):

```cmd
;; opcode: QUERY, status: NOERROR, id: 37862
;; flags: rd; QUERY: 7, ANSWER: 0, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;cachesize.bind.   CH	 TXT
;insertions.bind.  CH	 TXT
;evictions.bind.   CH	 TXT
;hits.bind.        CH	 TXT
;misses.bind.      CH	 TXT
;auth.bind.        CH	 TXT
;servers.bind.     CH	 TXT
```

## Charts

- Queries forwarded to the upstream servers in `queries/s`
- Cache entries in `entries`
- Cache operations in `operations/s`
- Cache performance in `events/s`

## Configuration

Edit the `go.d/dnsmasq.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dnsmasq.conf
```

Needs only `address`, here is an example with two jobs:

```yaml
jobs:
  - name: local
    address: '127.0.0.1:53'

  - name: remote
    address: '203.0.113.0:53'
```

For all available options, see the `dnsmasq`
collector's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dnsmasq.conf).

## Troubleshooting

To troubleshoot issues with the `dnsmasq` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m dnsmasq
```
