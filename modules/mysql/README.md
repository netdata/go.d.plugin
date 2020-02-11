# MySQL monitoring with Netdata

[`MySQL`](https://www.mysql.com/) is an open-source relational database management system.

This module will monitor one or more `MySQL` servers, depending on your configuration.

## Requirements

MySQL user specified in configuration should have at least `USAGE, REPLICATION CLIENT` permissions.

To create the user, enter following to MySQL shell:

```sql
CREATE USER 'netdata'@'localhost';
GRANT USAGE, REPLICATION CLIENT ON *.* TO 'netdata'@'localhost';
FLUSH PRIVILEGES;
```

## Charts

It will produce following charts:

-   Bandwidth in `kilobits/s`
-   Queries in `queries/s`
-   Queries By Type in `queries/s`
-   Handlerse in `handlers/s`
-   Table Locks in `locks/s`
-   Table Select Join Issuess in `joins/s`
-   Table Sort Issuess in `joins/s`
-   Tmp Operations in `created/s`
-   Connections in `connections/s`
-   Connections Active in `connections/s`
-   Binlog Cache in `threads`
-   Threads in `transactions/s`
-   Threads Creation Rate in `threads/s`
-   Threads Cache Misses in `misses`
-   InnoDB I/O Bandwidth in `KiB/s`
-   InnoDB I/O Operations in `operations/s`
-   InnoDB Pending I/O Operations in `operations/s`
-   InnoDB Log Operations in `operations/s`
-   InnoDB OS Log Pending Operations in `operations`
-   InnoDB OS Log Operations in `operations/s`
-   InnoDB OS Log Bandwidth in `KiB/s`
-   InnoDB Current Row Locks in `operations`
-   InnoDB Row Operations in `operations/s`
-   InnoDB Buffer Pool Pagess in `pages`
-   InnoDB Buffer Pool Flush Pages Requests in `requests/s`
-   InnoDB Buffer Pool Bytes in `MiB`
-   InnoDB Buffer Pool Operations in `operations/s`
-   QCache Operations in `queries/s`
-   QCache Queries in Cache in `queries`
-   QCache Free Memory in `MiB`
-   QCache Memory Blocks in `blocks`
-   MyISAM Key Cache Blocks in `blocks`
-   MyISAM Key Cache Requests in `requests/s`
-   MyISAM Key Cache Requests in `requests/s`
-   MyISAM Key Cache Disk Operations in `operations/s`
-   Open Files in `files`
-   Opened Files Rate in `files/s`
-   Binlog Statement Cache in `statements/s`
-   Connection Errors in `errors/s`
-   Slave Behind Seconds in `seconds`
-   I/O / SQL Thread Running Statein `bool`
-   Replicated Writesets in `writesets/s`
-   Replicated Bytes in `KiB/s`
-   Galera Queue in `writesets`
-   Replication Conflicts in `transactions`
-   Flow Control in `ms`

## Configuration

Edit the `go.d/mysql.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/mysql.conf
```

[DSN syntax in details](https://github.com/go-sql-driver/mysql#dsn-data-source-name).

```yaml
jobs:
  - name: local
    dsn: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
    # username:password@protocol(address)/dbname?param=value
    # user:password@/dbname
    # Examples:
    # - name: local
    #   dsn: user:pass@unix(/usr/local/var/mysql/mysql.sock)/
    # - name: remote
    #   dsn: user:pass5@localhost/mydb?charset=utf8
```

If no configuration is given, module will attempt to connect to mysql server via unix socket in the following order:

-   `/var/run/mysqld/mysqld.sock` without password and with username `root`;
-   `/usr/local/var/mysql/mysql.sock` without password and with username `root`;
-   `localhost:3306` without password and with username `root`.


For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/mysql.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m mysql
