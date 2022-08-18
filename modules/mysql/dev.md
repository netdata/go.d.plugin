# MySQL/MariaDB lab setup

## MariaDB

### Standalone

- [GitHub Link](https://github.com/bitnami/containers/tree/main/bitnami/mariadb#tldr)

```bash
# old version
docker run -d \
--network=host \
--name mariadb-5.5.64 \
--env MYSQL_ROOT_PASSWORD=password  \
mariadb:5.5.64

# recent version
docker run -d \
--network=host \
--name mariadb-10.8.3 \
--env MARIADB_USER=netdata \
--env MARIADB_PASSWORD=password \
--env MARIADB_ROOT_PASSWORD=password  \
mariadb:10.8.3 --port 3808
```

### Replication Cluster (Master-Slave)

- [GitHub Link](https://github.com/bitnami/containers/tree/main/bitnami/mariadb#setting-up-a-replication-cluster)

### Replication Cluster (Galera)

- [GitHub Link](https://github.com/bitnami/containers/tree/main/bitnami/mariadb-galera#setting-up-a-multi-master-cluster)
- [DockerHub Link](https://hub.docker.com/r/bitnami/mariadb-galera)

```bash
docker run -d --name mariadb-galera-0 \
  -e MARIADB_GALERA_CLUSTER_NAME=my_galera \
  -e MARIADB_GALERA_MARIABACKUP_USER=my_mariabackup_user \
  -e MARIADB_GALERA_MARIABACKUP_PASSWORD=my_mariabackup_password \
  -e MARIADB_ROOT_PASSWORD=my_root_password \
  -e MARIADB_GALERA_CLUSTER_BOOTSTRAP=yes \
  -e MARIADB_USER=my_user \
  -e MARIADB_PASSWORD=my_password \
  -e MARIADB_DATABASE=my_database \
  -e MARIADB_REPLICATION_USER=my_replication_user \
  -e MARIADB_REPLICATION_PASSWORD=my_replication_password \
  bitnami/mariadb-galera:10.8.4

  docker run -d --name mariadb-galera-1 --link mariadb-galera-0:mariadb-galera \
  -e MARIADB_GALERA_CLUSTER_NAME=my_galera \
  -e MARIADB_GALERA_CLUSTER_ADDRESS=gcomm://mariadb-galera:4567,0.0.0.0:4567 \
  -e MARIADB_GALERA_MARIABACKUP_USER=my_mariabackup_user \
  -e MARIADB_GALERA_MARIABACKUP_PASSWORD=my_mariabackup_password \
  -e MARIADB_ROOT_PASSWORD=my_root_password \
  -e MARIADB_REPLICATION_USER=my_replication_user \
  -e MARIADB_REPLICATION_PASSWORD=my_replication_password \
  bitnami/mariadb-galera:10.8.4
```

## Create user

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
