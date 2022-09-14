<!--
title: "ProxySQL monitoring with Netdata"
description: "Monitor connections, slow queries, lagging, backends status and more with zero configuration and per-second metric granularity."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/proxysql/README.md
sidebar_label: "ProxySQL"
-->

# ProxySQL monitoring with Netdata

[`ProxySQL`](https://www.proxysql.com/) is an open-source proxy for mySQL.

This module monitors one or more `ProxySQL` servers, depending on your configuration.

## Requirements

Executed queries:

- `SELECT Variable_Name, Variable_Value FROM global_variables;`
- `SELECT Variable_Name, Variable_Value FROM stats_memory_metrics;`
- `SELECT * FROM stats_mysql_commands_counters;`
- `SELECT Variable_Name, Variable_Value FROM stats_mysql_global;`
- `SELECT * FROM stats_mysql_users;`


Netdata uses `stats` username and password which is enabled by default

## Charts

It produces the following charts:

- Uptime in `seconds`
- Questions in `questions`
- Active transactions in `transanctions`
- Slow queries in `queries`
- Backend lagging during query in `backends`
- Backend offline during query in `backends`
- Generated error packets in `packets/s`
- Max connection timeouts in `connections`
- Client connections in `connections`
- Server connections in `connections`
- Query time in `nanoseconds`
- ProxySQL commands in `commands`
- Connection  pool requests in `connections`
- Connection  pool connections in `connections`
- Mysql monitor threads in `threads`
- Mysql thread workers in `workers`
- Network in `bytes/s`
- Query cache in `number of entries`
- Prepared statements in `prepared statements`
- MySQL max allowed packet in `bytes`
- Memory in `bytes`
- Jemalloc memory in `bytes`
- MySQL command counts in `commands`
- MySQL user connections in `connections` 


## Configuration

Edit the `go.d/proxysql.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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

