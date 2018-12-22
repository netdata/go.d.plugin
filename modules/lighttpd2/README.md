# lighttpd2

This module will monitor one or more Lighttpd2 servers depending on configuration.

**Requirements:**
 * lighttpd2 with enabled `mod_status`

It produces the following charts:

1. **Requests** in requests/s
 * requests

2. **Status Codes** in requests/s
 * 1xx
 * 2xx
 * 3xx
 * 4xx
 * 5xx

3. **Traffic** in kilobits/s
 * in
 * out

4. **Connections** in connections
 * connections
 
5. **Connection States** amount connection in every state
 * start
 * read header
 * handle request
 * write response
 * keepalive
 * upgraded
 
6. **Memory Usage** in KiB
 * usage

7. **Uptime** in seconds
 * uptime


### configuration

Needs only `url` to server's `server-status?format=plain`

Here is an example for 2 servers:

```yaml
jobs:
  - name: local
    url : http://localhost/server-status?format=plain
      
  - name: remote
    url : http://100.64.0.1/server-status?format=plain
```

Without configuration, module attempts to connect to `http://localhost/server-status?format=plain`

---
