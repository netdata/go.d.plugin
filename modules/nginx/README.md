# nginx

[`Nginx`](https://www.nginx.com/) is a web server which can also be used as a reverse proxy, load balancer, mail proxy and HTTP cache. 

This module will monitor one or more [`Nginx`](https://www.nginx.com/) depending on configuration.

## Requirements

 -   `Nginx` with configured [`ngx_http_stub_status_module`](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html).

## Charts

It produces following charts:

-   Active Client Connections Including Waiting Connections in `connections`
-   Active Connections Per Status in `connections`
-   Accepted And Handled Connections in `connections/s`
-   Requests in `requests/s`

## Configuration

Needs only `url` to server's `stub_status`. Here is an example for local server:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/stub_status
      
  - name: remote
    url : http://203.0.113.10/stub_status
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginx.conf).


## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m nginx
