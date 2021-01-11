<!--
title: "FreeRADIUS monitoring with Netdata"
description: "Monitor the health and performance of FreeRADIUS servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/freeradius/README.md
sidebar_label: "FreeRADIUS"
-->

# FreeRADIUS monitoring with Netdata

[`FreeRADIUS`](https://freeradius.org/) is a modular, high performance free RADIUS suite.

This module will monitor one or more `FreeRADIUS` servers, depending on your configuration.

## Requirements

- `FreeRADIUS` with enabled status feature.

The configuration for the status server is automatically created in the sites-available directory. By default, server is
enabled and can be queried from every client.

To enable status feature do the following:

- `cd sites-enabled`
- `ln -s ../sites-available/status status`
- restart FreeRADIUS server

## Charts

It produces following charts:

- Authentication in `pps`
- Authentication Responses in `pps`
- Bad Authentication Requests in `pps`
- Proxy Authentication in `pps`
- Proxy Authentication Responses in `pps`
- Proxy Bad Authentication Requests in `pps`
- Accounting in `pps`
- Bad Accounting Requests in `pps`
- Proxy Accounting in `pps`
- Proxy Bad Accounting Requests in `pps`

## Configuration

Edit the `go.d/freeradius.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/freeradius.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    host: 127.0.0.1

  - name: remote
    host: 203.0.113.10
    secret: secret 
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/freeradius.conf).

## Troubleshooting

To troubleshoot issues with the `freeradius` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m freeradius
```

