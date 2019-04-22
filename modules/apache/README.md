# apache

This module will monitor one or more [`Apache`](https://httpd.apache.org/) servers depending on configuration.

**Requirements:**
 * `Apache` with enabled [`mod_status`](https://httpd.apache.org/docs/2.4/mod/mod_status.html)

It produces the following charts:

1. **Requests** in requests/s
 * requests

2. **Connections** in connections
 * connections

3. **Async Connections** in connections
 * keepalive
 * closing
 * writing
 
4. **Scoreboard** in connections
 * waiting
 * starting
 * reading
 * sending
 * keepalive
 * dns lookup
 * closing
 * logging
 * finishing
 * idle cleanup
 * open

4. **Bandwidth** in kilobits/s
 * sent

5. **Workers** in workers
 * idle
 * busy

6. **Lifetime Average Number Of Requests Per Second** in requests/s
 * requests

7. **Lifetime Average Number Of Bytes Served Per Second** in KiB/s
 * served

8. **Lifetime Average Response Size** in KiB
 * size

### configuration

Needs only `url` to server's `server-status?auto`

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://127.0.0.1/server-status?auto
      
  - name: remote
    url : http://100.64.0.1/server-status?auto
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/apache.conf).

Without configuration, module attempts to connect to `http://127.0.0.1/server-status?auto`

---
