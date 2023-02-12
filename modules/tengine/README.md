<!--
title: "Tengine monitoring with Netdata"
description: "Monitor the health and performance of Tengine web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/tengine/README.md"
sidebar_label: "Tengine"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitoring/Webapps"
-->

# Tengine monitoring with Netdata

[`Tengine`](https://tengine.taobao.org/) is a web server originated by Taobao, the largest e-commerce website in Asia.
It is based on the Nginx HTTP server and has many advanced features.

This module monitors one or more `Tengine` instances, depending on your configuration.

## Requirements

- `tengine` with configured [`ngx_http_reqstat_module`](http://tengine.taobao.org/document/http_reqstat.html).
- collector expects [default line format](http://tengine.taobao.org/document/http_reqstat.html).

## Metrics

All metrics have "tengine." prefix.

| Metric                                           | Scope  |                               Dimensions                               |     Units     |
|--------------------------------------------------|:------:|:----------------------------------------------------------------------:|:-------------:|
| bandwidth_total                                  | global |                                in, out                                 |      B/s      |
| connections_total                                | global |                                accepted                                | connections/s |
| requests_total                                   | global |                               processed                                |  requests/s   |
| requests_per_response_code_family_total          | global |                       2xx, 3xx, 4xx, 5xx, other                        |  requests/s   |
| requests_per_response_code_detailed_total        | global | 200, 206, 302, 304, 403, 404, 419, 499, 500, 502, 503, 504, 508, other |  requests/s   |
| requests_upstream_total                          | global |                                requests                                |  requests/s   |
| tries_upstream_total                             | global |                                 calls                                  |    calls/s    |
| requests_upstream_per_response_code_family_total | global |                                4xx, 5xx                                |  requests/s   |

## Configuration

Edit the `go.d/tengine.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/tengine.conf
```

Needs only `url` to server's `/us`. Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/us

  - name: remote
    url: http://203.0.113.10/us
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/tengine.conf).

## Troubleshooting

To troubleshoot issues with the `tengine` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m tengine
  ```
