# tengine

> Tengine is a web server originated by Taobao, the largest e-commerce website in Asia. It is based on the Nginx HTTP server and has many advanced features.

This module will monitor one or more [tengine](https://tengine.taobao.org/) instances via the `ngx_http_reqstat_module` module.

**Requirements:**
 * tengine with configured [`ngx_http_reqstat_module`](http://tengine.taobao.org/document/http_reqstat.html)
 * collector expects [default line format](http://tengine.taobao.org/document/http_reqstat.html)

It produces the following charts:

<br>
Summary:

1. **Bandwidth** in B/s

2. **Connections** in connections/s

3. **Requests** in requests/s

4. **Requests Per Response Code Family** in requests/s
  * 2xx, 3xx, 4xx, 5xx, other
  
5. **Requests Per Response Code Detailed** in requests/s
  * 200, 206, 302, 304, 403, ..., other
  
6. **Number Of Requests Calling For Upstream** in requests/s

7. **Number Of Times Calling For Upstream** in calls/s

7. **Requests Per Response Code Family** in requests/s
  * 4xx, 5xx



### configuration

Needs only `url`.

Here is an example:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/us
      
  - name: remote
    url : http://100.64.0.1us
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/tengine.conf).

Without configuration, module attempts to connect to `http://127.0.0.1/us`

---
