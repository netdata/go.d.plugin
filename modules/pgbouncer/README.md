<!--
title: "PgBouncer monitoring with Netdata"
description: "Monitor client and server connections, databases statistics."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/postgres/README.md
sidebar_label: "PgBouncer"
-->

# PgBouncer monitoring with Netdata

[PgBouncer](https://www.pgbouncer.org/) is an open-source connection pooler for PostgreSQL.

This module monitors one or more PgBouncer servers, depending on your configuration.

## Metrics

All metrics have "pgbouncer." prefix.

| Metric                       | Scope  | Dimensions |    Units    |
|------------------------------|:------:|:----------:|:-----------:|
| pgbouncer.client_connections | global | free, used | connections |
| pgbouncer.server_connections | global | free, used | connections |

## Configuration

Edit the `go.d/pgbouncer.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

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
    dsn: 'host=/var/run/postgresql dbname=postgres user=postgres port=6432'

  - name: remote
    dsn: 'postgres://postgres:postgres@203.0.113.10:6432/pgbouncer'
```

For all available options see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/pgbouncer.conf).

## Troubleshooting

To troubleshoot issues with the `pgbouncer` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

First, navigate to your plugins' directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m pgbouncer
```
