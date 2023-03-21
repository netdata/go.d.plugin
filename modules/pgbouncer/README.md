<!--
title: "PgBouncer monitoring with Netdata"
description: "Monitor client and server connections and databases statistics."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/pgbouncer/README.md"
sidebar_label: "PgBouncer"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# PgBouncer collector

[PgBouncer](https://www.pgbouncer.org/) is an open-source connection pooler
for [PostgreSQL](https://www.postgresql.org/).

This module monitors one or more PgBouncer servers, depending on your configuration.

Executed queries:

- `SHOW VERSION;`
- `SHOW CONFIG;`
- `SHOW DATABASES;`
- `SHOW STATS;`
- `SHOW POOLS;`

Information about the queries can be found in the [PgBouncer Documentation](http://pgbouncer.org/usage.html).

## Requirements

- PgBouncer v1.8.0+.
- A user with `stats_users` permissions to query your PgBouncer instance.

To create the `netdata` user:

- Add `netdata` user to the `pgbouncer.ini` file:

  ```text
  stats_users = netdata
  ```

- Add a password for the `netdata` user to the `userlist.txt` file:

  ```text
  "netdata" "<PASSWORD>"
  ```

- To verify the credentials, run the following command

  ```bash
  psql -h localhost -U netdata -p 6432 pgbouncer -c "SHOW VERSION;" >/dev/null 2>&1 && echo OK || echo FAIL
  ```

  When it prompts for a password, enter the password you added to `userlist.txt`.

## Metrics

All metrics have "pgbouncer." prefix.

Labels per scope:

- global: no labels.
- database: database, postgres_database.

| Metric                            |  Scope   |            Dimensions             |     Units      |
|-----------------------------------|:--------:|:---------------------------------:|:--------------:|
| client_connections_utilization    |  global  |               used                |   percentage   |
| db_client_connections             | database |    active, waiting, cancel_req    |  connections   |
| db_server_connections             | database | active, idle, used, tested, login |  connections   |
| db_server_connections_utilization | database |               used                |   percentage   |
| db_clients_wait_time              | database |               time                |    seconds     |
| db_client_max_wait_time           | database |               time                |    seconds     |
| db_transactions                   | database |           transactions            | transactions/s |
| db_transactions_time              | database |               time                |    seconds     |
| db_transaction_avg_time           | database |               time                |    seconds     |
| db_queries                        | database |              queries              |   queries/s    |
| db_queries_time                   | database |               time                |    seconds     |
| db_query_avg_time                 | database |               time                |    seconds     |
| db_network_io                     | database |          received, sent           |      B/s       |

## Configuration

Edit the `go.d/pgbouncer.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/pgbouncer.conf
```

DSN (Data Source Name) may either be in URL format or key=word format.
See [Connection Strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING) for details.

```yaml
jobs:
  - name: local
    dsn: 'postgres://postgres:postgres@127.0.0.1:6432/pgbouncer'

  - name: local
    dsn: 'host=/tmp dbname=pgbouncer user=postgres port=6432'

  - name: remote
    dsn: 'postgres://postgres:postgres@203.0.113.10:6432/pgbouncer'
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pgbouncer.conf).

## Troubleshooting

To troubleshoot issues with the `pgbouncer` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m pgbouncer
  ```
