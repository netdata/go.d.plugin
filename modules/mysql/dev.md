# MySQL/MariaDB lab setup

- MariaDB
    - [Standalone](https://github.com/bitnami/containers/tree/main/bitnami/mariadb#tldr)
    - [Replication Cluster (Master-Slave)](https://github.com/bitnami/containers/tree/main/bitnami/mariadb#setting-up-a-replication-cluster)
    - [Replication Cluster (Galera)](https://github.com/bitnami/containers/tree/main/bitnami/mariadb-galera#setting-up-a-multi-master-cluster)
- MySQL
    - [Replication Cluster (Multi-Source)](https://github.com/wagnerjfr/mysql-multi-source-replication-docker)

MySQL Slave instance "Authentication requires secure connection" connection fix:
> ALTER USER 'repl1'@'%' IDENTIFIED WITH mysql_native_password BY 'slavepass';

Create `netdata` user with needed permissions:

```mysql
CREATE USER 'netdata'@'%' IDENTIFIED BY 'password';
GRANT USAGE ON *.* TO 'netdata';
GRANT REPLICATION CLIENT ON *.* TO 'netdata';
GRANT PROCESS on *.* to 'netdata';
FLUSH PRIVILEGES;
```

Enables User Statistics metrics collection in `MariaDB`:

- https://mariadb.com/kb/en/user-statistics/
