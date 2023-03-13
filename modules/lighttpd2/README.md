<!--
title: "Lighttpd2 monitoring with Netdata"
description: "Monitor the health and performance of Lighttpd2 web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/lighttpd2/README.md"
sidebar_label: "Lighttpd2"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# Lighttpd2 collector

[`Lighttpd2`](https://redmine.lighttpd.net/projects/lighttpd2) is a work in progress version of open-source web server.

This module will monitor one or more `Lighttpd2` servers, depending on your configuration.

## Requirements

- `lighttpd2` with enabled [`mod_status`](https://doc.lighttpd.net/lighttpd2/mod_status.html)

## Metrics

All metrics have "lighttpd2." prefix.

| Metric            | Scope  |                               Dimensions                                |    Units    |
|-------------------|:------:|:-----------------------------------------------------------------------:|:-----------:|
| requests          | global |                                requests                                 | requests/s  |
| status_codes      | global |                         1xx, 2xx, 3xx, 4xx, 5xx                         | requests/s  |
| traffic           | global |                                 in, out                                 | kilobits/s  |
| connections       | global |                               connections                               | connections |
| connection_states | global | start, read_header, handle_request, write_response, keepalive, upgraded |    state    |
| memory_usage      | global |                                  usage                                  |     KiB     |
| uptime            | global |                                 uptime                                  |   seconds   |

## Configuration

Edit the `go.d/lighttpd2.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/lighttpd2.conf
```

Needs only `url` to server's `server-status?format=plain`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/server-status?format=plain

  - name: remote
    url: http://203.0.113.10/server-status?format=plain
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/lighttpd2.conf).

## Troubleshooting

To troubleshoot issues with the `lighttpd2` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m lighttpd2
  ```
