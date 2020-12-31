<!--
title: "PHP-FPM monitoring with Netdata"
description: "Monitor the health and performance of PHP-FPM instances with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/phpfpm/README.md
sidebar_label: "PHP-FPM"
-->

# PHP-FPM monitoring with Netdata

[`PHP-FPM`](https://php-fpm.org/) is an alternative PHP FastCGI implementation with some additional features useful for
sites of any size, especially busier sites.

This module will monitor one or more `php-fpm` instances, depending on your configuration.

## Requirements

- `php-fpm` with enabled `status` page
- access to `status` page via web server

## Charts

It produces following charts:

- Active Connections in `connections`
- Requests in `requests/s`
- Performance in `status`
- Requests Duration Among All Idle Processes in `milliseconds`
- Last Request CPU Usage Among All Idle Processes in `percentage`
- Last Request Memory Usage Among All Idle Processes in `KB`

## Configuration

Edit the `go.d/phpfpm.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/phpfpm.conf
```

Needs only `url` or `socket`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://localhost/status?full&json

  - name: local
    url: http://[::1]/status?full&json

  - name: local_socket
    socket: '/tmp/php-fpm.sock'

  - name: remote
    url: http://203.0.113.10/status?full&json
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/phpfpm.conf).

## Troubleshooting

To troubleshoot issues with the `phpfpm` collector, run the `go.d.plugin` with the debug option enabled. The output
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
./go.d.plugin -d -m phpfpm
```
