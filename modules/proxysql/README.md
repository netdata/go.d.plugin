# ProxySQL collector

## Overview

[ProxySQL](https://www.proxysql.com/) is an open-source proxy for mySQL.

This module monitors one or more ProxySQL servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                             |                                                                                          Dimensions                                                                                           |     Unit      |
|----------------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------:|
| proxysql.client_connections_count                  |                                                                             connected, non_idle, hostgroup_locked                                                                             |  connections  |
| proxysql.client_connections_rate                   |                                                                                       created, aborted                                                                                        | connections/s |
| proxysql.server_connections_count                  |                                                                                           connected                                                                                           |  connections  |
| proxysql.server_connections_rate                   |                                                                                   created, aborted, delayed                                                                                   | connections/s |
| proxysql.backends_traffic                          |                                                                                          recv, sent                                                                                           |      B/s      |
| proxysql.clients_traffic                           |                                                                                          recv, sent                                                                                           |      B/s      |
| proxysql.active_transactions_count                 |                                                                                            client                                                                                             |  connections  |
| proxysql.questions_rate                            |                                                                                           questions                                                                                           |  questions/s  |
| proxysql.slow_queries_rate                         |                                                                                             slow                                                                                              |   queries/s   |
| proxysql.queries_rate                              | autocommit, autocommit_filtered, commit_filtered, rollback, rollback_filtered, backend_change_user, backend_init_db, backend_set_names, frontend_init_db, frontend_set_names, frontend_use_db |   queries/s   |
| proxysql.backend_statements_count                  |                                                                                         total, unique                                                                                         |  statements   |
| proxysql.backend_statements_rate                   |                                                                                    prepare, execute, close                                                                                    | statements/s  |
| proxysql.client_statements_count                   |                                                                                         total, unique                                                                                         |  statements   |
| proxysql.client_statements_rate                    |                                                                                    prepare, execute, close                                                                                    | statements/s  |
| proxysql.cached_statements_count                   |                                                                                            cached                                                                                             |  statements   |
| proxysql.query_cache_entries_count                 |                                                                                            entries                                                                                            |    entries    |
| proxysql.query_cache_memory_used                   |                                                                                             used                                                                                              |       B       |
| proxysql.query_cache_io                            |                                                                                            in, out                                                                                            |      B/s      |
| proxysql.query_cache_requests_rate                 |                                                                                   read, write, read_success                                                                                   |  requests/s   |
| proxysql.mysql_monitor_workers_count               |                                                                                      workers, auxiliary                                                                                       |    threads    |
| proxysql.mysql_monitor_workers_rate                |                                                                                            started                                                                                            |   workers/s   |
| proxysql.mysql_monitor_connect_checks_rate         |                                                                                        succeed, failed                                                                                        |   checks/s    |
| proxysql.mysql_monitor_ping_checks_rate            |                                                                                        succeed, failed                                                                                        |   checks/s    |
| proxysql.mysql_monitor_read_only_checks_rate       |                                                                                        succeed, failed                                                                                        |   checks/s    |
| proxysql.mysql_monitor_replication_lag_checks_rate |                                                                                        succeed, failed                                                                                        |   checks/s    |
| proxysql.jemalloc_memory_used                      |                                                                    active, allocated, mapped, metadata, resident, retained                                                                    |       B       |
| proxysql.memory_used                               |       auth, sqlite3, query_digest, query_rules, firewall_users_table, firewall_users_config, firewall_rules_table, firewall_rules_config, mysql_threads, admin_threads, cluster_threads       |       B       |
| proxysql.uptime                                    |                                                                                            uptime                                                                                             |    seconds    |

### command

These metrics refer to the SQL command.

Labels:

| Label   | Description  |
|---------|--------------|
| command | SQL command. |

Metrics:

| Metric                                    |                             Dimensions                              |     Unit     |
|-------------------------------------------|:-------------------------------------------------------------------:|:------------:|
| proxysql.mysql_command_execution_rate     |                               uptime                                |   seconds    |
| proxysql.mysql_command_execution_time     |                                time                                 | microseconds |
| proxysql.mysql_command_execution_duration | 100us, 500us, 1ms, 5ms, 10ms, 50ms, 100ms, 500ms, 1s, 5s, 10s, +Inf | microseconds |

### user

These metrics refer to the user.

Labels:

| Label | Description                         |
|-------|-------------------------------------|
| user  | username from the mysql_users table |

Metrics:

| Metric                                      | Dimensions |    Unit     |
|---------------------------------------------|:----------:|:-----------:|
| proxysql.mysql_user_connections_utilization |    used    | percentage  |
| proxysql.mysql_user_connections_count       |    used    | connections |

### backend

These metrics refer to the backend server.

Labels:

| Label | Description         |
|-------|---------------------|
| host  | backend server host |
| port  | backend server port |

Metrics:

| Metric                             |                 Dimensions                  |     Unit      |
|------------------------------------|:-------------------------------------------:|:-------------:|
| proxysql.backend_status            | online, shunned, offline_soft, offline_hard |    status     |
| proxysql.backend_connections_usage |                 free, used                  |  connections  |
| proxysql.backend_connections_rate  |               succeed, failed               | connections/s |
| proxysql.backend_queries_rate      |                   queries                   |   queries/s   |
| proxysql.backend_traffic           |                 recv, send                  |      B/s      |
| proxysql.backend_latency           |                   latency                   | microseconds  |

## Setup

### Prerequisites

No action required.

### Configuration

#### File

The configuration file name is `go.d/proxysql.conf`.

The file format is YAML. Generally, the format is:

```yaml
update_every: 1
autodetection_retry: 0
jobs:
  - name: some_name1
  - name: some_name1
```

You can edit the configuration file using the `edit-config` script from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md#the-netdata-config-directory).

```bash
cd /etc/netdata 2>/dev/null || cd /opt/netdata/etc/netdata
sudo ./edit-config go.d/proxysql.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                      |             Default              | Required |
|:-------------------:|--------------------------------------------------------------------------------------------------|:--------------------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                       |                5                 |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                               |                0                 |          |
|         dsn         | Data Source Name. See [DSN syntax](https://github.com/go-sql-driver/mysql#dsn-data-source-name). | stats:stats@tcp(127.0.0.1:6032)/ |   yes    |
|       my.cnf        | Specifies my.cnf file to read connection parameters from under the [client] section.             |                                  |          |
|       timeout       | Query timeout in seconds.                                                                        |                1                 |          |

</details>

#### Examples

##### TCP socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: stats:stats@tcp(127.0.0.1:6032)/
```

</details>

##### my.cnf

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    my.cnf: '/etc/my.cnf'
```

</details>

##### Multi-instance

> **Note**: When you define multiple jobs, their names must be unique.

Local and remote instances.

<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: stats:stats@tcp(127.0.0.1:6032)/

  - name: remote
    dsn: stats:stats@tcp(203.0.113.0:6032)/
```

</details>

## Troubleshooting

### Debug mode

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
