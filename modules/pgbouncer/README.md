# PostgreSQL collector

## Overview

[PgBouncer](https://www.pgbouncer.org/) is an open-source connection pooler
for [PostgreSQL](https://www.postgresql.org/).

This collector monitors one or more PgBouncer servers, depending on your configuration.

Executed queries:

- `SHOW VERSION;`
- `SHOW CONFIG;`
- `SHOW DATABASES;`
- `SHOW STATS;`
- `SHOW POOLS;`

Information about the queries can be found in the [PgBouncer Documentation](http://pgbouncer.org/usage.html).

## Collected metrics

Metrics grouped by *scope*.

The scope defines the instance that the metric belongs to. An instance is uniquely identified by a set of labels.

### global

These metrics refer to the entire monitored application.

This scope has no labels.

Metrics:

| Metric                                   | Dimensions |    Unit    |
|------------------------------------------|:----------:|:----------:|
| pgbouncer.client_connections_utilization |    used    | percentage |

### database

These metrics refer to the database.

Labels:

| Label             | Description            |
|-------------------|------------------------|
| database          | database name          |
| postgres_database | Postgres database name |

Metrics:

| Metric                                      |            Dimensions             |      Unit      |
|---------------------------------------------|:---------------------------------:|:--------------:|
| pgbouncer.db_client_connections             |    active, waiting, cancel_req    |  connections   |
| pgbouncer.db_server_connections             | active, idle, used, tested, login |  connections   |
| pgbouncer.db_server_connections_utilization |               used                |   percentage   |
| pgbouncer.db_clients_wait_time              |               time                |    seconds     |
| pgbouncer.db_client_max_wait_time           |               time                |    seconds     |
| pgbouncer.db_transactions                   |           transactions            | transactions/s |
| pgbouncer.db_transactions_time              |               time                |    seconds     |
| pgbouncer.db_transaction_avg_time           |               time                |    seconds     |
| pgbouncer.db_queries                        |              queries              |   queries/s    |
| pgbouncer.db_queries_time                   |               time                |    seconds     |
| pgbouncer.db_query_avg_time                 |               time                |    seconds     |
| pgbouncer.db_network_io                     |          received, sent           |      B/s       |

## Setup

### Prerequisites

#### Create netdata user

Create a user with `stats_users` permissions to query your PgBouncer instance.

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

### Configuration

#### File

The configuration file name is `go.d/postgresql.conf`.

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
sudo ./edit-config go.d/postgresql.conf
```

#### Options

The following options can be defined globally: update_every, autodetection_retry.

<details>
<summary>Config options</summary>

|        Name         | Description                                                                                                                             |                        Default                        | Required |
|:-------------------:|-----------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------------------------:|:--------:|
|    update_every     | Data collection frequency.                                                                                                              |                           5                           |          |
| autodetection_retry | Re-check interval in seconds. Zero means not to schedule re-check.                                                                      |                           0                           |          |
|         dsn         | PgBouncer server DSN (Data Source Name). See [DSN syntax](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING). | postgres://postgres:postgres@127.0.0.1:6432/pgbouncer |   yes    |
|       timeout       | Query timeout in seconds.                                                                                                               |                           1                           |          |

</details>

#### Examples

##### TCP socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: 'postgres://postgres:postgres@127.0.0.1:6432/pgbouncer'
```

</details>

##### Unix socket

An example configuration.
<details>
<summary>Config</summary>

```yaml
jobs:
  - name: local
    dsn: 'host=/tmp dbname=pgbouncer user=postgres port=6432'
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
    dsn: 'postgres://postgres:postgres@127.0.0.1:6432/pgbouncer'

  - name: remote
    dsn: 'postgres://postgres:postgres@203.0.113.10:6432/pgbouncer'
```

</details>

## Troubleshooting

### Debug mode

To troubleshoot issues with the `postgresql` collector, run the `go.d.plugin` with the debug option enabled.
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
  ./go.d.plugin -d -m postgresql
  ```
