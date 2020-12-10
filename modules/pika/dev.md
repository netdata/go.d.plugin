#### Setup Pika

Run the command ([image page](https://hub.docker.com/r/pikadb/pika)):

```cmd
docker run \
  --name pika \
  -d \
  -p 9221:9221/tcp \
  pikadb/pika:v3.4.0
```

#### Gather metrics

```cmd
echo "INFO ALL" | nc 127.0.0.1 9221
```

Links:

- [github wiki](https://github.com/Qihoo360/pika/wiki) (Chinese).
- [`pika_admin.cc`](https://github.com/Qihoo360/pika/blob/master/src/pika_admin.cc).
