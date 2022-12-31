#### Setup Mongo

- For a standalone setup, run the command ([image page](https://hub.docker.com/_/mongo)):

  ```shell
  docker run \
    --name mongodb \
    -d \
    -p 27017:27017 \
    mongo:5.0.0
  ```

- Replica set setup:
  use [MongoDB packaged by Bitnami](https://github.com/bitnami/containers/tree/main/bitnami/mongodb#setting-up-replication) (
  use `--net` instead of `--link`).

- Sharding setup:
  use [MongoDB Sharded packaged by Bitnami](https://github.com/bitnami/containers/tree/main/bitnami/mongodb-sharded)
