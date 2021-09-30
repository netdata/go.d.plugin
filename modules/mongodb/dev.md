#### Setup Mongo

Run the command ([image page](https://hub.docker.com/_/mongo)):

```shell
docker run \
  --name mongodb \
  -d \
  -p 27017:27017 \
  mongo:5.0.0
```

for testing shards you could use
a [pre-made docker-compose.yml](https://github.com/bitnami/bitnami-docker-mongodb-sharded/blob/master/docker-compose.yml)

### run the module

```shell
go build -o go.d.plugin github.com/netdata/go.d.plugin/cmd/godplugin
go.d.plugin -d -m=mongodb
```
