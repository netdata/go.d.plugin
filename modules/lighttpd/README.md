<!--
title: "Lighttpd monitoring with Netdata"
description: "Monitor the health and performance of Lighttpd web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/lighttpd/README.md"
sidebar_label: "Lighttpd"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# Lighttpd collector

[`Lighttpd`](https://www.lighttpd.net/) is an open-source web server optimized for speed-critical environments while
remaining standards-compliant, secure and flexible

This module will monitor one or more `Lighttpd` servers, depending on your configuration.

## Requirements

- `lighttpd` with enabled [`mod_status`](https://redmine.lighttpd.net/projects/lighttpd/wiki/Docs_ModStatus).

## Metrics

All metrics have "lighttpd." prefix.

| Metric     | Scope  |                                                    Dimensions                                                    |    Units    |
|------------|:------:|:----------------------------------------------------------------------------------------------------------------:|:-----------:|
| requests   | global |                                                     requests                                                     | requests/s  |
| net        | global |                                                       sent                                                       | kilobits/s  |
| workers    | global |                                                    idle, busy                                                    |   servers   |
| scoreboard | global | waiting, open, close, hard_error, keepalive, read, read_post, write, handle_request, request_start, request_end, | connections |
| uptime     | global |                                                      uptime                                                      |   seconds   |

## Configuration

Edit the `go.d/lighttpd.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/lighttpd.conf
```

Needs only `url` to server's `server-status?auto`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/server-status?auto

  - name: remote
    url: http://203.0.113.10/server-status?auto
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/lighttpd.conf).

## Troubleshooting

To troubleshoot issues with the `lighttpd` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m lighttpd
  ```
