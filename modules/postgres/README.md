<!--
title: "PostgreSQL monitoring with Netdata"
description: "Monitor connections, replication, databases, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/postgres/README.md"
sidebar_label: "PostgresSQL"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# PostgreSQL monitoring with Netdata

[PostgreSQL](https://www.postgresql.org/), also known as Postgres, is a free and open-source relational database
management system emphasizing extensibility and SQL compliance.

This module monitors one or more Postgres servers, depending on your configuration.

To find out more about these metrics and why they are important to monitor, read our blog post on [PostgreSQL monitoring with Netdata](https://www.netdata.cloud/blog/postgresql-monitoring)

## Requirements

- PostgreSQL v9.4+
- User with granted `pg_monitor`
  or `pg_read_all_stat` [built-in role](https://www.postgresql.org/docs/current/predefined-roles.html).

Without additional configuration Netdata will attempt to use the default postgres user - but a separate netdata user can be created for this purpose. If you have PostgreSQL 10+ and want to use Netdata to monitor statistics normally reserved for superusers grant the netdata user pg_monitor permissions. Some of the advanced metrics also require additional permissions as mentioned in [metrics](#metrics).

To create the `netdata` user with these permissions, execute the following in the psql session, as a user with CREATEROLE priviliges:

```postgresql
CREATE USER netdata;
GRANT pg_monitor TO netdata;
```

After creating the new user, restart the Netdata agent with `sudo systemctl restart netdata`, or the [appropriate method](https://learn.netdata.cloud/docs/configure/start-stop-restart) for your system

## Metrics

- all metrics have "postgres." prefix.
- db_size need CONNECT privilege to the database.
- table_* and index_* metrics need [additional configuration](#database-detailed-metrics).
- table_bloat* and index_bloat* metrics need read (SELECT) permission to the table.
- wal_files_count, wal_archiving_files_count and replication_slot_files_count
  need [superuser](https://www.postgresql.org/docs/current/role-attributes.html) status.

Labels per scope:

- global: no labels.
- repl application: application.
- repl slot: slot.
- database: database.
- table: database, schema, table.
- index: database, schema, table, index.

| Metric                                  |      Scope       |                                                                 Dimensions                                                                 |     Units      |
|-----------------------------------------|:----------------:|:------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| connections_utilization                 |      global      |                                                                    used                                                                    |   percentage   |
| connections_usage                       |      global      |                                                              available, used                                                               |  connections   |
| connections_state_count                 |      global      |                                  active, idle, idle_in_transaction, idle_in_transaction_aborted, disabled                                  |  connections   |
| transactions_duration                   |      global      |                                                       <i>a dimension per bucket</i>                                                        | transactions/s |
| queries_duration                        |      global      |                                                       <i>a dimension per bucket</i>                                                        |   queries/s    |
| locks_utilization                       |      global      |                                                                    used                                                                    |   percentage   |
| checkpoints_rate                        |      global      |                                                            scheduled, requested                                                            | checkpoints/s  |
| checkpoints_time                        |      global      |                                                                write, sync                                                                 |  milliseconds  |
| bgwriter_halts_rate                     |      global      |                                                                 maxwritten                                                                 |    events/s    |
| buffers_io_rate                         |      global      |                                                       checkpoint, backend, bgwriter                                                        |      B/s       |
| buffers_backend_fsync_rate              |      global      |                                                                   fsync                                                                    |    calls/s     |
| buffers_allocated_rate                  |      global      |                                                                 allocated                                                                  |      B/s       |
| wal_io_rate                             |      global      |                                                                   write                                                                    |      B/s       |
| wal_files_count                         |      global      |                                                             written, recycled                                                              |     files      |
| wal_archiving_files_count               |      global      |                                                                ready, done                                                                 |    files/s     |
| autovacuum_workers_count                |      global      |                                       analyze, vacuum_analyze, vacuum, vacuum_freeze, brin_summarize                                       |    workers     |
| txid_exhaustion_towards_autovacuum_perc |      global      |                                                            emergency_autovacuum                                                            |   percentage   |
| txid_exhaustion_perc                    |      global      |                                                              txid_exhaustion                                                               |   percentage   |
| txid_exhaustion_oldest_txid_num         |      global      |                                                                    xid                                                                     |      xid       |
| catalog_relations_count                 |      global      | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |   relations    |
| catalog_relations_size                  |      global      | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |       B        |
| uptime                                  |      global      |                                                                   uptime                                                                   |    seconds     |
| databases_count                         |      global      |                                                                 databases                                                                  |   databases    |
| replication_app_wal_lag_size            | repl application |                                                 sent_lag, write_lag, flush_lag, replay_lag                                                 |       B        |
| replication_app_wal_lag_time            | repl application |                                                      write_lag, flush_lag, replay_lag                                                      |    seconds     |
| replication_slot_files_count            |    repl slot     |                                                        wal_keep, pg_replslot_files                                                         |     files      |
| db_transactions_ratio                   |     database     |                                                            committed, rollback                                                             |   percentage   |
| db_transactions_rate                    |     database     |                                                            committed, rollback                                                             | transactions/s |
| db_connections_utilization              |     database     |                                                                    used                                                                    |   percentage   |
| db_connections_count                    |     database     |                                                                connections                                                                 |  connections   |
| db_cache_io_ratio                       |     database     |                                                                    miss                                                                    |   percentage   |
| db_io_rate                              |     database     |                                                                memory, disk                                                                |      B/s       |
| db_ops_fetched_rows_ratio               |     database     |                                                                  fetched                                                                   |   percentage   |
| db_ops_read_rows_rate                   |     database     |                                                             returned, fetched                                                              |     rows/s     |
| db_ops_write_rows_rate                  |     database     |                                                         inserted, deleted, updated                                                         |     rows/s     |
| db_conflicts_rate                       |     database     |                                                                 conflicts                                                                  |   queries/s    |
| db_conflicts_reason_rate                |     database     |                                              tablespace, lock, snapshot, bufferpin, deadlock                                               |   queries/s    |
| db_deadlocks_rate                       |     database     |                                                                 deadlocks                                                                  |  deadlocks/s   |
| db_locks_held_count                     |     database     |               access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive                |     locks      |
| db_locks_awaited_count                  |     database     |               access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive                |     locks      |
| db_temp_files_created_rate              |     database     |                                                                  created                                                                   |    files/s     |
| db_temp_files_io_rate                   |     database     |                                                                  written                                                                   |      B/s       |
| db_size                                 |     database     |                                                                    size                                                                    |       B        |
| table_rows_dead_ratio                   |      table       |                                                                    dead                                                                    |   percentage   |
| table_rows_count                        |      table       |                                                                 live, dead                                                                 |      rows      |
| table_ops_rows_rate                     |      table       |                                                         inserted, deleted, updated                                                         |     rows/s     |
| table_ops_rows_hot_ratio                |      table       |                                                                    hot                                                                     |   percentage   |
| table_ops_rows_hot_rate                 |      table       |                                                                    hot                                                                     |     rows/s     |
| table_cache_io_ratio                    |      table       |                                                                    miss                                                                    |   percentage   |
| table_io_rate                           |      table       |                                                                memory, disk                                                                |      B/s       |
| table_index_cache_io_ratio              |      table       |                                                                    miss                                                                    |   percentage   |
| table_index_io_rate                     |      table       |                                                                memory, disk                                                                |      B/s       |
| table_toast_cache_io_ratio              |      table       |                                                                    miss                                                                    |   percentage   |
| table_toast_io_rate                     |      table       |                                                                memory, disk                                                                |      B/s       |
| table_toast_index_cache_io_ratio        |      table       |                                                                    miss                                                                    |   percentage   |
| table_toast_index_io_rate               |      table       |                                                                memory, disk                                                                |      B/s       |
| table_scans_rate                        |      table       |                                                             index, sequential                                                              |    scans/s     |
| table_scans_rows_rate                   |      table       |                                                             index, sequential                                                              |     rows/s     |
| table_autovacuum_since_time             |      table       |                                                                    time                                                                    |    seconds     |
| table_vacuum_since_time                 |      table       |                                                                    time                                                                    |    seconds     |
| table_autoanalyze_since_time            |      table       |                                                                    time                                                                    |    seconds     |
| table_analyze_since_time                |      table       |                                                                    time                                                                    |    seconds     |
| table_null_columns                      |      table       |                                                                    null                                                                    |    columns     |
| table_size                              |      table       |                                                                    size                                                                    |       B        |
| table_bloat_size_perc                   |      table       |                                                                   bloat                                                                    |   percentage   |
| table_bloat_size                        |      table       |                                                                   bloat                                                                    |       B        |
| index_size                              |      index       |                                                                    size                                                                    |       B        |
| index_bloat_size_perc                   |      index       |                                                                   bloat                                                                    |   percentage   |
| index_bloat_size                        |      index       |                                                                   bloat                                                                    |       B        |
| index_usage_status                      |      index       |                                                                used, unused                                                                |     status     |

## Configuration

Edit the `go.d/postgres.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/postgres.conf
```

DSN (Data Source Name) may either be in URL format or key=word format.
See [Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for details.

```yaml
jobs:
  - name: local
    dsn: 'postgres://postgres:postgres@127.0.0.1:5432/postgres'

  - name: local
    dsn: 'host=/var/run/postgresql dbname=postgres user=postgres'

  - name: remote
    dsn: 'postgres://postgres:postgres@203.0.113.10:5432/postgres'
```

### Database detailed metrics

Detailed metrics include table_* and index_*.

By default, this module only collects detailed metrics for the database it is connected to. Collection from all
databases on a database server is disabled because each database requires an additional connection.

Use the `collect_databases_matching` configuration option to select the databases from which you want to collect
detailed metrics. The value
supports [Netdata simple patterns](https://learn.netdata.cloud/docs/agent/libnetdata/simple_pattern).

```yaml
jobs:
  - name: local
    dsn: 'postgres://postgres:postgres@127.0.0.1:5432/postgres'
    collect_databases_matching: 'mydb1 mydb2 !mydb3 mydb4'
```

---

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/postgres.conf).

## Troubleshooting

To troubleshoot issues with the `postgres` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m postgres
  ```
