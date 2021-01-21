#### Setup NGINX with nginx-vts

- Run the command ([image page](https://hub.docker.com/r/xcgd/nginx-vts/)):

```cmd
docker run \
  --name nginxvts \
  -d \
  -p 80:80/tcp \
  xcgd/nginx-vts:1.16.1-0.1.18
```

- Add `vhost_traffic_status_zone;` to the `/etc/nginx/nginx.conf`

```
...
http {
    vhost_traffic_status_zone;
    ...
```

- Add to the `/etc/nginx/conf.d/default.conf`

```
server {
    ...
    location /status {
        vhost_traffic_status_display;
        vhost_traffic_status_display_format html;
    }
```

- Restart the docker container

```
docker restart nginxvts
```

#### Gather metrics

See [nginx-vts module readme](https://github.com/vozlt/nginx-module-vts#description) for metrics description.

```cmd
curl http://127.0.0.1/status/format/json
```
