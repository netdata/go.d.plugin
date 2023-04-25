<!--
title: "PostgreSQL monitoring with Netdata"
description: "Monitor connections, replication, databases, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/postgres/README.md"
sidebar_label: "PostgresSQL"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# PostgreSQL collector

[PostgreSQL](https://www.postgresql.org/), also known as Postgres, is a free and open-source relational database
management system emphasizing extensibility and SQL compliance.

This module monitors one or more Postgres servers, depending on your configuration.

To find out more about these metrics and why they are important to monitor, read our blog post
on [PostgreSQL monitoring with Netdata](https://www.netdata.cloud/blog/postgresql-monitoring)

## Requirements

- PostgreSQL v9.4+
- User with granted `pg_monitor`
  or `pg_read_all_stat` [built-in role](https://www.postgresql.org/docs/current/predefined-roles.html).

Without additional configuration Netdata will attempt to use the default postgres user - but a separate netdata user can
be created for this purpose. If you have PostgreSQL 10+ and want to use Netdata to monitor statistics normally reserved
for superusers grant the netdata user pg_monitor permissions. Some of the advanced metrics also require additional
permissions as mentioned in [metrics](#metrics).

To create the `netdata` user with these permissions, execute the following in the psql session, as a user with
CREATEROLE priviliges:

```postgresql
CREATE USER netdata;
GRANT pg_monitor TO netdata;
```

After creating the new user, restart the Netdata agent with `sudo systemctl restart netdata`, or
the [appropriate method](https://github.com/netdata/netdata/blob/master/docs/configure/start-stop-restart.md) for your
system

## Metrics

- db_size need CONNECT privilege to the database.
- table_* and index_* metrics need [additional configuration](#database-detailed-metrics).
- table_bloat* and index_bloat* metrics need read (SELECT) permission to the table.
- wal_files_count, wal_archiving_files_count and replication_slot_files_count
  need [superuser](https://www.postgresql.org/docs/current/role-attributes.html) status.

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/postgres/metrics.csv) for a list
of metrics.

## Configuration

Edit the `go.d/postgres.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/postgres.conf
```

DSN (Data Source Name) may either be in URL format or key=word format.
See [Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for details.

```yaml
jobs:
  - name: local
    dsn: 'postgresql://postgres:postgres@127.0.0.1:5432/postgres'

  - name: local
    dsn: 'host=/var/run/postgresql dbname=postgres user=postgres'

  - name: remote
    dsn: 'postgresql://postgres:postgres@203.0.113.10:5432/postgres'
```

### Database detailed metrics

Detailed metrics include table_* and index_*.

By default, this module only collects detailed metrics for the database it is connected to. Collection from all
databases on a database server is disabled because each database requires an additional connection.

Use the `collect_databases_matching` configuration option to select the databases from which you want to collect
detailed metrics. The value
supports [Netdata simple patterns](https://github.com/netdata/netdata/blob/master/libnetdata/simple_pattern/README.md).

```yaml
jobs:
  - name: local
    dsn: 'postgresql://postgres:postgres@127.0.0.1:5432/postgres'
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
