<!--
title: "MySQL monitoring with Netdata"
description: "Monitor connections, slow queries, InnoDB memory and disk utilization, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/mysql/README.md"
sidebar_label: "MySQL"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# MySQL monitoring with Netdata

[MySQL](https://www.mysql.com/) is an open-source relational database management system.

This module monitors one or more MySQL servers, depending on your configuration.

## Requirements

Executed queries:

- `SELECT VERSION();`
- `SHOW GLOBAL STATUS;`
- `SHOW GLOBAL VARIABLES;`
- `SHOW SLAVE STATUS;` or `SHOW ALL SLAVES STATUS;` (MariaDBv10.2+)
- `SHOW USER_STATISTICS;` (MariaDBv10.1.1+)
- `SELECT TIME,USER FROM INFORMATION_SCHEMA.PROCESSLIST;`

[User Statistics](https://mariadb.com/kb/en/user-statistics/) query is [MariaDB](https://mariadb.com/) specific.

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

## Metrics

All metrics have "mysql." prefix.

- userstats_* metrics need [User Statistics](https://mariadb.com/kb/en/user-statistics/#enabling-the-plugin) plugin
  enabled. MariaDB and Percona MySQL only.

Labels per scope:

- global: no labels.
- connection: no labels.
- user: user.

| Metric                              |   Scope    |                                                                     Dimensions                                                                      |     Units      |
|-------------------------------------|:----------:|:---------------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| net                                 |   global   |                                                                       in, out                                                                       |   kilobits/s   |
| queries                             |   global   |                                                          queries, questions, slow_queries                                                           |   queries/s    |
| queries_type                        |   global   |                                                       select, delete, update, insert, replace                                                       |   queries/s    |
| handlers                            |   global   | commit, delete, prepare, read_first, read_key, read_next, read_prev, read_rnd, read_rnd_next, rollback, savepoint, savepointrollback, update, write |   handlers/s   |
| table_open_cache_overflows          |   global   |                                                                     open_cache                                                                      |  overflows/s   |
| table_locks                         |   global   |                                                                  immediate, waited                                                                  |    locks/s     |
| join_issues                         |   global   |                                                full_join, full_range_join, range, range_check, scan                                                 |    joins/s     |
| sort_issues                         |   global   |                                                              merge_passes, range, scan                                                              |    issues/s    |
| tmp                                 |   global   |                                                             disk_tables, files, tables                                                              |    events/s    |
| connections                         |   global   |                                                                    all, aborted                                                                     | connections/s  |
| connections_active                  |   global   |                                                              active, limit, max_active                                                              |  connections   |
| threads                             |   global   |                                                             connected, cached, running                                                              |    threads     |
| threads_created                     |   global   |                                                                       created                                                                       |   threads/s    |
| thread_cache_misses                 |   global   |                                                                       misses                                                                        |     misses     |
| innodb_io                           |   global   |                                                                     read, write                                                                     |     KiB/s      |
| innodb_io_ops                       |   global   |                                                                reads, writes, fsyncs                                                                |  operations/s  |
| innodb_io_pending_ops               |   global   |                                                                reads, writes, fsyncs                                                                |   operations   |
| innodb_log                          |   global   |                                                            waits, write_requests, writes                                                            |  operations/s  |
| innodb_cur_row_lock                 |   global   |                                                                    current waits                                                                    |   operations   |
| innodb_rows                         |   global   |                                                          inserted, read, updated, deleted                                                           |  operations/s  |
| innodb_buffer_pool_pages            |   global   |                                                           data, dirty, free, misc, total                                                            |     pages      |
| innodb_buffer_pool_pages_flushed    |   global   |                                                                     flush_pages                                                                     |   requests/s   |
| innodb_buffer_pool_bytes            |   global   |                                                                     data, dirty                                                                     |      MiB       |
| innodb_buffer_pool_read_ahead       |   global   |                                                                    all, evicted                                                                     |    pages/s     |
| innodb_buffer_pool_read_ahead_rnd   |   global   |                                                                     read-ahead                                                                      |  operations/s  |
| innodb_buffer_pool_ops              |   global   |                                                                disk_reads, wait_free                                                                |  operations/s  |
| innodb_os_log                       |   global   |                                                                   fsyncs, writes                                                                    |   operations   |
| innodb_os_log_fsync_writes          |   global   |                                                                       fsyncs                                                                        |  operations/s  |
| innodb_os_log_io                    |   global   |                                                                        write                                                                        |     KiB/s      |
| innodb_deadlocks                    |   global   |                                                                      deadlocks                                                                      |  operations/s  |
| files                               |   global   |                                                                        files                                                                        |     files      |
| files_rate                          |   global   |                                                                        files                                                                        |    files/s     |
| connection_errors                   |   global   |                                                  accept, internal, max, peer_addr, select, tcpwrap                                                  |    errors/s    |
| opened_tables                       |   global   |                                                                       tables                                                                        |    tables/s    |
| open_tables                         |   global   |                                                                    cache, tables                                                                    |     tables     |
| process_list_fetch_query_duration   |   global   |                                                                      duration                                                                       |  milliseconds  |
| process_list_queries_count          |   global   |                                                                    system, user                                                                     |    queries     |
| process_list_longest_query_duration |   global   |                                                                      duration                                                                       |    seconds     |
| qcache_ops                          |   global   |                                                      hits, lowmem_prunes, inserts, not_cached                                                       |   queries/s    |
| qcache                              |   global   |                                                                       queries                                                                       |    queries     |
| qcache_freemem                      |   global   |                                                                        free                                                                         |      MiB       |
| qcache_memblocks                    |   global   |                                                                     free, total                                                                     |     blocks     |
| galera_writesets                    |   global   |                                                                       rx, tx                                                                        |  writesets/s   |
| galera_bytes                        |   global   |                                                                       rx, tx                                                                        |     KiB/s      |
| galera_queue                        |   global   |                                                                       rx, tx                                                                        |   writesets    |
| galera_conflicts                    |   global   |                                                                bf_aborts, cert_fails                                                                |  transactions  |
| galera_flow_control                 |   global   |                                                                       paused                                                                        |       ms       |
| galera_cluster_status               |   global   |                                                         primary, non_primary, disconnected                                                          |     status     |
| galera_cluster_state                |   global   |                                                  undefined, joining, donor, joined, synced, error                                                   |     state      |
| galera_cluster_size                 |   global   |                                                                        nodes                                                                        |     nodes      |
| galera_cluster_weight               |   global   |                                                                       weight                                                                        |     weight     |
| galera_connected                    |   global   |                                                                      connected                                                                      |    boolean     |
| galera_ready                        |   global   |                                                                        ready                                                                        |    boolean     |
| galera_open_transactions            |   global   |                                                                        open                                                                         |  transactions  |
| galera_thread_count                 |   global   |                                                                       threads                                                                       |    threads     |
| key_blocks                          |   global   |                                                              unused, used, not_flushed                                                              |     blocks     |
| key_requests                        |   global   |                                                                    reads, writes                                                                    |   requests/s   |
| key_disk_ops                        |   global   |                                                                    reads, writes                                                                    |  operations/s  |
| binlog_cache                        |   global   |                                                                      disk, all                                                                      | transactions/s |
| binlog_stmt_cache                   |   global   |                                                                      disk, all                                                                      |  statements/s  |
| slave_behind                        | connection |                                                                       seconds                                                                       |    seconds     |
| slave_status                        | connection |                                                               sql_running, io_running                                                               |    boolean     |
| userstats_cpu                       |    user    |                                                                        used                                                                         |   percentage   |
| userstats_rows                      |    user    |                                                       read, sent, updated, inserted, deleted                                                        |  operations/s  |
| userstats_commands                  |    user    |                                                                select, update, other                                                                |   commands/s   |
| userstats_denied_commands           |    user    |                                                                       denied                                                                        |   commands/s   |
| userstats_created_transactions      |    user    |                                                                  commit, rollback                                                                   | transactions/s |
| userstats_binlog_written            |    user    |                                                                       written                                                                       |      B/s       |
| userstats_empty_queries             |    user    |                                                                        empty                                                                        |   queries/s    |
| userstats_connections               |    user    |                                                                       created                                                                       | connections/s  |
| userstats_lost_connections          |    user    |                                                                        lost                                                                         | connections/s  |
| userstats_denied_connections        |    user    |                                                                       denied                                                                        | connections/s  |

## Configuration

Edit the `go.d/mysql.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/mysql.conf
```

[DSN syntax in details](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

```yaml
jobs:
  - name: local
    dsn: '[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]'
    # username:password@protocol(address)/dbname?param=value
    # user:password@/dbname
    # Examples:
    # - name: local
    #   dsn: user:pass@unix(/usr/local/var/mysql/mysql.sock)/
    # - name: remote
    #   dsn: user:pass5@localhost/mydb?charset=utf8
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/mysql.conf).

## Troubleshooting

To troubleshoot issues with the `mysql` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m mysql
  ```

