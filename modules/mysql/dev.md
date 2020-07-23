`MariaDB` Galera cluster setup:

-   https://github.com/bitnami/bitnami-docker-mariadb-galera
-   https://hub.docker.com/r/bitnami/mariadb-galera

`MariaDB` cluster setup:

-   https://github.com/bitnami/bitnami-docker-mariadb
-   https://hub.docker.com/r/bitnami/mariadb/

`MySQL` multi-source replication setup:

-   https://github.com/wagnerjfr/mysql-multi-source-replication-docker 

Create `netdata` user with needed permissions:

```sql
CREATE USER 'netdata'@'%' IDENTIFIED BY 'password';
GRANT USAGE ON *.* TO 'netdata';
GRANT REPLICATION CLIENT ON *.* TO 'netdata';
GRANT PROCESS on *.* to 'netdata';
FLUSH PRIVILEGES;
```

Enables User Statistics metrics collection in `MariaDB`:

-   https://mariadb.com/kb/en/user-statistics/
