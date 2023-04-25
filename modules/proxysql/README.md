<!--
title: "ProxySQL monitoring with Netdata"
description: "Monitor connections, slow queries, lagging, backends status and more with zero configuration and per-second metric granularity."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/proxysql/README.md"
sidebar_label: "proxysql-go.d.plugin (Recommended)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# ProxySQL collector

[ProxySQL](https://www.proxysql.com/) is an open-source proxy for mySQL.

This module monitors one or more ProxySQL servers, depending on your configuration.

## Metrics

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/proxysql/metrics.csv) for a list
of metrics.

## Configuration

> **Note**: this collector uses `stats` username and password which is enabled by default.

Edit the `go.d/proxysql.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/proxysql.conf
```

[DSN syntax in details](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

```yaml
jobs:
  - name: local
    dsn: '[username[:password]@][protocol[(address)]]/'
    # username:password@protocol(address)/
    # Examples:
    # - name: remote
    #   dsn: stats:stats@localhost/
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/proxysql.conf).

## Troubleshooting

To troubleshoot issues with the `proxysql` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m proxysql
  ```
