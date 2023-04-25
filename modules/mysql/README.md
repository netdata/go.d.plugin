<!--
title: "MySQL monitoring with Netdata"
description: "Monitor connections, slow queries, InnoDB memory and disk utilization, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/mysql/README.md"
sidebar_label: "MySQL"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Databases"
-->

# MySQL collector

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

> userstats_* metrics need [User Statistics](https://mariadb.com/kb/en/user-statistics/#enabling-the-plugin) plugin
> enabled. MariaDB and Percona MySQL only.

See [metrics.csv](https://github.com/netdata/go.d.plugin/blob/master/modules/mysql/metrics.csv) for a list of
metrics.

## Configuration

Edit the `go.d/mysql.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically
at `/etc/netdata`.

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

