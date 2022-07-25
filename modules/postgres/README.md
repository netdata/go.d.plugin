<!--
title: "PostgreSQL monitoring with Netdata"
description: "Monitor connections, slow queries, InnoDB memory and disk utilization, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/postgres/README.md
sidebar_label: "PostgresSQL"
-->

# PostgreSQL monitoring with Netdata

[PostgreSQL](https://www.postgresql.org/), also known as Postgres, is a free and open-source relational database
management system emphasizing extensibility and SQL compliance.

This module monitors one or more Postgres servers, depending on your configuration.

## Metrics

All metrics have "postgres." prefix.

| Metric              |                   Dimensions                    |     Units      |
|---------------------|:-----------------------------------------------:|:--------------:|
| db_transactions     |               committed, rollback               | transactions/s |
| db_connections      |                   connections                   |  connections   |
| db_buffer_cache     |                    hit, miss                    |    blocks/s    |
| db_read_operations  |                returned, fetched                |     rows/s     |
| db_write_operations |           inserted, deleted, updated            |     rows/s     |
| db_conflicts        |                    conflicts                    |   queries/s    |
| db_conflicts_stat   | tablespace, lock, snapshot, bufferpin, deadlock |   queries/s    |
| db_deadlocks        |                    deadlocks                    |  deadlocks/s   |
| db_temp_files       |                     written                     |    files/s     |
| db_temp_files_data  |                     written                     |      B/s       |
| db_size             |                      size                       |       B        |

## Configuration

Edit the `go.d/postgres.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/postgres.conf
```

[DSN syntax in details](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

```yaml
jobs:
  - name: local
    # the format is: postgres://[username[:password]]@host:port[/dbname]?sslmode=[disable|verify-ca|verify-full]
    dsn: 'postgres://postgres:postgres@127.0.0.1:5432/postgres'
  - name: remote
    # the format is: postgres://[username[:password]]@host:port[/dbname]?sslmode=[disable|verify-ca|verify-full]
    dsn: 'postgres://postgres:postgres@203.0.113.10:5432/postgres'
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/postgres.conf).

## Troubleshooting

To troubleshoot issues with the `postgres` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m postgres
```
