# nginx

> Nginx is a web server which can also be used as a reverse proxy, load balancer, mail proxy and HTTP cache. 

This module will monitor one or more [`nginx`](https://www.nginx.com/) servers via `ngx_http_stub_status_module`.

**Requirements:**
 * `nginx` with configured [`ngx_http_stub_status_module`](http://nginx.org/en/docs/http/ngx_http_stub_status_module.html)


It produces following charts:

1. **Active Client Connections Including Waiting Connections** in connections
 
2. **Active Connections Per Status** in connections
 * reading, writing, waiting
 
3. **Accepted And Handled Connections** in connections/s
 * accepts, handled

4. **Requests** in requests/s

### configuration

Needs only `url` to server's `stub_status`

Here is an example for local server:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/stub_status
      
  - name: remote
    url : http://100.64.0.1/stub_status
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/nginx.conf).

Without configuration, module attempts to connect to `http://127.0.0.1/stub_status`

---
