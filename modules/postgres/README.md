# PostgreSQL collector

## Overview

[PostgreSQL](https://www.postgresql.org/), also known as Postgres, is a free and open-source relational database
management system emphasizing extensibility and SQL compliance.

This collector monitors one or more Postgres servers, depending on your configuration.

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                           |                                                                 Dimensions                                                                 |      Unit      |
|--------------------------------------------------|:------------------------------------------------------------------------------------------------------------------------------------------:|:--------------:|
| postgres.connections_utilization                 |                                                                    used                                                                    |   percentage   |
| postgres.connections_usage                       |                                                              available, used                                                               |  connections   |
| postgres.connections_state_count                 |                                  active, idle, idle_in_transaction, idle_in_transaction_aborted, disabled                                  |  connections   |
| postgres.transactions_duration                   |                                                           a dimension per bucket                                                           | transactions/s |
| postgres.queries_duration                        |                                                           a dimension per bucket                                                           |   queries/s    |
| postgres.locks_utilization                       |                                                                    used                                                                    |   percentage   |
| postgres.checkpoints_rate                        |                                                            scheduled, requested                                                            | checkpoints/s  |
| postgres.checkpoints_time                        |                                                                write, sync                                                                 |  milliseconds  |
| postgres.bgwriter_halts_rate                     |                                                                 maxwritten                                                                 |    events/s    |
| postgres.buffers_io_rate                         |                                                       checkpoint, backend, bgwriter                                                        |      B/s       |
| postgres.buffers_backend_fsync_rate              |                                                                   fsync                                                                    |    calls/s     |
| postgres.buffers_allocated_rate                  |                                                                 allocated                                                                  |      B/s       |
| postgres.wal_io_rate                             |                                                                   write                                                                    |      B/s       |
| postgres.wal_files_count                         |                                                             written, recycled                                                              |     files      |
| postgres.wal_archiving_files_count               |                                                                ready, done                                                                 |    files/s     |
| postgres.autovacuum_workers_count                |                                       analyze, vacuum_analyze, vacuum, vacuum_freeze, brin_summarize                                       |    workers     |
| postgres.txid_exhaustion_towards_autovacuum_perc |                                                            emergency_autovacuum                                                            |   percentage   |
| postgres.txid_exhaustion_perc                    |                                                              txid_exhaustion                                                               |   percentage   |
| postgres.txid_exhaustion_oldest_txid_num         |                                                                    xid                                                                     |      xid       |
| postgres.catalog_relations_count                 | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |   relations    |
| postgres.catalog_relations_size                  | ordinary_table, index, sequence, toast_table, view, materialized_view, composite_type, foreign_table, partitioned_table, partitioned_index |       B        |
| postgres.uptime                                  |                                                                   uptime                                                                   |    seconds     |
| postgres.databases_count                         |                                                                 databases                                                                  |   databases    |

### repl application

These metrics refer to the replication application.

Labels:

| Label       | Description      |
|-------------|------------------|
| application | application name |

Metrics:

| Metric                                |                 Dimensions                 |  Unit   |
|---------------------------------------|:------------------------------------------:|:-------:|
| postgres.replication_app_wal_lag_size | sent_lag, write_lag, flush_lag, replay_lag |    B    |
| postgres.replication_app_wal_lag_time |      write_lag, flush_lag, replay_lag      | seconds |

### repl slot

These metrics refer to the replication slot.

Labels:

| Label | Description           |
|-------|-----------------------|
| slot  | replication slot name |

Metrics:

| Metric                                |         Dimensions          | Unit  |
|---------------------------------------|:---------------------------:|:-----:|
| postgres.replication_slot_files_count | wal_keep, pg_replslot_files | files |

### database

These metrics refer to the database.

Labels:

| Label    | Description   |
|----------|---------------|
| database | database name |

Metrics:

| Metric                              |                                                  Dimensions                                                   |      Unit      |
|-------------------------------------|:-------------------------------------------------------------------------------------------------------------:|:--------------:|
| postgres.db_transactions_ratio      |                                              committed, rollback                                              |   percentage   |
| postgres.db_transactions_rate       |                                              committed, rollback                                              | transactions/s |
| postgres.db_connections_utilization |                                                     used                                                      |   percentage   |
| postgres.db_connections_count       |                                                  connections                                                  |  connections   |
| postgres.db_cache_io_ratio          |                                                     miss                                                      |   percentage   |
| postgres.db_io_rate                 |                                                 memory, disk                                                  |      B/s       |
| postgres.db_ops_fetched_rows_ratio  |                                                    fetched                                                    |   percentage   |
| postgres.db_ops_read_rows_rate      |                                               returned, fetched                                               |     rows/s     |
| postgres.db_ops_write_rows_rate     |                                          inserted, deleted, updated                                           |     rows/s     |
| postgres.db_conflicts_rate          |                                                   conflicts                                                   |   queries/s    |
| postgres.db_conflicts_reason_rate   |                                tablespace, lock, snapshot, bufferpin, deadlock                                |   queries/s    |
| postgres.db_deadlocks_rate          |                                                   deadlocks                                                   |  deadlocks/s   |
| postgres.db_locks_held_count        | access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive |     locks      |
| postgres.db_locks_awaited_count     | access_share, row_share, row_exclusive, share_update, share, share_row_exclusive, exclusive, access_exclusive |     locks      |
| postgres.db_temp_files_created_rate |                                                    created                                                    |    files/s     |
| postgres.db_temp_files_io_rate      |                                                    written                                                    |      B/s       |
| postgres.db_size                    |                                                     size                                                      |       B        |

### table

These metrics refer to the database table.

Labels:

| Label        | Description       |
|--------------|-------------------|
| database     | database name     |
| schema       | schema name       |
| table        | table name        |
| parent_table | parent table name |

Metrics:

| Metric                                    |         Dimensions         |    Unit    |
|-------------------------------------------|:--------------------------:|:----------:|
| postgres.table_rows_dead_ratio            |            dead            | percentage |
| postgres.table_rows_count                 |         live, dead         |    rows    |
| postgres.table_ops_rows_rate              | inserted, deleted, updated |   rows/s   |
| postgres.table_ops_rows_hot_ratio         |            hot             | percentage |
| postgres.table_ops_rows_hot_rate          |            hot             |   rows/s   |
| postgres.table_cache_io_ratio             |            miss            | percentage |
| postgres.table_io_rate                    |        memory, disk        |    B/s     |
| postgres.table_index_cache_io_ratio       |            miss            | percentage |
| postgres.table_index_io_rate              |        memory, disk        |    B/s     |
| postgres.table_toast_cache_io_ratio       |            miss            | percentage |
| postgres.table_toast_io_rate              |        memory, disk        |    B/s     |
| postgres.table_toast_index_cache_io_ratio |            miss            | percentage |
| postgres.table_toast_index_io_rate        |        memory, disk        |    B/s     |
| postgres.table_scans_rate                 |     index, sequential      |  scans/s   |
| postgres.table_scans_rows_rate            |     index, sequential      |   rows/s   |
| postgres.table_autovacuum_since_time      |            time            |  seconds   |
| postgres.table_vacuum_since_time          |            time            |  seconds   |
| postgres.table_autoanalyze_since_time     |            time            |  seconds   |
| postgres.table_analyze_since_time         |            time            |  seconds   |
| postgres.table_null_columns               |            null            |  columns   |
| postgres.table_size                       |            size            |     B      |
| postgres.table_bloat_size_perc            |           bloat            | percentage |
| postgres.table_bloat_size                 |           bloat            |     B      |

### index

These metrics refer to the table index.

Labels:

| Label        | Description       |
|--------------|-------------------|
| database     | database name     |
| schema       | schema name       |
| table        | table name        |
| parent_table | parent table name |
| index        | index name        |

Metrics:

| Metric                         |  Dimensions  |    Unit    |
|--------------------------------|:------------:|:----------:|
| postgres.index_size            |     size     |     B      |
| postgres.index_bloat_size_perc |    bloat     | percentage |
| postgres.index_bloat_size      |    bloat     |     B      |
| postgres.index_usage_status    | used, unused |   status   |

## Setup

### Prerequisites

#### Create netdata user

Create a user with granted `pg_monitor`
or `pg_read_all_stat` [built-in role](https://www.postgresql.org/docs/current/predefined-roles.html).

To create the `netdata` user with these permissions, execute the following in the psql session, as a user with
CREATEROLE privileges:

```postgresql
CREATE USER netdata;
GRANT pg_monitor TO netdata;
```

After creating the new user, restart the Netdata agent with `sudo systemctl restart netdata`, or
the [appropriate method](https://github.com/netdata/netdata/blob/master/docs/configure/start-stop-restart.md) for your
system.

### Configuration

#### File

The configuration file name is `go.d/postgres.conf`.

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
sudo ./edit-config go.d/postgres.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|            Name            | Description                                                                                                                                                                                   |                       Default                        | Required |
|:--------------------------:|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:----------------------------------------------------:|:--------:|
|        update_every        | Data collection frequency.                                                                                                                                                                    |                          5                           |          |
|    autodetection_retry     | Re-check interval in seconds. Zero means not to schedule re-check.                                                                                                                            |                          0                           |          |
|            dsn             | Postgres server DSN (Data Source Name). See [DSN syntax](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING).                                                        | postgres://postgres:postgres@127.0.0.1:5432/postgres |   yes    |
|          timeout           | Query timeout in seconds.                                                                                                                                                                     |                          2                           |          |
| collect_databases_matching | Databases selector. Determines which database metrics will be collected. Syntax is [simple patterns](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#simple-patterns-matcher). |                                                      |          |
|       max_db_tables        | Maximum number of tables in the database. Table metrics will not be collected for databases that have more tables than max_db_tables. 0 means no limit.                                       |                          50                          |          |
|       max_db_indexes       | Maximum number of indexes in the database. Index metrics will not be collected for databases that have more indexes than max_db_indexes. 0 means no limit.                                    |                         250                          |          |

</details>

#### Examples

##### TCP socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: 'postgresql://netdata@127.0.0.1:5432/postgres'
```

</details>

##### Unix socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: 'host=/var/run/postgresql dbname=postgres user=netdata'
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
    dsn: 'postgresql://netdata@127.0.0.1:5432/postgres'

  - name: remote
    dsn: 'postgresql://netdata@203.0.113.0:5432/postgres'
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `postgres` collector, run the `go.d.plugin` with the debug option enabled.
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
  ./go.d.plugin -d -m postgres
  ```
