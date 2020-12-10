#### Setup Redis

Run the command ([image page](https://hub.docker.com/_/redis)):

```cmd
docker run \
  --name redis \
  -d \
  -p 6379:6379/tcp \
  redis:6.0.9
```

#### Gather metrics

```cmd
echo "INFO ALL" | nc 127.0.0.1 6379
```

Links:

- [`INFO` command docs](https://redis.io/commands/info).
