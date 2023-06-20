# MySQL collector

## Overview

[MySQL](https://www.mysql.com/) is an open-source relational database management system.

This collector monitors one or more MySQL servers, depending on your configuration.

Executed queries:

- `SELECT VERSION();`
- `SHOW GLOBAL STATUS;`
- `SHOW GLOBAL VARIABLES;`
- `SHOW SLAVE STATUS;` or `SHOW ALL SLAVES STATUS;` (MariaDBv10.2+)
- `SHOW USER_STATISTICS;` (MariaDBv10.1.1+)
- `SELECT TIME,USER FROM INFORMATION_SCHEMA.PROCESSLIST;`

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                    |                                                                     Dimensions                                                                      |      Unit      |
|-------------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| mysql.net                                 |                                                                       in, out                                                                       |   kilobits/s   |
| mysql.queries                             |                                                          queries, questions, slow_queries                                                           |   queries/s    |
| mysql.queries_type                        |                                                       select, delete, update, insert, replace                                                       |   queries/s    |
| mysql.handlers                            | commit, delete, prepare, read_first, read_key, read_next, read_prev, read_rnd, read_rnd_next, rollback, savepoint, savepointrollback, update, write |   handlers/s   |
| mysql.table_open_cache_overflows          |                                                                     open_cache                                                                      |  overflows/s   |
| mysql.table_locks                         |                                                                  immediate, waited                                                                  |    locks/s     |
| mysql.join_issues                         |                                                full_join, full_range_join, range, range_check, scan                                                 |    joins/s     |
| mysql.sort_issues                         |                                                              merge_passes, range, scan                                                              |    issues/s    |
| mysql.tmp                                 |                                                             disk_tables, files, tables                                                              |    events/s    |
| mysql.connections                         |                                                                    all, aborted                                                                     | connections/s  |
| mysql.connections_active                  |                                                              active, limit, max_active                                                              |  connections   |
| mysql.threads                             |                                                             connected, cached, running                                                              |    threads     |
| mysql.threads_created                     |                                                                       created                                                                       |   threads/s    |
| mysql.thread_cache_misses                 |                                                                       misses                                                                        |     misses     |
| mysql.innodb_io                           |                                                                     read, write                                                                     |     KiB/s      |
| mysql.innodb_io_ops                       |                                                                reads, writes, fsyncs                                                                |  operations/s  |
| mysql.innodb_io_pending_ops               |                                                                reads, writes, fsyncs                                                                |   operations   |
| mysql.innodb_log                          |                                                            waits, write_requests, writes                                                            |  operations/s  |
| mysql.innodb_cur_row_lock                 |                                                                    current waits                                                                    |   operations   |
| mysql.innodb_rows                         |                                                          inserted, read, updated, deleted                                                           |  operations/s  |
| mysql.innodb_buffer_pool_pages            |                                                           data, dirty, free, misc, total                                                            |     pages      |
| mysql.innodb_buffer_pool_pages_flushed    |                                                                     flush_pages                                                                     |   requests/s   |
| mysql.innodb_buffer_pool_bytes            |                                                                     data, dirty                                                                     |      MiB       |
| mysql.innodb_buffer_pool_read_ahead       |                                                                    all, evicted                                                                     |    pages/s     |
| mysql.innodb_buffer_pool_read_ahead_rnd   |                                                                     read-ahead                                                                      |  operations/s  |
| mysql.innodb_buffer_pool_ops              |                                                                disk_reads, wait_free                                                                |  operations/s  |
| mysql.innodb_os_log                       |                                                                   fsyncs, writes                                                                    |   operations   |
| mysql.innodb_os_log_fsync_writes          |                                                                       fsyncs                                                                        |  operations/s  |
| mysql.innodb_os_log_io                    |                                                                        write                                                                        |     KiB/s      |
| mysql.innodb_deadlocks                    |                                                                      deadlocks                                                                      |  operations/s  |
| mysql.files                               |                                                                        files                                                                        |     files      |
| mysql.files_rate                          |                                                                        files                                                                        |    files/s     |
| mysql.connection_errors                   |                                                  accept, internal, max, peer_addr, select, tcpwrap                                                  |    errors/s    |
| mysql.opened_tables                       |                                                                       tables                                                                        |    tables/s    |
| mysql.open_tables                         |                                                                    cache, tables                                                                    |     tables     |
| mysql.process_list_fetch_query_duration   |                                                                      duration                                                                       |  milliseconds  |
| mysql.process_list_queries_count          |                                                                    system, user                                                                     |    queries     |
| mysql.process_list_longest_query_duration |                                                                      duration                                                                       |    seconds     |
| mysql.qcache_ops                          |                                                      hits, lowmem_prunes, inserts, not_cached                                                       |   queries/s    |
| mysql.qcache                              |                                                                       queries                                                                       |    queries     |
| mysql.qcache_freemem                      |                                                                        free                                                                         |      MiB       |
| mysql.qcache_memblocks                    |                                                                     free, total                                                                     |     blocks     |
| mysql.galera_writesets                    |                                                                       rx, tx                                                                        |  writesets/s   |
| mysql.galera_bytes                        |                                                                       rx, tx                                                                        |     KiB/s      |
| mysql.galera_queue                        |                                                                       rx, tx                                                                        |   writesets    |
| mysql.galera_conflicts                    |                                                                bf_aborts, cert_fails                                                                |  transactions  |
| mysql.galera_flow_control                 |                                                                       paused                                                                        |       ms       |
| mysql.galera_cluster_status               |                                                         primary, non_primary, disconnected                                                          |     status     |
| mysql.galera_cluster_state                |                                                  undefined, joining, donor, joined, synced, error                                                   |     state      |
| mysql.galera_cluster_size                 |                                                                        nodes                                                                        |     nodes      |
| mysql.galera_cluster_weight               |                                                                       weight                                                                        |     weight     |
| mysql.galera_connected                    |                                                                      connected                                                                      |    boolean     |
| mysql.galera_ready                        |                                                                        ready                                                                        |    boolean     |
| mysql.galera_open_transactions            |                                                                        open                                                                         |  transactions  |
| mysql.galera_thread_count                 |                                                                       threads                                                                       |    threads     |
| mysql.key_blocks                          |                                                              unused, used, not_flushed                                                              |     blocks     |
| mysql.key_requests                        |                                                                    reads, writes                                                                    |   requests/s   |
| mysql.key_disk_ops                        |                                                                    reads, writes                                                                    |  operations/s  |
| mysql.binlog_cache                        |                                                                      disk, all                                                                      | transactions/s |
| mysql.binlog_stmt_cache                   |                                                                      disk, all                                                                      |  statements/s  |

### connection

These metrics refer to the replication connection.

This scope has no labels.

Metrics:

| Metric             |       Dimensions        |  Unit   |
|--------------------|:-----------------------:|:-------:|
| mysql.slave_behind |         seconds         | seconds |
| mysql.slave_status | sql_running, io_running | boolean |

### user

These metrics refer to the MySQL user.

Labels:

| Label | Description |
|-------|-------------|
| user  | username    |

Metrics:

| Metric                               |               Dimensions               |      Unit      |
|--------------------------------------|:--------------------------------------:|:--------------:|
| mysql.userstats_cpu                  |                  used                  |   percentage   |
| mysql.userstats_rows                 | read, sent, updated, inserted, deleted |  operations/s  |
| mysql.userstats_commands             |         select, update, other          |   commands/s   |
| mysql.userstats_denied_commands      |                 denied                 |   commands/s   |
| mysql.userstats_created_transactions |            commit, rollback            | transactions/s |
| mysql.userstats_binlog_written       |                written                 |      B/s       |
| mysql.userstats_empty_queries        |                 empty                  |   queries/s    |
| mysql.userstats_connections          |                created                 | connections/s  |
| mysql.userstats_lost_connections     |                  lost                  | connections/s  |
| mysql.userstats_denied_connections   |                 denied                 | connections/s  |

## Setup

### Prerequisites

#### Create netdata user

A user account should have the
following [permissions](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html):

- [`USAGE`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_usage)
- [`REPLICATION CLIENT`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_replication-client)
- [`PROCESS`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_process)

To create the `netdata` user with these permissions, execute the following in the MySQL shell:

```mysql
CREATE USER 'netdata'@'localhost';
GRANT USAGE, REPLICATION CLIENT, PROCESS ON *.* TO 'netdata'@'localhost';
FLUSH PRIVILEGES;
```

The `netdata` user will have the ability to connect to the MySQL server on localhost without a password. It will only
be able to gather statistics without being able to alter or affect operations in any way.

### Configuration

#### File

The configuration file name is `go.d/mysql.conf`.

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
sudo ./edit-config go.d/mysql.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                         |          Default          | Required |
|:-------------------:|---------------------------------------------------------------------------------------------------------------------|:-------------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                          |             5             |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                                  |             0             |          |
|         dsn         | MySQL server DSN (Data Source Name). See [DSN syntax](https://github.com/go-sql-driver/mysql#dsn-data-source-name). | root@tcp(localhost:3306)/ |   yes    |
|       my.cnf        | Specifies the my.cnf file to read the connection settings from the [client] section.                                |                           |          |
|       timeout       | Query timeout in seconds.                                                                                           |             1             |          |

</details>

#### Examples

##### TCP socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: netdata@tcp(127.0.0.1:3306)/
```

</details>

##### Unix socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: netdata@unix(/var/lib/mysql/mysql.sock)/
```

</details>

##### Connection with password

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: netdata:password@tcp(127.0.0.1:3306)/
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
    dsn: netdata@tcp(127.0.0.1:3306)/

  - name: remote
    dsn: netdata:password@tcp(203.0.113.0:3306)/
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `mysql` collector, run the `go.d.plugin` with the debug option enabled.
The output should give you clues as to why the collector isn't working.

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
  ./go.d.plugin -d -m mysql
  ```
