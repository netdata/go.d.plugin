<!--
title: "NGINX monitoring"
description: "Monitor the health and performance of NGINX web servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/nginx/README.md"
sidebar_label: "NGINX"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Webapps"
-->

# NGINX collector

[`NGINX`](https://www.nginx.com/) is a web server which can also be used as a reverse proxy, load balancer, mail proxy
and HTTP cache.

This module will monitor one or more `NGINX` servers, depending on your configuration.

## Requirements

- `NGINX` with
  configured [`ngx_http_stub_status_module`](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html).

## Metrics

All metrics have "nginx." prefix.

| Metric                       | Scope  |       Dimensions       |     Units     |
|------------------------------|:------:|:----------------------:|:-------------:|
| connections                  | global |         active         |  connections  |
| connections_status           | global | reading, writing, idle |  connections  |
| connections_accepted_handled | global |   accepted, handled    | connections/s |
| requests                     | global |        requests        |  requests/s   |

## Configuration

Edit the `go.d/nginx.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/nginx.conf
```

Needs only `url` to server's `stub_status`. Here is an example for local and remote servers:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1/stub_status

  - name: remote
    url: http://203.0.113.10/stub_status
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginx.conf).

## Troubleshooting

To troubleshoot issues with the `nginx` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m nginx
  ```
