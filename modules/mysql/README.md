<!--
title: "MySQL monitoring with Netdata"
description: "Monitor connections, slow queries, InnoDB memory and disk utilization, locks, and more with zero configuration and per-second metric granularity."
custom_edit_url: https://github.com/netdata/go.d.plugin/edit/master/modules/mysql/README.md
sidebar_label: "MySQL"
-->

# MySQL monitoring with Netdata

[`MySQL`](https://www.mysql.com/) is an open-source relational database management system.

This module monitors one or more `MySQL` servers, depending on your configuration.

## Requirements

Executed queries:

- `SELECT VERSION();`
- `SHOW GLOBAL STATUS;`
- `SHOW GLOBAL VARIABLES;`
- `SHOW SLAVE STATUS;` or `SHOW ALL SLAVES STATUS;` (MariaDBv10.2+)
- `SHOW USER_STATISTICS;` (MariaDBv10.1.1+)

[User Statistics](https://mariadb.com/kb/en/user-statistics/) query is [`MariaDB`](https://mariadb.com/) specific.

`MySQL` user should have the following [permissions](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html):

- [`USAGE`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_usage)
- [`REPLICATION CLIENT`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_replication-client)
- [`PROCESS`](https://dev.mysql.com/doc/refman/8.0/en/privileges-provided.html#priv_process)

To create the `netdata` user with these permissions, execute the following in the `MySQL` shell:

```mysql
CREATE USER 'netdata'@'localhost';
GRANT USAGE, REPLICATION CLIENT, PROCESS ON *.* TO 'netdata'@'localhost';
FLUSH PRIVILEGES;
```

The `netdata` user will have the ability to connect to the `MySQL` server on localhost without a password. It will only
be able to gather statistics without being able to alter or affect operations in any way.

## Charts

It produces the following charts:

- Bandwidth in `kilobits/s`
- Queries in `queries/s`
- Queries By Type in `queries/s`
- Handlers in `handlers/s`
- Table Locks in `locks/s`
- Table Select Join Issues in `joins/s`
- Table Sort Issues in `joins/s`
- Tmp Operations in `events/s`
- Connections in `connections/s`
- Active Connections in `connections`
- Binlog Cache in `transactions/s`
- Threads in `threads`
- Threads Creation Rate in `threads/s`
- Threads Cache Misses in `misses`
- InnoDB I/O Bandwidth in `KiB/s`
- InnoDB I/O Operations in `operations/s`
- InnoDB Pending I/O Operations in `operations`
- InnoDB Log Operations in `operations/s`
- InnoDB OS Log Pending Operations in `operations`
- InnoDB OS Log Operations in `operations/s`
- InnoDB OS Log Bandwidth in `KiB/s`
- InnoDB Current Row Locks in `operations`
- InnoDB Row Operations in `operations/s`
- InnoDB Buffer Pool Pages in `pages`
- InnoDB Buffer Pool Flush Pages Requests in `requests/s`
- InnoDB Buffer Pool Bytes in `MiB`
- InnoDB Buffer Pool Operations in `operations/s`
- MyISAM Key Cache Blocks in `blocks`
- MyISAM Key Cache Requests in `requests/s`
- MyISAM Key Cache Disk Operations in `operations/s`
- Open Files in `files`
- Opened Files Rate in `files/s`
- Binlog Statement Cache in `statements/s`
- Connection Errors in `errors/s`
- Opened Tables in `tables/s`
- Open Tables in `tables`

If [Query Cache](https://dev.mysql.com/doc/refman/5.7/en/query-cache.html) metrics are available (`MariaDB`
and [old versions of `MySQL`](https://mysqlserverteam.com/mysql-8-0-retiring-support-for-the-query-cache/)):

- QCache Operations in `queries/s`
- QCache Queries in Cache in `queries`
- QCache Free Memory in `MiB`
- QCache Memory Blocks in `blocks`

If [WSRep](https://galeracluster.com/library/documentation/galera-status-variables.html) metrics are available:

- Replicated Writesets in `writesets/s`
- Replicated Bytes in `KiB/s`
- Galera Queue in `writesets`
- Replication Conflicts in `transactions`
- Flow Control in `ms`
- Cluster Component Status in `status`
- Cluster Component State in `state`
- Number of Nodes in the Cluster in `num`
- The Total Weight of the Current Members in the Cluster in `weight`
- Cluster Connection Status in `boolean`
- Accept Queries Readiness Status in `boolean`
- Open Transactions in `num`
- Total Number of WSRep (applier/rollbacker) Threads in `num`

If [Slave Status](https://dev.mysql.com/doc/refman/8.0/en/show-slave-status.html) metrics are available:

- Slave Behind Seconds in `seconds`
- I/O / SQL Thread Running State in `boolean`

If [User Statistics](https://mariadb.com/kb/en/user-statistics/) metrics are available:

- User CPU Time in `percentage`
- Rows Operations in `operations/s`
- Commands in `commands/s`

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

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugin's directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` to debug the collector:

```bash
./go.d.plugin -d -m mysql
```
