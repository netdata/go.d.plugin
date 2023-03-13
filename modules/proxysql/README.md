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

All metrics have "proxysql." prefix.

Labels per scope:

- global: no labels.
- command: command.
- user: user.
- backend: host, port.

| Metric                                    |  Scope  |                                                                                          Dimensions                                                                                           |     Units     |
|-------------------------------------------|:-------:|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------:|
| client_connections_count                  | global  |                                                                             connected, non_idle, hostgroup_locked                                                                             |  connections  |
| client_connections_rate                   | global  |                                                                                       created, aborted                                                                                        | connections/s |
| server_connections_count                  | global  |                                                                                           connected                                                                                           |  connections  |
| server_connections_rate                   | global  |                                                                                   created, aborted, delayed                                                                                   | connections/s |
| backends_traffic                          | global  |                                                                                          recv, sent                                                                                           |      B/s      |
| clients_traffic                           | global  |                                                                                          recv, sent                                                                                           |      B/s      |
| active_transactions_count                 | global  |                                                                                            client                                                                                             |  connections  |
| questions_rate                            | global  |                                                                                           questions                                                                                           |  questions/s  |
| slow_queries_rate                         | global  |                                                                                             slow                                                                                              |   queries/s   |
| queries_rate                              | global  | autocommit, autocommit_filtered, commit_filtered, rollback, rollback_filtered, backend_change_user, backend_init_db, backend_set_names, frontend_init_db, frontend_set_names, frontend_use_db |   queries/s   |
| backend_statements_count                  | global  |                                                                                         total, unique                                                                                         |  statements   |
| backend_statements_rate                   | global  |                                                                                    prepare, execute, close                                                                                    | statements/s  |
| client_statements_count                   | global  |                                                                                         total, unique                                                                                         |  statements   |
| client_statements_rate                    | global  |                                                                                    prepare, execute, close                                                                                    | statements/s  |
| cached_statements_count                   | global  |                                                                                            cached                                                                                             |  statements   |
| query_cache_entries_count                 | global  |                                                                                            entries                                                                                            |    entries    |
| query_cache_memory_used                   | global  |                                                                                             used                                                                                              |       B       |
| query_cache_io                            | global  |                                                                                            in, out                                                                                            |      B/s      |
| query_cache_requests_rate                 | global  |                                                                                   read, write, read_success                                                                                   |  requests/s   |
| mysql_monitor_workers_count               | global  |                                                                                      workers, auxiliary                                                                                       |    threads    |
| mysql_monitor_workers_rate                | global  |                                                                                            started                                                                                            |   workers/s   |
| mysql_monitor_connect_checks_rate         | global  |                                                                                        succeed, failed                                                                                        |   checks/s    |
| mysql_monitor_ping_checks_rate            | global  |                                                                                        succeed, failed                                                                                        |   checks/s    |
| mysql_monitor_read_only_checks_rate       | global  |                                                                                        succeed, failed                                                                                        |   checks/s    |
| mysql_monitor_replication_lag_checks_rate | global  |                                                                                        succeed, failed                                                                                        |   checks/s    |
| jemalloc_memory_used                      | global  |                                                                    active, allocated, mapped, metadata, resident, retained                                                                    |       B       |
| memory_used                               | global  |       auth, sqlite3, query_digest, query_rules, firewall_users_table, firewall_users_config, firewall_rules_table, firewall_rules_config, mysql_threads, admin_threads, cluster_threads       |       B       |
| uptime                                    | global  |                                                                                            in, out                                                                                            |      B/s      |
| mysql_command_execution_rate              | command |                                                                                            uptime                                                                                             |    seconds    |
| mysql_command_execution_time              | command |                                                                                             time                                                                                              | microseconds  |
| mysql_command_execution_duration          | command |                                                              100us, 500us, 1ms, 5ms, 10ms, 50ms, 100ms, 500ms, 1s, 5s, 10s, +Inf                                                              | microseconds  |
| mysql_user_connections_utilization        |  user   |                                                                                             used                                                                                              |  percentage   |
| mysql_user_connections_count              |  user   |                                                                                             used                                                                                              |  connections  |
| backend_status                            | backend |                                                                          online, shunned, offline_soft, offline_hard                                                                          |    status     |
| backend_connections_usage                 | backend |                                                                                          free, used                                                                                           |  connections  |
| backend_connections_rate                  | backend |                                                                                        succeed, failed                                                                                        | connections/s |
| backend_queries_rate                      | backend |                                                                                            queries                                                                                            |   queries/s   |
| backend_traffic                           | backend |                                                                                          recv, send                                                                                           |      B/s      |
| backend_latency                           | backend |                                                                                            latency                                                                                            | microseconds  |

## Configuration

> **Note**: this collector uses `stats` username and password which is enabled by default.

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
